package fasthttpmw

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	//"log"
	"fmt"
)

type MW = (func(fasthttp.RequestHandler) fasthttp.RequestHandler)

type RouterWithMW struct {
	*fasthttprouter.Router
	premiddleware []MW
	middleware    []MW
}

// Skipper defines a function to skip middleware. Returning true skips processing
// the middleware.
type Skipper = func(c *fasthttp.RequestCtx) bool

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
	return &RouterWithMW{Router: fasthttprouter.New(), premiddleware: []MW{}, middleware: []MW{}}
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
