package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/controller"
	"github.com/reflet-devops/go-media-resizer/http/controller/health"
	"strings"
)

type Host struct {
	Echo *echo.Echo
}

func CreateServerHTTP(ctx *context.Context) *echo.Echo {
	e := createServerHTTP()

	e.GET("/health/ping", health.GetPing)
	if ctx.Config.ResizeCGI.Enabled {
		e.GET("/cdn-cgi/image/:options/:source", controller.MediaCGI)
	}

	hosts := initRouter(ctx, ctx.Config)
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()

		hostname := strings.Split(req.Host, ":")[0]
		host := hosts[hostname]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	return e
}

func createServerHTTP() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return e
}

func initRouter(ctx *context.Context, cfg *config.Config) map[string]*Host {
	hosts := map[string]*Host{}

	for _, project := range cfg.Projects {
		e := createServerHTTP()

		if _, ok := hosts[project.Hostname]; !ok {
			hosts[project.Hostname] = &Host{Echo: e}
		}
		host := hosts[project.Hostname]

		host.Echo.GET(fmt.Sprintf("%s/*", project.PrefixPath), controller.GetMedia(ctx, &project))
		host.Echo.GET(fmt.Sprintf("%s/webhook", project.PrefixPath), controller.GetWebhook(ctx, &project))

	}
	return hosts
}
