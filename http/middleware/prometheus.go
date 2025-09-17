package middleware

import (
	"crypto/subtle"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/route"
)

func ConfigurePrometheusMiddleware(ctx *context.Context, e *echo.Echo) {
	metricsCfg := ctx.Config.HTTP.Metrics

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Registerer: ctx.MetricsRegistry,
		Skipper:    prometheusSkipperFunc,
	}))

	basicAuthMid := basicAuthMidDefaultFunc

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

func basicAuthMidDefaultFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func prometheusSkipperFunc(c echo.Context) bool {
	if c.Request().URL.Path == route.HealthCheckPingRoute {
		return true
	}
	return false
}
