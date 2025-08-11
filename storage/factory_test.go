package storage

import (
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateStorage_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	want := &fs{fs: ctx.Fs, cfg: ConfigFs{PrefixPath: "/app"}}
	cfg := config.StorageConfig{Type: FsKey, Config: map[string]interface{}{"prefix_path": "/app"}}
	got, err := CreateStorage(ctx, cfg)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCreateStorage_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	cfg := config.StorageConfig{Type: "wrong", Config: map[string]interface{}{}}
	got, err := CreateStorage(ctx, cfg)
	assert.Nil(t, got)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config storage type 'wrong' does not exist")
}
