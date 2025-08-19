package http

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/route"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
	"net/http"
	"os"
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

	mandatoryRoutes := make([]string, len(route.MandatoryRoutes))
	copy(mandatoryRoutes, route.MandatoryRoutes)
	for _, r := range routes {
		index := slices.Index(mandatoryRoutes, r.Path)
		if index != -1 {
			mandatoryRoutes = slices.Delete(mandatoryRoutes, index, index+1)
		}
	}
	if len(mandatoryRoutes) > 0 {
		assert.Fail(t, fmt.Sprintf("Missing mandatory routes: %v", mandatoryRoutes))
	}
	assert.Equal(t, os.Stdout, e.Logger.Output())
}

func Test_CreateServerHTTP_MidllewareLogger_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.HTTP.AccessLogPath = "/path/log.txt"
	ctx.Fs = afero.NewReadOnlyFs(ctx.Fs)

	_, err := CreateServerHTTP(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can't set access log middleware: ")

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

	mandatoryRoutes := make([]string, len(route.CgiExtraRoutes))
	copy(mandatoryRoutes, route.CgiExtraRoutes)
	for _, r := range routes {
		index := slices.Index(mandatoryRoutes, r.Path)
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

func Test_createServerHTTP(t *testing.T) {
	e := createServerHTTP()
	assert.NotNil(t, e)
	assert.Equal(t, e.Logger.Output(), io.Discard)
}

func Test_CreateServerHTTP(t *testing.T) {

	tests := []struct {
		name       string
		assertEcho func(t *testing.T, ctx *context.Context) *echo.Echo
		checkFn    func(t *testing.T, e *echo.Echo, buff *bytes.Buffer)
	}{
		{
			name: "HostNotFound",
			assertEcho: func(t *testing.T, ctx *context.Context) *echo.Echo {
				ctx.LogLevel.Set(slog.LevelDebug)
				e, err := CreateServerHTTP(ctx)
				assert.NoError(t, err)
				assert.NotNil(t, e)
				return e
			},
			checkFn: func(t *testing.T, e *echo.Echo, buff *bytes.Buffer) {
				go assert.NotPanics(t, func() {
					err := e.Start("127.0.0.1:8081")
					assert.Contains(t, err.Error(), "http: Server closed")
				})
				time.Sleep(200 * time.Millisecond)
				_, err := http.Get(fmt.Sprintf("http://%s", e.Server.Addr))
				assert.Nil(t, err)
				assert.Contains(t, buff.String(), "host not found:")
				_ = e.Close()
			},
		},
		{
			name: "HostFound",
			assertEcho: func(t *testing.T, ctx *context.Context) *echo.Echo {
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
				return e
			},
			checkFn: func(t *testing.T, e *echo.Echo, buff *bytes.Buffer) {
				go assert.NotPanics(t, func() {
					err := e.Start("127.0.0.1:8082")
					assert.Contains(t, err.Error(), "http: Server closed")
				})
				time.Sleep(200 * time.Millisecond)
				_, err := http.Get(fmt.Sprintf("http://%s", e.Server.Addr))
				assert.Nil(t, err)
				assert.Contains(t, buff.String(), "host found:")
				_ = e.Close()
			},
		},
		{
			name: "CGIMiddlewareFail",
			assertEcho: func(t *testing.T, ctx *context.Context) *echo.Echo {

				ctx.Config.ResizeCGI.Enabled = true
				ctx.Config.ResizeCGI.AllowSelfDomain = true

				e, err := CreateServerHTTP(ctx)
				assert.NoError(t, err)
				assert.NotNil(t, e)
				return e
			},
			checkFn: func(t *testing.T, e *echo.Echo, buff *bytes.Buffer) {
				go assert.NotPanics(t, func() {
					err := e.Start("127.0.0.1:8084")
					assert.Contains(t, err.Error(), "http: Server closed")
				})
				time.Sleep(200 * time.Millisecond)

				resp, _ := http.Get(fmt.Sprintf("http://%s/cdn-cgi/image/width=200/127.0.0.2/my/resource", e.Server.Addr))
				assert.Equal(t, http.StatusForbidden, resp.StatusCode)

				_ = e.Close()
			},
		},
		{
			name: "InitRouterFail",
			assertEcho: func(t *testing.T, ctx *context.Context) *echo.Echo {

				ctx.Config.Projects = []config.Project{
					{ID: "id", Hostname: "example.com", Storage: config.StorageConfig{Type: "wrong"}},
				}

				e, err := CreateServerHTTP(ctx)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "project=id, failed to create storage instance: config storage type 'wrong' does not exist")
				return e
			},
			checkFn: func(t *testing.T, e *echo.Echo, buff *bytes.Buffer) {},
		},
		{
			name: "getExtractorTrustedOptionFail",
			assertEcho: func(t *testing.T, ctx *context.Context) *echo.Echo {

				ctx.Config.HTTP.ForwardedHeadersTrustedIP = []string{"a.a.a.a/12"}

				e, err := CreateServerHTTP(ctx)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error getting extractor trusted options:")
				return e
			},
			checkFn: func(t *testing.T, e *echo.Echo, buff *bytes.Buffer) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buff := bytes.NewBufferString("")
			ctx := context.TestContext(buff)
			e := tt.assertEcho(t, ctx)
			tt.checkFn(t, e, buff)
		})
	}
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

	chanEvents := make(chan types.Events, 2024)
	events := types.Events{{Type: types.EventTypePurge, Path: "/test.txt"}}

	storageMock := mockTypes.NewMockStorage(ctrl)
	storageMock.EXPECT().NotifyFileChange(gomock.Any()).Times(1)

	purgeMock := mockTypes.NewMockPurgeCache(ctrl)
	purgeMock.EXPECT().Purge(gomock.Eq(events)).Times(1)
	purgeCaches := []types.PurgeCache{purgeMock}

	listenFileChange(ctx, chanEvents, purgeCaches, storageMock)
	time.Sleep(100 * time.Millisecond)
	chanEvents <- events
	ctx.Cancel()
}

func Test_getExtractorTrustOptions(t *testing.T) {
	tests := []struct {
		name            string
		trustedIpRanges []string
		testFn          func(t *testing.T, got []echo.TrustOption, err error)
	}{
		{
			name:            "Success",
			trustedIpRanges: []string{"192.168.0.0/24", "20.10.1.0/24"},
			testFn: func(t *testing.T, got []echo.TrustOption, err error) {
				assert.NoError(t, err)
				assert.Len(t, got, 3)
			},
		},
		{
			name:            "Fail",
			trustedIpRanges: []string{"not/cidr", "20.10.1.0/24"},
			testFn: func(t *testing.T, got []echo.TrustOption, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "can't parse CIDR:")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			ctx.Config.HTTP.ForwardedHeadersTrustedIP = tt.trustedIpRanges
			got, err := getExtractorTrustOptions(ctx)
			tt.testFn(t, got, err)
		})
	}

}
