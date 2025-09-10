package http

import (
	"crypto/subtle"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/reflet-devops/go-media-resizer/cache_purge"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/controller"
	"github.com/reflet-devops/go-media-resizer/http/controller/health"
	"github.com/reflet-devops/go-media-resizer/http/middleware"
	"github.com/reflet-devops/go-media-resizer/http/route"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/storage"
	"github.com/reflet-devops/go-media-resizer/types"
)

type Host struct {
	Echo *echo.Echo
}

func CreateServerHTTP(ctx *context.Context) (*echo.Echo, error) {
	e := createServerHTTP()
	e.Logger.SetOutput(os.Stdout)

	if ctx.Config.HTTP.Metrics.Enable {
		configureMetrics(ctx, e)
	}

	extractorTrustOptions, err := getExtractorTrustOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting extractor trusted options: %w", err)
	}

	e.IPExtractor = echo.ExtractIPFromXFFHeader(
		extractorTrustOptions...,
	)

	e.Use(echoMiddleware.RequestID())
	err = middleware.ConfigureAccessLogMiddleware(e, ctx)
	if err != nil {
		return nil, fmt.Errorf("can't set access log middleware: %v", err)
	}

	e.GET(route.HealthCheckPingRoute, health.GetPing)

	if ctx.Config.ResizeCGI.Enabled {
		cgiMiddleware := middleware.NewDomainAcceptedBySource(ctx)

		ctx.Config.ResizeCGI.Headers = types.Headers{}
		for k, v := range ctx.Config.Headers {
			ctx.Config.ResizeCGI.Headers[k] = v
		}
		for k, v := range ctx.Config.ResizeCGI.ExtraHeaders {
			ctx.Config.ResizeCGI.Headers[k] = v
		}

		e.GET(route.CgiExtraResizeRoute, controller.GetMediaCGI(ctx), cgiMiddleware.Handler)
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
			return c.String(http.StatusNotFound, "file not found")
		} else {
			ctx.Logger.Debug(fmt.Sprintf("host found: %s", hostname))
			host.Echo.ServeHTTP(res, req)
		}
		return
	})

	return e, nil
}

func getExtractorTrustOptions(ctx *context.Context) ([]echo.TrustOption, error) {
	trustOptions := []echo.TrustOption{
		echo.TrustLoopback(true),
	}

	for _, ipRange := range ctx.Config.HTTP.ForwardedHeadersTrustedIP {
		_, ipNet, err := net.ParseCIDR(ipRange)
		if err != nil {
			return nil, fmt.Errorf("can't parse CIDR: %v", err)
		}

		trustOptions = append(trustOptions, echo.TrustIPRange(ipNet))
	}

	return trustOptions, nil
}

func createServerHTTP() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetOutput(io.Discard)
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

		if len(project.PurgeCaches) > 0 {
			chanEvents := make(chan types.Events, 2024)
			host.Echo.POST(fmt.Sprintf("%s/webhook", project.PrefixPath), controller.GetWebhook(ctx, chanEvents, &project))
			purgeCaches := []types.PurgeCache{}
			for _, purgeCacheCfg := range project.PurgeCaches {
				purgeCache, errCreatePurge := cache_purge.CreatePurgeCache(ctx, &project, purgeCacheCfg)
				if errCreatePurge != nil {
					return hosts, errCreatePurge
				}
				purgeCaches = append(purgeCaches, purgeCache)
			}
			listenFileChange(ctx, chanEvents, purgeCaches, storageInstance)
		}

	}
	return hosts, nil
}

func configureMetrics(ctx *context.Context, e *echo.Echo) {
	metricsCfg := ctx.Config.HTTP.Metrics

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Registerer: ctx.MetricsRegistry,
	}))

	basicAuthMid := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}

	if metricsCfg.BasicAuth.Enable() {
		basicAuthMid = echoMiddleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte(metricsCfg.BasicAuth.Username)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(metricsCfg.BasicAuth.Password)) == 1 {
				return true, nil
			}
			return false, nil
		})
	}

	e.GET(route.MetricsRoute, echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
		Gatherer: ctx.MetricsRegistry,
	}), basicAuthMid)
}

func listenFileChange(ctx *context.Context, chanEvents chan types.Events, purgeCaches []types.PurgeCache, storageInstance types.Storage) {

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
