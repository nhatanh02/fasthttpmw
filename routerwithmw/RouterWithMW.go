package routerwithmw

import (
	fastrouter "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	//"log"
)

type MW (func(fasthttp.RequestHandler) fasthttp.RequestHandler)

type RouterWithMW struct {
	*fastrouter.Router
	premiddleware []MW
	middleware    []MW
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
