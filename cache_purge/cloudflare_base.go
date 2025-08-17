package cache_purge

import (
	"encoding/json"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/valyala/fasthttp"
)

type CloudflareCachePurge struct {
	Files []string `json:"files,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

type ConfigCloudflare struct {
	ZoneId string `mapstructure:"zone_id" validate:"required"`

	AuthEmail string `mapstructure:"auth_email" validate:"required_with=AuthKey"`
	AuthKey   string `mapstructure:"auth_key" validate:"required_with=AuthEmail"`

	AuthToken string `mapstructure:"auth_token" validate:"required_without=AuthEmail AuthKey"`
}

func CloudflareDoRequest(ctx *context.Context, cfg ConfigCloudflare, opts CloudflareCachePurge) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set(fasthttp.HeaderContentType, types.MimeTypeJSON)
	if cfg.AuthEmail != "" && cfg.AuthKey != "" {
		req.Header.Set("X-Auth-Email", cfg.AuthEmail)
		req.Header.Set("X-Auth-Key", cfg.AuthKey)
	} else {
		req.Header.Set(fasthttp.HeaderAuthorization, fmt.Sprintf("Bearer %s", cfg.AuthToken))
	}
	req.SetRequestURI(fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", cfg.ZoneId))
	body, _ := json.Marshal(opts)

	req.SetBodyRaw(body)
	errDoRequest := ctx.HttpClient.DoTimeout(req, resp, ctx.Config.RequestTimeout)

	if errDoRequest != nil {
		ctx.Logger.Error(fmt.Sprintf("cloudflare cache purge: Purge %s", errDoRequest))
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		ctx.Logger.Error(fmt.Sprintf("cloudflare cache purge: Purge %d", resp.StatusCode()))
	}

}
