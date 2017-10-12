package fasthttpmw

import (
	"encoding/base64"
	//	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
)

type (
	// BasicAuthConfig defines the config for BasicAuth middleware.
	BasicAuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Validator is a function to validate BasicAuth credentials.
		// Required.
		Validator BasicAuthValidator

		// Realm is a string to define realm attribute of BasicAuth.
		// Default value "Restricted".
		Realm string
	}

	// BasicAuthValidator defines a function to validate BasicAuth credentials.
	BasicAuthValidator = func(string, string, *fasthttp.RequestCtx) (bool, error)
)

const (
	basic        = "Basic"
	defaultRealm = "Restricted"
)

var (
	// DefaultBasicAuthConfig is the default BasicAuth middleware config.
	DefaultBasicAuthConfig = BasicAuthConfig{
		Skipper: DefaultSkipper,
		Realm:   defaultRealm,
	}
)

// BasicAuth returns an BasicAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func BasicAuth(fn BasicAuthValidator) MW {
	c := DefaultBasicAuthConfig
	c.Validator = fn
	return BasicAuthWithConfig(c)
}

// BasicAuthWithConfig returns an BasicAuth middleware with config.
// See `BasicAuth()`.
func BasicAuthWithConfig(config BasicAuthConfig) MW {
	// Defaults
	if config.Validator == nil {
		panic("echo: basic-auth middleware requires a validator function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultBasicAuthConfig.Skipper
	}
	if config.Realm == "" {
		config.Realm = defaultRealm
	}

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			auth := string(c.Request.Header.Peek(HeaderAuthorization))
			l := len(basic)
			if len(auth) > l+1 && auth[:l] == basic {
				b, err := base64.StdEncoding.DecodeString(auth[l+1:])
				if err != nil {
					//panic(fmt.Errorf("fasthttprouter: invalid Authorization=%s", auth))

					c.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
					return
				}
				cred := string(b)
				for i := 0; i < len(cred); i++ {
					if cred[i] == ':' {
						// Verify credentials
						valid, err := config.Validator(cred[:i], cred[i+1:], c)
						if err != nil {
							//panic(fmt.Errorf("fasthttprouter: unable to validate: err=%v", err))

							c.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
							return
						} else if valid {
							next(c)
							return
						}
					}
				}
			}

			realm := ""
			if config.Realm == defaultRealm {
				realm = defaultRealm
			} else {
				realm = strconv.Quote(config.Realm)
			}

			// Need to return `401` for browsers to pop-up login box.
			c.Response.Header.Set(HeaderWWWAuthenticate, basic+" realm="+realm)
			c.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
			return

		}
	}
}
