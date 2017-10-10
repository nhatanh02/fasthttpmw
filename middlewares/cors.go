package middlewares

import (
	"strconv"
	"strings"

	http "github.com/valyala/fasthttp"
)

type (
	// CORSConfig defines the config for CORS middleware.
	CORSConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// AllowOrigin defines a list of origins that may access the resource.
		// Optional. Default value []string{"*"}.
		AllowOrigins []string `json:"allow_origins"`

		// AllowMethods defines a list methods allowed when accessing the resource.
		// This is used in response to a preflight request.
		// Optional. Default value DefaultCORSConfig.AllowMethods.
		AllowMethods []string `json:"allow_methods"`

		// AllowHeaders defines a list of request headers that can be used when
		// making the actual request. This in response to a preflight request.
		// Optional. Default value []string{}.
		AllowHeaders []string `json:"allow_headers"`

		// AllowCredentials indicates whether or not the response to the request
		// can be exposed when the credentials flag is true. When used as part of
		// a response to a preflight request, this indicates whether or not the
		// actual request can be made using credentials.
		// Optional. Default value false.
		AllowCredentials bool `json:"allow_credentials"`

		// ExposeHeaders defines a whitelist headers that clients are allowed to
		// access.
		// Optional. Default value []string{}.
		ExposeHeaders []string `json:"expose_headers"`

		// MaxAge indicates how long (in seconds) the results of a preflight request
		// can be cached.
		// Optional. Default value 0.
		MaxAge int `json:"max_age"`
	}
)

var (
	// DefaultCORSConfig is the default CORS middleware config.
	DefaultCORSConfig = CORSConfig{
		Skipper:      DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
	}
)

// CORS returns a Cross-Origin Resource Sharing (CORS) middleware.
// See: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
func CORS() MW {
	return CORSWithConfig(DefaultCORSConfig)
}

// CORSWithConfig returns a CORS middleware with config.
// See: `CORS()`.
func CORSWithConfig(config CORSConfig) MW {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCORSConfig.Skipper
	}
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(next http.RequestHandler) http.RequestHandler {
		return func(c *http.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			req := c.Request
			res := c.Response
			origin := req.Header.Peek(HeaderOrigin)
			allowOrigin := ""

			// Check allowed origins
			for _, o := range config.AllowOrigins {
				if o == "*" || o == string(origin) {
					allowOrigin = o
					break
				}
			}

			// Simple request
			if string(req.Header.Method()) != "OPTIONS" {
				res.Header.Add(HeaderVary, HeaderOrigin)
				res.Header.Set(HeaderAccessControlAllowOrigin, allowOrigin)
				if config.AllowCredentials {
					res.Header.Set(HeaderAccessControlAllowCredentials, "true")
				}
				if exposeHeaders != "" {
					res.Header.Set(HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				next(c)
				return
			}

			// Preflight request
			res.Header.Add(HeaderVary, HeaderOrigin)
			res.Header.Add(HeaderVary, HeaderAccessControlRequestMethod)
			res.Header.Add(HeaderVary, HeaderAccessControlRequestHeaders)
			res.Header.Set(HeaderAccessControlAllowOrigin, allowOrigin)
			res.Header.Set(HeaderAccessControlAllowMethods, allowMethods)
			if config.AllowCredentials {
				res.Header.Set(HeaderAccessControlAllowCredentials, "true")
			}
			if allowHeaders != "" {
				res.Header.Set(HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := string(req.Header.Peek(HeaderAccessControlRequestHeaders))
				if h != "" {
					res.Header.Set(HeaderAccessControlAllowHeaders, h)
				}
			}
			if config.MaxAge > 0 {
				res.Header.Set(HeaderAccessControlMaxAge, maxAge)
			}
			c.SetStatusCode(http.StatusNoContent)
			return
		}
	}
}
