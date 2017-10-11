package main

import (
	"fmt"
	//fastrouter "github.com/buaazp/fasthttprouter"
	"bytes"
	"encoding/json"
	"fasthttp-mw/middlewares"
	"fasthttp-mw/routerwithmw"
	"github.com/valyala/fasthttp"
	"log"
	//"strings"
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
	postHandler := func(ctx *fasthttp.RequestCtx) {
		var v *Resp = new(Resp)
		var u Resp
		ctx.Request.BodyWriteTo(v)
		fmt.Println(v)
		jsval, _ := json.Marshal(*v)
		fmt.Println(jsval)
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(jsval)
		//test unmarshalling
		json.Unmarshal(jsval, &u)
		fmt.Println(u)
	}
	router := routerwithmw.New()
	router.Use(middlewares.Recover())
	//router.Use(middlewares.BodyLimit("1B"))
	router.Use(middlewares.BasicAuth(func(username string, password string, c *fasthttp.RequestCtx) (bool, error) {
		if username == "joe" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))
	//router.Use(middlewares.Secure())
	//router.Use(middlewares.JWT(""))
	//router.Use(middlewares.CORS())
	router.POST("/*a", postHandler)
	router.GET("/*a", requestHandler)
	//router.GET("/:a", requestHandler)
	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// ListenAndServe returns only on error, so usually it blocks forever.
	if err := fasthttp.ListenAndServe(listenAddr, router.Handler); err != nil {
		log.Fatalf("error in ListenAndServe: %s", err)
	}
}

//write body to struct example

type Resp struct {
	A string `json:"a"`
	B string `json:"b"`
}

func (r *Resp) Write(p []byte) (n int, err error) {
	//lines(p), with some light trimming
	fields := bytes.Split(bytes.Trim(p, " \n"), []byte("\n"))
	for _, field := range fields {
		//nil check
		if bytes.Compare(field, nil) == 0 {
			continue
		}
		//split each line to a k-v pair
		kv := bytes.SplitN(field, []byte{':'}, 2)
		if kv == nil {
			continue
		}
		k := bytes.Trim(kv[0], " ")
		v := bytes.Trim(kv[1], " ")
		//writing to struct
		switch string(bytes.ToLower(k)) {
		case "a":
			r.A = string(v)
		case "b":
			r.B = string(v)
		}
	}

	return 0, nil
}
