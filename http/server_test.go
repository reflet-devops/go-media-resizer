package http

import (
	"bytes"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"slices"
	"testing"
	"time"
)

func Test_CreateServerHTTP_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	e := CreateServerHTTP(ctx)

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
	e := CreateServerHTTP(ctx)
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
}

func Test_CreateServerHTTP_HostNotFound_Fail(t *testing.T) {
	buffer := bytes.NewBufferString("")
	ctx := context.TestContext(buffer)
	ctx.LogLevel.Set(slog.LevelDebug)

	e := CreateServerHTTP(ctx)
	assert.NotNil(t, e)

	go assert.NotPanics(t, func() {
		err := e.Start("127.0.0.1:8081")
		assert.Contains(t, err.Error(), "http: Server closed")
	})
	time.Sleep(200 * time.Millisecond)

	_, err := http.Get(fmt.Sprintf("http://%s", e.Server.Addr))
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
		},
	}

	e := CreateServerHTTP(ctx)
	assert.NotNil(t, e)

	go assert.NotPanics(t, func() {
		err := e.Start("127.0.0.1:8082")
		assert.Contains(t, err.Error(), "http: Server closed")
	})
	time.Sleep(200 * time.Millisecond)

	_, err := http.Get(fmt.Sprintf("http://%s", e.Server.Addr))
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

	e := CreateServerHTTP(ctx)
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

func Test_initRouter_WithPrefix_Success(t *testing.T) {
	ctx := context.TestContext(nil)

	ctx.Config.Projects = []config.Project{
		{
			ID:         "with-prefix",
			Hostname:   "with-prefix.com",
			PrefixPath: "prefix",
		},
	}

	hosts := initRouter(ctx, ctx.Config)
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
		},
	}

	hosts := initRouter(ctx, ctx.Config)
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
