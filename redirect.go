package fasthttpmw

import (
	"github.com/valyala/fasthttp"
)

type (
	// RedirectConfig defines the config for Redirect middleware.
	RedirectConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Status code to be used when redirecting the request.
		// Optional. Default value http.StatusMovedPermanently.
		Code int `json:"code"`
	}
)

const (
	www = "www"
)

var (
	// DefaultRedirectConfig is the default Redirect middleware config.
	DefaultRedirectConfig = RedirectConfig{
		Skipper: DefaultSkipper,
		Code:    fasthttp.StatusMovedPermanently,
	}
)

// HTTPSRedirect redirects http requests to https.
// For example, http://labstack.com will be redirect to https://labstack.com.
//
// Usage `Echo#Pre(HTTPSRedirect())`
func HTTPSRedirect() MW {
	return HTTPSRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSRedirect()`.
func HTTPSRedirectWithConfig(config RedirectConfig) MW {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRedirectConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			req := c.Request
			host := string(req.Host())
			uri := req.URI().String()
			if !c.IsTLS() {
				c.Redirect("https://"+host+uri, config.Code)
				return
			}
			next(c)
			return
		}
	}
}

// HTTPSWWWRedirect redirects http requests to https www.
// For example, http://labstack.com will be redirect to https://www.labstack.com.
//
// Usage `Echo#Pre(HTTPSWWWRedirect())`
func HTTPSWWWRedirect() MW {
	return HTTPSWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSWWWRedirect()`.
func HTTPSWWWRedirectWithConfig(config RedirectConfig) MW {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRedirectConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			req := c.Request
			host := string(req.Host())
			uri := req.URI().String()
			if !c.IsTLS() && host[:3] != www {
				c.Redirect("https://www."+host+uri, config.Code)
				return
			}
			next(c)
			return
		}
	}
}

// HTTPSNonWWWRedirect redirects http requests to https non www.
// For example, http://www.labstack.com will be redirect to https://labstack.com.
//
// Usage `Echo#Pre(HTTPSNonWWWRedirect())`
func HTTPSNonWWWRedirect() MW {
	return HTTPSNonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSNonWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSNonWWWRedirect()`.
func HTTPSNonWWWRedirectWithConfig(config RedirectConfig) MW {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRedirectConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			req := c.Request
			host := string(req.Host())
			uri := req.URI().String()
			if !c.IsTLS() {
				if host[:3] == www {
					c.Redirect("https://"+host[4:]+uri, config.Code)
					return
				}
				c.Redirect("https://"+host+uri, config.Code)
				return
			}
			next(c)
			return
		}
	}
}

// WWWRedirect redirects non www requests to www.
// For example, http://labstack.com will be redirect to http://www.labstack.com.
//
// Usage `Echo#Pre(WWWRedirect())`
func WWWRedirect() MW {
	return WWWRedirectWithConfig(DefaultRedirectConfig)
}

// WWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `WWWRedirect()`.
func WWWRedirectWithConfig(config RedirectConfig) MW {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRedirectConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			req := c.Request
			scheme := string(c.URI().Scheme())
			host := string(req.Host())
			if host[:3] != www {
				uri := req.URI().String()
				c.Redirect(scheme+"://www."+host+uri, config.Code)
				return
			}
			next(c)
			return
		}
	}
}

// NonWWWRedirect redirects www requests to non www.
// For example, http://www.labstack.com will be redirect to http://labstack.com.
//
// Usage `Echo#Pre(NonWWWRedirect())`
func NonWWWRedirect() MW {
	return NonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// NonWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `NonWWWRedirect()`.
func NonWWWRedirectWithConfig(config RedirectConfig) MW {
	if config.Skipper == nil {
		config.Skipper = DefaultRedirectConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			req := c.Request
			scheme := string(c.URI().Scheme())
			host := string(req.Host())
			if host[:3] == www {
				uri := req.URI().String()
				c.Redirect(scheme+"://"+host[4:]+uri, config.Code)
				return
			}
			next(c)
			return
		}
	}
}
