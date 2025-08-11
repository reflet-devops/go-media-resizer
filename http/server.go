package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/controller"
	"github.com/reflet-devops/go-media-resizer/http/controller/health"
	"github.com/reflet-devops/go-media-resizer/http/middleware"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
)

type Host struct {
	Echo *echo.Echo
}

const ROUTE_HEALTH_CHECK_PING = "/health/ping"
const ROUTE_CGI_EXTRA_RESIZE = "/cdn-cgi/image/:options/:source"

var MANDATORY_ROUTES = []string{
	ROUTE_HEALTH_CHECK_PING,
}

var CGI_EXTRA_ROUTES = []string{
	ROUTE_CGI_EXTRA_RESIZE,
}

func CreateServerHTTP(ctx *context.Context) *echo.Echo {
	e := createServerHTTP()

	for _, route := range MANDATORY_ROUTES {
		e.GET(route, health.GetPing)
	}
	if ctx.Config.ResizeCGI.Enabled {
		cgiMiddleware := middleware.NewDomainAcceptedBySource(ctx)
		for _, route := range CGI_EXTRA_ROUTES {
			e.GET(route, controller.MediaCGI)
			e.Use(cgiMiddleware.Handler)
		}
	}

	hosts := initRouter(ctx, ctx.Config)
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()

		hostname := urltools.GetHostname(req.Host)
		host := hosts[hostname]

		if host == nil {
			ctx.Logger.Debug(fmt.Sprintf("host not found: %s", hostname))
			err = echo.ErrNotFound
		} else {
			ctx.Logger.Debug(fmt.Sprintf("host found: %s", hostname))
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
		_, found := hosts[project.Hostname]
		if !found {
			hosts[project.Hostname] = &Host{
				Echo: e,
			}
		}
		host := hosts[project.Hostname]
		host.Echo.GET(fmt.Sprintf("%s/*", project.PrefixPath), controller.GetMedia(ctx, &project))
		host.Echo.GET(fmt.Sprintf("%s/webhook", project.PrefixPath), controller.GetWebhook(ctx, &project))

	}
	return hosts
}
