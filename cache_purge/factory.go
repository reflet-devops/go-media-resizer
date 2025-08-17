package cache_purge

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
)

var TypePurgeCacheMapping = map[string]CreatePurgeCacheFn{}

type CreatePurgeCacheFn func(ctx *context.Context, projectCfg *config.Project, cfg config.PurgeCacheConfig) (types.PurgeCache, error)

func CreatePurgeCache(ctx *context.Context, projectCfg *config.Project, cfg config.PurgeCacheConfig) (types.PurgeCache, error) {
	if fn, ok := TypePurgeCacheMapping[cfg.Type]; ok {
		return fn(ctx, projectCfg, cfg)
	}
	return nil, fmt.Errorf("config purge cache type '%s' does not exist", cfg.Type)
}
