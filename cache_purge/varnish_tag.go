package cache_purge

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	"strings"
)

const (
	VarnishTagKey = "varnish-tag"
)

func init() {
	TypePurgeCacheMapping[VarnishTagKey] = createVarnishTagPurgeCache
}

var _ types.PurgeCache = &varnishTag{}

type varnishTag struct {
	ctx        *context.Context
	cfg        ConfigVarnish
	projectCfg *config.Project
}

func (v varnishTag) Type() string {
	return VarnishTagKey
}

func (v varnishTag) Purge(events types.Events) {
	for _, event := range events {
		headers := map[string]string{types.HeaderCachePurge: fmt.Sprintf("(%s)", types.GetTagSourcePathHash(event.Path))}
		VarnishDoRequest(
			v.ctx,
			"BAN",
			strings.Join([]string{v.cfg.Server, urltools.JoinUri(v.projectCfg.PrefixPath, event.Path)}, "/"),
			headers,
		)
	}
}

func createVarnishTagPurgeCache(ctx *context.Context, projectCfg *config.Project, cfg config.PurgeCacheConfig) (types.PurgeCache, error) {
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

	instance := &varnishTag{ctx: ctx, projectCfg: projectCfg, cfg: instanceConfig}

	return instance, nil
}
