package cache_purge

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
)

const (
	CloudflareUrlKey = "cloudflare-url"
)

func init() {
	TypePurgeCacheMapping[CloudflareUrlKey] = createCloudflareUrlPurgeCache
}

var _ types.PurgeCache = &cloudflareUrl{}

type cloudflareUrl struct {
	ctx        *context.Context
	cfg        ConfigCloudflare
	projectCfg *config.Project
}

func (v cloudflareUrl) Type() string {
	return CloudflareUrlKey
}

func (v cloudflareUrl) Purge(events types.Events) {
	for _, event := range events {
		path := urltools.JoinUri(v.projectCfg.PrefixPath, event.Path)
		opts := CloudflareCachePurge{
			Files: []string{
				fmt.Sprintf("http://%s/%s", v.projectCfg.Hostname, path),
				fmt.Sprintf("https://%s/%s", v.projectCfg.Hostname, path),
			},
		}
		CloudflareDoRequest(
			v.ctx,
			v.cfg,
			opts,
		)
	}
}

func createCloudflareUrlPurgeCache(ctx *context.Context, projectCfg *config.Project, cfg config.PurgeCacheConfig) (types.PurgeCache, error) {
	instanceConfig := ConfigCloudflare{}
	err := mapstructure.Decode(cfg.Config, &instanceConfig)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(instanceConfig)
	if err != nil {
		return nil, err
	}

	instance := &cloudflareUrl{ctx: ctx, projectCfg: projectCfg, cfg: instanceConfig}

	return instance, nil
}
