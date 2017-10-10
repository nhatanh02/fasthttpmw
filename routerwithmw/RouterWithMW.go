package routerwithmw

import (
	fastrouter "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	//"log"
	"fmt"
)

type MW (func(fasthttp.RequestHandler) fasthttp.RequestHandler)

type RouterWithMW struct {
	*fastrouter.Router
	premiddleware []MW
	middleware    []MW
}

// Skipper defines a function to skip middleware. Returning true skips processing
// the middleware.
type Skipper func(c *fasthttp.RequestCtx) bool

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(*fasthttp.RequestCtx) bool {
	return false
}

func (r *RouterWithMW) Pre(premw MW) {
	r.premiddleware = append(r.premiddleware, premw)
}

func (r *RouterWithMW) Use(mw MW) {
	r.premiddleware = append(r.premiddleware, mw)
}

func New() *RouterWithMW {
	return &RouterWithMW{Router: fastrouter.New(), premiddleware: []MW{}, middleware: []MW{}}
}

func (r *RouterWithMW) Handler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	//Middleware
	//handler = foldr apply routedHandler r.middleware
	//		where routedHandler = r.Lookup(method,path,ctx)
	handler := func(c *fasthttp.RequestCtx) {
		if h, _ := r.Lookup(method, path, ctx); h != nil {
			for i := len(r.middleware) - 1; i >= 0; i-- {
				h = r.middleware[i](h)
			}
			h(ctx)
		} else {
			r.Router.Handler(ctx)
		}
		return
	}

	// Premiddleware
	//handler' = foldr apply handler r.premiddleware
	for i := len(r.premiddleware) - 1; i >= 0; i-- {
		handler = r.premiddleware[i](handler)
	}
	fmt.Println("All middlewares applied!")
	handler(ctx) //=> premid[0] (premid[1] .... routedHandler...) $ ctx
	return
}

// Auxillary consts
// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=UTF-8"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)
