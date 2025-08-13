package storage

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
)

var TypeStorageMapping = map[string]CreateStorageFn{}

type CreateStorageFn func(ctx *context.Context, cfg config.StorageConfig) (types.Storage, error)

func CreateStorage(ctx *context.Context, cfg config.StorageConfig) (types.Storage, error) {
	if fn, ok := TypeStorageMapping[cfg.Type]; ok {
		return fn(ctx, cfg)
	}
	return nil, fmt.Errorf("config storage type '%s' does not exist", cfg.Type)
}
