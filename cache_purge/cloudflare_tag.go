package cache_purge

import (
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
)

const (
	CloudflareTagKey = "cloudflare-tag"
)

func init() {
	TypePurgeCacheMapping[CloudflareTagKey] = createCloudflareTagPurgeCache
}

var _ types.PurgeCache = &cloudflareTag{}

type cloudflareTag struct {
	ctx        *context.Context
	cfg        ConfigCloudflare
	projectCfg *config.Project
}

func (v cloudflareTag) Type() string {
	return CloudflareTagKey
}

func (v cloudflareTag) Purge(events types.Events) {
	for _, event := range events {
		fullPath := urltools.FormatPathWithPrefix(v.projectCfg.PrefixPath, event.Path)
		opts := CloudflareCachePurge{
			Tags: []string{types.GetTagSourcePathHash(types.FormatProjectPathHash(v.projectCfg.ID, fullPath))},
		}
		CloudflareDoRequest(
			v.ctx,
			v.cfg,
			opts,
		)
	}
}

func createCloudflareTagPurgeCache(ctx *context.Context, projectCfg *config.Project, cfg config.PurgeCacheConfig) (types.PurgeCache, error) {
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

	instance := &cloudflareTag{ctx: ctx, projectCfg: projectCfg, cfg: instanceConfig}

	return instance, nil
}
