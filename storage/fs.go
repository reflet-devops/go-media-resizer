package storage

import (
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
	"strings"
)

const (
	FsKey = "fs"
)

func init() {
	TypeStorageMapping[FsKey] = createFsStorage
}

var _ types.Storage = &fs{}

type ConfigFs struct {
	PrefixPath string `mapstructure:"prefix_path"`
}

type fs struct {
	fs  afero.Fs
	cfg ConfigFs
}

func (f fs) Type() string {
	return FsKey
}

func (f fs) NotifyFileChange(_ chan types.Events) {}

func (f fs) GetFile(path string) (io.Reader, error) {
	if f.cfg.PrefixPath != "" {
		path = filepath.Join(f.cfg.PrefixPath, path)
	}
	return f.fs.Open(path)
}

func createFsStorage(ctx *context.Context, cfg config.StorageConfig) (types.Storage, error) {
	instanceConfig := ConfigFs{}
	err := mapstructure.Decode(cfg.Config, &instanceConfig)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	// no rule validation
	_ = validate.Struct(instanceConfig)

	instanceConfig.PrefixPath = strings.TrimRight(instanceConfig.PrefixPath, "/")

	instance := &fs{fs: ctx.GetFS(), cfg: instanceConfig}

	return instance, nil
}
