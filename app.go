package main

import (
	"fmt"
	//fastrouter "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"testfasthttp/routerwithmw"
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
	router.Use(routerwithmw.BodyLimit("1B"))
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
