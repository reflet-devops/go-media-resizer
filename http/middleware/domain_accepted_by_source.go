package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"net/http"
	"strings"
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
		if strings.HasPrefix("http://", source) {
			d.ctx.Logger.Warn(fmt.Sprintf("source has unsecured http protocol: %s", source))
		}
		hostname := urltools.GetHostname(source)

		isAccepted := false

		if d.ctx.Config.ResizeCGI.AllowSelfDomain {
			selfHostname := urltools.GetHostname(c.Request().Host)
			if strings.Compare(hostname, selfHostname) == 0 {
				isAccepted = true
			}
		}

		for _, acceptedHostname := range d.ctx.Config.ResizeCGI.AllowDomains {
			if isAccepted || strings.Compare(hostname, acceptedHostname) == 0 {
				isAccepted = true
				break
			}
		}
		if !isAccepted {
			d.ctx.Logger.Error(fmt.Sprintf("domain not allowed: %s", source))
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("domain not allowed: %s", source))
		}

		return next(c)
	}
}
