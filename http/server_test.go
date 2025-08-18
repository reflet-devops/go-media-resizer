package http

import (
	"bytes"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log/slog"
	"net/http"
	"slices"
	"testing"
	"time"
)

func Test_CreateServerHTTP_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	e, err := CreateServerHTTP(ctx)
	assert.NoError(t, err)

	assert.NotNil(t, e)
	routes := e.Routes()

	mandatoryRoutes := make([]string, len(MandatoryRoutes))
	copy(mandatoryRoutes, MandatoryRoutes)
	for _, route := range routes {
		index := slices.Index(mandatoryRoutes, route.Path)
		if index != -1 {
			mandatoryRoutes = slices.Delete(mandatoryRoutes, index, index+1)
		}
	}
	if len(mandatoryRoutes) > 0 {
		assert.Fail(t, fmt.Sprintf("Missing mandatory routes: %v", mandatoryRoutes))
	}
}

func Test_CreateServerHTTP_CGI_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.ResizeCGI.Enabled = true
	ctx.Config.Headers = types.Headers{"X-Custom": "foo", "X-Override": "foo"}
	ctx.Config.ResizeCGI.ExtraHeaders = types.Headers{"X-Override": "bar"}

	e, err := CreateServerHTTP(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, e)
	routes := e.Routes()

	mandatoryRoutes := make([]string, len(CgiExtraRoutes))
	copy(mandatoryRoutes, CgiExtraRoutes)
	for _, route := range routes {
		index := slices.Index(mandatoryRoutes, route.Path)
		if index != -1 {
			mandatoryRoutes = slices.Delete(mandatoryRoutes, index, index+1)
		}
	}
	if len(mandatoryRoutes) > 0 {
		assert.Fail(t, fmt.Sprintf("Missing mandatory CGI routes: %v", mandatoryRoutes))
	}

	wantHeader := types.Headers{"X-Custom": "foo", "X-Override": "bar"}
	assert.Equal(t, wantHeader, ctx.Config.ResizeCGI.Headers)
}

func Test_CreateServerHTTP_HostNotFound_Fail(t *testing.T) {
	buffer := bytes.NewBufferString("")
	ctx := context.TestContext(buffer)
	ctx.LogLevel.Set(slog.LevelDebug)

	e, err := CreateServerHTTP(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, e)

	go assert.NotPanics(t, func() {
		err := e.Start("127.0.0.1:8081")
		assert.Contains(t, err.Error(), "http: Server closed")
	})
	time.Sleep(200 * time.Millisecond)

	_, err = http.Get(fmt.Sprintf("http://%s", e.Server.Addr))
	assert.Nil(t, err)
	assert.Contains(t, buffer.String(), "host not found:")

	_ = e.Close()
}

func Test_CreateServerHTTP_HostFound_Success(t *testing.T) {
	buffer := bytes.NewBufferString("")
	ctx := context.TestContext(buffer)
	ctx.LogLevel.Set(slog.LevelDebug)

	ctx.Config.Projects = []config.Project{
		{
			ID:       "localhost",
			Hostname: "127.0.0.1",
			Storage:  config.StorageConfig{Type: "fs", Config: map[string]interface{}{"prefix_path": "/app"}},
		},
	}

	e, err := CreateServerHTTP(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, e)

	go assert.NotPanics(t, func() {
		err := e.Start("127.0.0.1:8082")
		assert.Contains(t, err.Error(), "http: Server closed")
	})
	time.Sleep(200 * time.Millisecond)

	_, err = http.Get(fmt.Sprintf("http://%s", e.Server.Addr))
	assert.Nil(t, err)
	assert.Contains(t, buffer.String(), "host found:")
	_ = e.Close()
}

func Test_CreateServerHTTP_CGIMiddleware_Fail(t *testing.T) {
	buffer := bytes.NewBufferString("")
	ctx := context.TestContext(buffer)
	ctx.Config.ResizeCGI.Enabled = true
	ctx.Config.ResizeCGI.AllowSelfDomain = true
	ctx.LogLevel.Set(slog.LevelDebug)

	e, err := CreateServerHTTP(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, e)

	go assert.NotPanics(t, func() {
		err := e.Start("127.0.0.1:8084")
		assert.Contains(t, err.Error(), "http: Server closed")
	})
	time.Sleep(200 * time.Millisecond)

	resp, _ := http.Get(fmt.Sprintf("http://%s/cdn-cgi/image/width=200/127.0.0.2/my/resource", e.Server.Addr))
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	_ = e.Close()
}

func Test_CreateServerHTTP_initRouter_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.Projects = []config.Project{
		{ID: "id", Hostname: "example.com", Storage: config.StorageConfig{Type: "wrong"}},
	}

	e, err := CreateServerHTTP(ctx)
	assert.Error(t, err)
	assert.NotNil(t, e)
	assert.Contains(t, err.Error(), "project=id, failed to create storage instance: config storage type 'wrong' does not exist")
}

func Test_initRouter_WithPrefix_Success(t *testing.T) {
	ctx := context.TestContext(nil)

	ctx.Config.Projects = []config.Project{
		{
			ID:         "with-prefix",
			Hostname:   "with-prefix.com",
			PrefixPath: "prefix",
			Storage:    config.StorageConfig{Type: "fs", Config: map[string]interface{}{"prefix_path": "/app"}},
		},
	}

	hosts, err := initRouter(ctx, ctx.Config)
	assert.NoError(t, err)
	host, found := hosts["with-prefix.com"]
	assert.True(t, found)

	mandatoryRoutes := []string{
		"/prefix/*",
		"/prefix/webhook",
	}
	for _, route := range host.Echo.Routes() {
		index := slices.Index(mandatoryRoutes, route.Path)
		if index != -1 {
			mandatoryRoutes = slices.Delete(mandatoryRoutes, index, index+1)
		}
	}
	if len(mandatoryRoutes) > 0 {
		assert.Fail(t, fmt.Sprintf("Missing mandatory routes: %v", mandatoryRoutes))
	}
}

func Test_initRouter_NoPrefix_Success(t *testing.T) {
	ctx := context.TestContext(nil)

	ctx.Config.Projects = []config.Project{
		{
			ID:       "no-prefix",
			Hostname: "no-prefix.com",
			Storage:  config.StorageConfig{Type: "fs", Config: map[string]interface{}{"prefix_path": "/app"}},
		},
	}

	hosts, err := initRouter(ctx, ctx.Config)
	assert.NoError(t, err)
	host, found := hosts["no-prefix.com"]
	assert.True(t, found)

	mandatoryRoutes := []string{
		"/*",
		"/webhook",
	}
	for _, route := range host.Echo.Routes() {
		index := slices.Index(mandatoryRoutes, route.Path)
		if index != -1 {
			mandatoryRoutes = slices.Delete(mandatoryRoutes, index, index+1)
		}
	}
	if len(mandatoryRoutes) > 0 {
		assert.Fail(t, fmt.Sprintf("Missing mandatory routes: %v", mandatoryRoutes))
	}
}

func Test_initRouter_WithPurgeCache_Success(t *testing.T) {
	ctx := context.TestContext(nil)

	ctx.Config.Projects = []config.Project{
		{
			ID:          "id",
			Hostname:    "example.com",
			Storage:     config.StorageConfig{Type: "fs", Config: map[string]interface{}{"prefix_path": "/app"}},
			PurgeCaches: []config.PurgeCacheConfig{{Type: "varnish-url", Config: map[string]interface{}{"server": "127.0.0.1"}}},
		},
	}

	_, err := initRouter(ctx, ctx.Config)
	assert.NoError(t, err)
}

func Test_initRouter_CreatePurgeCache_Failed(t *testing.T) {
	ctx := context.TestContext(nil)

	ctx.Config.Projects = []config.Project{
		{
			ID:          "id",
			Hostname:    "example.com",
			Storage:     config.StorageConfig{Type: "fs", Config: map[string]interface{}{"prefix_path": "/app"}},
			PurgeCaches: []config.PurgeCacheConfig{{Type: "wrong"}},
		},
	}

	_, err := initRouter(ctx, ctx.Config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config purge cache type 'wrong' does not exist")
}

func Test_initRouter_CreateStorage_Fail(t *testing.T) {
	ctx := context.TestContext(nil)

	ctx.Config.Projects = []config.Project{
		{
			ID:       "id",
			Hostname: "example.com",
			Storage:  config.StorageConfig{Type: "wrong"},
		},
	}

	_, err := initRouter(ctx, ctx.Config)
	assert.Error(t, err)
}

func Test_listenFileChange_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var chanEvents chan types.Events
	events := types.Events{{Type: types.EventTypePurge, Path: "/test.txt"}}

	storageMock := mockTypes.NewMockStorage(ctrl)
	storageMock.EXPECT().NotifyFileChange(gomock.Any()).Do(func(cEvents chan types.Events) { chanEvents = cEvents }).Times(1)

	purgeMock := mockTypes.NewMockPurgeCache(ctrl)
	purgeMock.EXPECT().Purge(gomock.Eq(events)).Times(1)
	purgeCaches := []types.PurgeCache{purgeMock}

	listenFileChange(ctx, purgeCaches, storageMock)
	time.Sleep(100 * time.Millisecond)
	chanEvents <- events
	ctx.Cancel()
}
