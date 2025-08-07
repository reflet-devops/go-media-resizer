package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"net/http"
)

func GetWebhook(ctx *context.Context, project *config.Project) func(c echo.Context) error {
	return func(c echo.Context) error {

		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
