package cache_purge

import (
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePurgeCache_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	prjCfg := &config.Project{}
	want := &varnishUrl{ctx: ctx, projectCfg: prjCfg, cfg: ConfigVarnish{Server: "127.0.0.1"}}
	cfg := config.PurgeCacheConfig{Type: VarnishUrlKey, Config: map[string]interface{}{"server": "127.0.0.1"}}
	got, err := CreatePurgeCache(ctx, prjCfg, cfg)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCreatePurgeCache_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	prjCfg := &config.Project{}
	cfg := config.PurgeCacheConfig{Type: "wrong", Config: map[string]interface{}{}}
	got, err := CreatePurgeCache(ctx, prjCfg, cfg)
	assert.Nil(t, got)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config purge cache type 'wrong' does not exist")
}
