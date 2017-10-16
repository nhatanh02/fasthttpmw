package fasthttpmw

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sync"

	//fastrouter "github.com/buaazp/fasthttprouter"
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

			req := &c.Request

			// Based on content length
			if len := int64(req.Header.ContentLength()); len > config.limit {
				fmt.Println("len vs limit: %d vs %d", len, config.limit)
				c.Error(fasthttp.StatusMessage(fasthttp.StatusRequestEntityTooLarge), fasthttp.StatusRequestEntityTooLarge)
				return
			}

			// Based on content read
			// How it works:
			// (1) Get a limitedReader instance from the pool. This pool is shared
			// among all request handlers wrapped by the same mw. This saves lots of
			// circles for allocating/deallocating potentially large blocks of memories
			// to hold the (potentially large) request's body
			// (2) Assign the limitedReader's reader and context fields with the current
			// Request.Body (an io.ReaderCloser in net/http) and echo.Context (rather pointless)
			// (3) After all computation and handling stuffs done, re-deposit the memory block
			// into the pool
			// (4) Assigned the limitedReader to the Request's Body ReaderCloser. This is where
			// the actual work is done, since limitedReader also implements ReaderCloser but
			// with a catch: its Read method wraps around the original Body's Read method to
			// supply read length checking semantics.
			//  		r := pool.Get().(*limitedReader)
			//			r.Reset(req.Body, c)
			//			defer pool.Put(r)
			//			req.Body = r
			// The question is, how to adapt this to fasthttp?
			// Keep it simple and stupid.
			if l := int64(len(req.Body())); l > config.limit {
				c.SetStatusCode(fasthttp.StatusRequestEntityTooLarge)
				c.SetContentType("application/json")
				resp, _ := json.Marshal(ErrStatusRequestEntityTooLarge.Message)
				c.SetBody(resp)
				return
			} else {
				next(c)
				return
			}
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
