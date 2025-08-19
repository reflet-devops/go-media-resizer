package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/logger"
	"net/http"
	"slices"
)

type DomainAcceptedBySource struct {
	ctx *context.Context
}

func NewDomainAcceptedBySource(ctx *context.Context) *DomainAcceptedBySource {
	return &DomainAcceptedBySource{ctx: ctx}
}

func (d DomainAcceptedBySource) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		source := c.Param("source")
		if source == "" {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		host := c.Request().Host

		isAccepted := d.Validate(source, host)
		if !isAccepted {
			d.ctx.Logger.Error(fmt.Sprintf("domain not allowed: %s", source), logger.RequestIDKey, c.Request().Header.Get(echo.HeaderXRequestID))
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("domain not allowed: %s", source))
		}

		return next(c)
	}
}

func (d DomainAcceptedBySource) Validate(source, server string) bool {
	srcHostname := urltools.GetHostname(source)
	isAccepted := false
	if d.ctx.Config.ResizeCGI.AllowSelfDomain {
		selfHostname := urltools.GetHostname(server)
		if srcHostname == selfHostname {
			isAccepted = true
		}
	}
	if !isAccepted && slices.Contains(d.ctx.Config.ResizeCGI.AllowDomains, srcHostname) {
		isAccepted = true
	}
	return isAccepted
}
