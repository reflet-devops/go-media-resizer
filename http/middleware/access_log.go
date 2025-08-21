package middleware

import (
	builtCtx "context"
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/route"
	"github.com/reflet-devops/go-media-resizer/logger"
	"log/slog"
)

func ConfigureAccessLogMiddleware(e *echo.Echo, ctx *context.Context) error {

	lw, err := logger.NewHandlerRotateWriter(ctx.Fs, ctx.Config.HTTP.AccessLogPath, ctx.Done())
	if err != nil {
		return fmt.Errorf("error creating rotate logger: %v", err)
	}

	sloger := slog.New(slog.NewTextHandler(lw, &slog.HandlerOptions{AddSource: false, Level: slog.LevelInfo}))
	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogStatus:        true,
		LogURI:           true,
		LogRequestID:     true,
		LogProtocol:      true,
		LogMethod:        true,
		LogResponseSize:  true,
		LogContentLength: true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogUserAgent:     true,
		LogError:         true,
		LogHeaders:       []string{"X-Forwarded-For"},
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == route.RouteHealthCheckPing
		},
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				xForwardedFor := c.Request().Header.Get("X-Forwarded-For")
				sloger.LogAttrs(builtCtx.Background(), slog.LevelInfo, "REQUEST",
					slog.String(logger.RemoteIPKey, v.RemoteIP),
					slog.String(logger.RealIPKey, c.RealIP()),
					slog.String(logger.HostKey, v.Host),
					slog.String(logger.ProtocolKey, v.Protocol),
					slog.String(logger.MethodKey, v.Method),
					slog.String(logger.UriKey, v.URI),
					slog.Int(logger.StatusKey, v.Status),
					slog.Int64(logger.ResponseSizeKey, v.ResponseSize),
					slog.String(logger.UserAgentKey, v.UserAgent),
					slog.String(logger.XForwardedForKey, xForwardedFor),
					slog.String(logger.RequestIDKey, v.RequestID),
				)
			} else {
				sloger.LogAttrs(builtCtx.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String(logger.HostKey, v.Host),
					slog.String(logger.RealIPKey, c.RealIP()),
					slog.String(logger.UriKey, v.URI),
					slog.Int(logger.StatusKey, v.Status),
					slog.String(logger.RequestIDKey, v.RequestID),
					slog.String(logger.ErrorKey, v.Error.Error()),
				)
			}
			return nil
		},
	}))
	return nil
}
