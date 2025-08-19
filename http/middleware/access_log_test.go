package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/controller/health"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ConfigureAccessLogMiddleware(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.GET("/200", health.GetPing)
	e.GET("/500", func(c echo.Context) error {
		return fmt.Errorf("hello world")
	})
	e.Use(middleware.RequestID())

	ctx := context.TestContext(nil)

	ctx.Config.HTTP.AccessLogPath = "/var/access.log"

	err := ConfigureAccessLogMiddleware(e, ctx)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/200"), nil)
	req.Host = "127.0.0.1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/200")

	e.ServeHTTP(rec, req)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/500"), nil)
	req.Host = "127.0.0.1"
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/500")

	e.ServeHTTP(rec, req)

	buff, _ := afero.ReadFile(ctx.Fs, "/var/access.log")

	assert.Contains(t, string(buff), "level=INFO msg=REQUEST remote_ip=192.0.2.1 real_ip=192.0.2.1 host=127.0.0.1 protocol=HTTP/1.1 method=GET uri=http://127.0.0.1/200 status=200 response_size")
	assert.Contains(t, string(buff), "level=ERROR msg=REQUEST_ERROR")
}

func Test_ConfigureAccessLogMiddleware_Fail(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.GET("/200", health.GetPing)
	e.Use(middleware.RequestID())

	ctx := context.TestContext(nil)
	ctx.Config.HTTP.AccessLogPath = "/var/access.log"
	ctx.Fs = afero.NewReadOnlyFs(ctx.Fs)

	err := ConfigureAccessLogMiddleware(e, ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating rotate logger: failed to open file: operation not permitted")
}
