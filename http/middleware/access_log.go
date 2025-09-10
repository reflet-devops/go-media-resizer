package middleware

import (
	builtCtx "context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/route"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/logger"
)

func ConfigureAccessLogMiddleware(e *echo.Echo, ctx *context.Context) error {

	lw, err := logger.NewHandlerRotateWriter(ctx.Fs, ctx.Config.HTTP.AccessLogPath, ctx.Done())
	if err != nil {
		return fmt.Errorf("error creating rotate logger: %v", err)
	}

	slogger := slog.New(slog.NewTextHandler(lw, &slog.HandlerOptions{AddSource: false, Level: slog.LevelInfo}))
	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogStatus:        true,
		LogURI:           true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogProtocol:      true,
		LogMethod:        true,
		LogResponseSize:  true,
		LogContentLength: true,
		LogHost:          true,
		LogUserAgent:     true,
		LogError:         true,
		LogLatency:       true,
		LogHeaders:       []string{"X-Forwarded-For"},
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == route.HealthCheckPingRoute
		},
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				xForwardedFor := v.Headers["X-Forwarded-For"]

				slogger.LogAttrs(builtCtx.Background(), slog.LevelInfo, "REQUEST",
					slog.String(logger.RemoteIPKey, urltools.RemovePortNumber(c.Request().RemoteAddr)),
					slog.String(logger.RealIPKey, v.RemoteIP),
					slog.String(logger.HostKey, v.Host),
					slog.String(logger.ProtocolKey, v.Protocol),
					slog.String(logger.MethodKey, v.Method),
					slog.String(logger.UriKey, v.URI),
					slog.Int(logger.StatusKey, v.Status),
					slog.Int64(logger.ResponseSizeKey, v.ResponseSize),
					slog.String(logger.UserAgentKey, v.UserAgent),
					slog.String(logger.XForwardedForKey, fmt.Sprintf("%v", strings.Join(xForwardedFor, ","))),
					slog.String(logger.RequestIDKey, v.RequestID),
					slog.String(logger.LatencyKey, v.Latency.String()),
				)
			} else {
				slogger.LogAttrs(builtCtx.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String(logger.RemoteIPKey, urltools.RemovePortNumber(c.Request().RemoteAddr)),
					slog.String(logger.HostKey, v.Host),
					slog.String(logger.RealIPKey, v.RemoteIP),
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
