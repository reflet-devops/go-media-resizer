package middleware

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/reflet-devops/go-media-resizer/http/route"
)

func ConfigureRequestIdMiddleware(e *echo.Echo) {
	e.Use(echoMiddleware.RequestIDWithConfig(echoMiddleware.RequestIDConfig{
		Skipper: requestIdSkipperFunc,
	}))
}

func requestIdSkipperFunc(c echo.Context) bool {
	if c.Request().URL.Path == route.HealthCheckPingRoute {
		return true
	}
	if c.Request().Header.Get(echo.HeaderXRequestID) != "" {
		return true
	}
	return false
}
