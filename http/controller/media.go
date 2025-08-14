package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/parser"
	"net/http"
	"strings"
)

func GetMedia(ctx *context.Context, project *config.Project) func(c echo.Context) error {
	return func(c echo.Context) error {
		requestPath := strings.Replace(c.Request().RequestURI, project.PrefixPath, "", 1)
		for _, endpoint := range project.Endpoints {
			opts, errMatch := parser.ParseOption(&endpoint, project, requestPath)
			if errMatch != nil {
				ctx.Logger.Debug(fmt.Sprintf("%s: %s", errMatch.Error(), requestPath))
				return echo.NewHTTPError(http.StatusBadRequest, errMatch.Error())
			}

			if opts == nil {
				continue
			}

			return c.String(http.StatusNotImplemented, "Not Implemented")
		}
		return c.NoContent(http.StatusNotFound)
	}
}
