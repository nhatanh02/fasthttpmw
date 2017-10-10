package middlewares

import (
	"bufio"
	"fmt"
	"sync"

	//fastrouter "github.com/buaazp/fasthttprouter"
	. "fasthttp-mw/routerwithmw"
	"github.com/labstack/gommon/bytes"
	"github.com/valyala/fasthttp"
)

type (
	// BodyLimitConfig defines the config for BodyLimit middleware.
	BodyLimitConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Maximum allowed size for a request body, it can be specified
		// as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
		Limit string `json:"limit"`
		limit int64
	}

	limitedReader struct {
		BodyLimitConfig
		reader  bufio.Reader
		read    int64
		context fasthttp.RequestCtx
	}
)

var (
	// DefaultBodyLimitConfig is the default BodyLimit middleware config.
	DefaultBodyLimitConfig = BodyLimitConfig{
		Skipper: DefaultSkipper,
	}
)

// BodyLimit returns a BodyLimit middleware.
//
// BodyLimit middleware sets the maximum allowed size for a request body, if the
// size exceeds the configured limit, it sends "413 - Request Entity Too Large"
// response. The BodyLimit is determined based on both `Content-Length` request
// header and actual content read, which makes it super secure.
// Limit can be specified as `4x` or `4xB`, where x is one of the multiple from K, M,
// G, T or P.
func BodyLimit(limit string) MW {
	c := DefaultBodyLimitConfig
	c.Limit = limit
	fmt.Println("Calling BodyLimitWithConfig..")
	return BodyLimitWithConfig(c)
}

// BodyLimitWithConfig returns a BodyLimit middleware with config.
// See: `BodyLimit()`.
func BodyLimitWithConfig(config BodyLimitConfig) MW {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultBodyLimitConfig.Skipper
	}

	limit, err := bytes.Parse(config.Limit)
	if err != nil {
		panic(fmt.Errorf("fasthttprouter: invalid body-limit=%s", config.Limit))
	}
	config.limit = limit
	fmt.Println("Making new handler with body limit = %d", config.limit)
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			if config.Skipper(c) {
				next(c)
				return
			}

			req := c.Request

			// Based on content length
			if len := int64(req.Header.ContentLength()); len > config.limit {
				fmt.Println("len vs limit: %d vs %d", len, config.limit)
				c.Error(fasthttp.StatusMessage(fasthttp.StatusRequestEntityTooLarge), fasthttp.StatusRequestEntityTooLarge)
				return
			}

			//TODO: adapting this logic, from for http.Request to
			//fasthttp.Request, or fasthttp.RequestCtx
			// Based on content read
			//	pool := limitedReaderPool(config)
			//	r := pool.Get().(*limitedReader)
			//	r.Reset(&req, c)
			//	defer pool.Put(r)
			//	req.SetBody(r)

			next(c)
			return
		}
	}
}

func (r *limitedReader) Read(b []byte) (n int, err error) {
	n, err = r.reader.Read(b)
	r.read += int64(n)
	if r.read > r.limit {
		return n, err
	}
	return
}

//func (r *limitedReader) Close() error {
//	return r.reader.Close()
//}

func (r *limitedReader) Reset(reader bufio.Reader, context fasthttp.RequestCtx) {
	r.reader = reader
	r.context = context
}

func limitedReaderPool(c BodyLimitConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &limitedReader{BodyLimitConfig: c}
		},
	}
}
