package cache_purge

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/valyala/fasthttp"
)

type ConfigVarnish struct {
	Server string `mapstructure:"server" validate:"required"`
}

func VarnishDoRequest(ctx *context.Context, method, uri string, headers map[string]string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	req.Header.SetMethod(method)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.SetRequestURI(uri)

	err := ctx.HttpClient.DoTimeout(req, resp, ctx.Config.RequestTimeout)

	if err != nil {
		ctx.Logger.Error(fmt.Sprintf("varnish cache purge: Purge %s: error with request: %v", uri, err))
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		ctx.Logger.Error(fmt.Sprintf("varnish cache purge: Purge %s: invalid status code status code: %d", uri, resp.StatusCode()))
	}

}
