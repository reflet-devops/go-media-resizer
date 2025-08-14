package types

import (
	"github.com/valyala/fasthttp"
	"time"
)

type Client interface {
	DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error
}
