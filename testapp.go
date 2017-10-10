package main

import (
	"fmt"
	//fastrouter "github.com/buaazp/fasthttprouter"
	"fasthttp-mw/middlewares"
	"fasthttp-mw/routerwithmw"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	listenAddr := "127.0.0.1:9876"

	// This function will be called by the server for each incoming request.
	//
	// RequestCtx provides a lot of functionality related to http request
	// processing. See RequestCtx docs for details.
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		fmt.Println(string(ctx.RequestURI()))
		switch path := string(ctx.Path()); path {
		case "/":
			fmt.Fprintf(ctx, "Root")
		default:
			//fmt.Fprintf(ctx, "Not found: %q", path)
			//ctx.NotFound()
			fmt.Fprintf(ctx, "%s", ctx.UserValue("a"))
		}
	}
	router := routerwithmw.New()
	router.Use(middlewares.Recover())
	router.Use(middlewares.BodyLimit("1B"))
	router.Use(middlewares.BasicAuth(func(username string, password string, c *fasthttp.RequestCtx) (bool, error) {
		if username == "joe" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))
	router.Use(middlewares.Secure())
	router.Use(middlewares.JWT(""))
	router.Use(middlewares.CORS())
	router.POST("/*a", requestHandler)
	//router.GET("/*a", requestHandler)
	//router.GET("/:a", requestHandler)
	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// ListenAndServe returns only on error, so usually it blocks forever.
	if err := fasthttp.ListenAndServe(listenAddr, router.Handler); err != nil {
		log.Fatalf("error in ListenAndServe: %s", err)
	}
}
