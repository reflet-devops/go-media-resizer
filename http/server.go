package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/cache_purge"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/controller"
	"github.com/reflet-devops/go-media-resizer/http/controller/health"
	"github.com/reflet-devops/go-media-resizer/http/middleware"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/storage"
	"github.com/reflet-devops/go-media-resizer/types"
)

type Host struct {
	Echo *echo.Echo
}

const RouteHealthCheckPing = "/health/ping"
const RouteCgiExtraResize = "/cdn-cgi/image/:options/:source"

var MandatoryRoutes = []string{
	RouteHealthCheckPing,
}

var CgiExtraRoutes = []string{
	RouteCgiExtraResize,
}

func CreateServerHTTP(ctx *context.Context) (*echo.Echo, error) {
	e := createServerHTTP()

	for _, route := range MandatoryRoutes {
		e.GET(route, health.GetPing)
	}
	if ctx.Config.ResizeCGI.Enabled {
		cgiMiddleware := middleware.NewDomainAcceptedBySource(ctx)
		for _, route := range CgiExtraRoutes {
			e.GET(route, controller.GetMediaCGI(ctx), cgiMiddleware.Handler)
		}
	}

	hosts, err := initRouter(ctx, ctx.Config)
	if err != nil {
		return e, err
	}

	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()

		hostname := urltools.RemovePortNumber(req.Host)
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

	return e, nil
}

func createServerHTTP() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return e
}

func initRouter(ctx *context.Context, cfg *config.Config) (map[string]*Host, error) {
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
		storageInstance, err := storage.CreateStorage(ctx, project.Storage)
		if err != nil {
			return hosts, fmt.Errorf("project=%s, failed to create storage instance: %v", project.ID, err)
		}
		host.Echo.GET(fmt.Sprintf("%s/*", project.PrefixPath), controller.GetMedia(ctx, &project, storageInstance))
		host.Echo.GET(fmt.Sprintf("%s/webhook", project.PrefixPath), controller.GetWebhook(ctx, &project))

		if len(project.PurgeCaches) > 0 {
			purgeCaches := []types.PurgeCache{}
			for _, purgeCacheCfg := range project.PurgeCaches {
				purgeCache, errCreatePurge := cache_purge.CreatePurgeCache(ctx, &project, purgeCacheCfg)
				if errCreatePurge != nil {
					return hosts, errCreatePurge
				}
				purgeCaches = append(purgeCaches, purgeCache)
			}
			listenFileChange(ctx, purgeCaches, storageInstance)
		}

	}
	return hosts, nil
}

func listenFileChange(ctx *context.Context, purgeCaches []types.PurgeCache, storageInstance types.Storage) {
	chanEvents := make(chan types.Events, 2024)

	go func(chanEvents chan types.Events, purgeCaches []types.PurgeCache) {
		for {
			select {
			case event := <-chanEvents:
				for _, purgeCache := range purgeCaches {
					purgeCache.Purge(event)
				}
			case <-ctx.Done():
				return
			}
		}
	}(chanEvents, purgeCaches)
	go storageInstance.NotifyFileChange(chanEvents)
}
