package controller

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"net/http"
)

func GetWebhook(_ *context.Context, chanEvents chan types.Events, project *config.Project) func(c echo.Context) error {
	return func(c echo.Context) error {
		if project.WebhookToken != "" && c.Request().Header.Get(echo.HeaderAuthorization) != fmt.Sprintf("Bearer %s", project.WebhookToken) {
			return c.NoContent(http.StatusUnauthorized)
		}

		events := types.Events{}
		errBind := c.Bind(&events)
		if errBind != nil {
			return c.String(http.StatusBadRequest, errBind.Error())
		}

		validate := validator.New()
		err := validate.Var(events, "required,min=1,dive")
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		chanEvents <- events

		return c.String(http.StatusAccepted, "ok")
	}
}
