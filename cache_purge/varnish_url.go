package cache_purge

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
)

const (
	VarnishUrlKey = "varnish-url"
)

func init() {
	TypePurgeCacheMapping[VarnishUrlKey] = createVarnishUrlPurgeCache
}

var _ types.PurgeCache = &varnishUrl{}

type varnishUrl struct {
	ctx        *context.Context
	cfg        ConfigVarnish
	projectCfg *config.Project
}

func (v varnishUrl) Type() string {
	return VarnishUrlKey
}

func (v varnishUrl) Purge(events types.Events) {
	for _, event := range events {
		fullPath := urltools.FormatPathWithPrefix(v.projectCfg.PrefixPath, event.Path)
		VarnishDoRequest(
			v.ctx,
			"PURGE",
			strings.Join([]string{v.cfg.Server, fullPath}, "/"),
			nil,
		)
	}
}

func createVarnishUrlPurgeCache(ctx *context.Context, projectCfg *config.Project, cfg config.PurgeCacheConfig) (types.PurgeCache, error) {
	instanceConfig := ConfigVarnish{}
	err := mapstructure.Decode(cfg.Config, &instanceConfig)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(instanceConfig)
	if err != nil {
		return nil, err
	}

	instance := &varnishUrl{ctx: ctx, projectCfg: projectCfg, cfg: instanceConfig}

	return instance, nil
}
