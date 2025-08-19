package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/parser"
	"github.com/reflet-devops/go-media-resizer/types"
	"net/http"
	"strings"
)

func GetMedia(ctx *context.Context, project *config.Project, storage types.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx = prepareContext(ctx, c)

		requestPath := strings.Replace(c.Request().RequestURI, project.PrefixPath, "", 1)
		for _, endpoint := range project.Endpoints {
			opts, errMatch := parser.ParseOption(&endpoint, project, requestPath)
			if errMatch != nil {
				ctx.Logger.Debug(fmt.Sprintf("%s: %s", errMatch.Error(), requestPath))
				return c.String(http.StatusBadRequest, errMatch.Error())
			}

			if opts == nil {
				continue
			}

			content, errGetFile := storage.GetFile(opts.Source)
			if errGetFile != nil {
				ctx.Logger.Debug(fmt.Sprintf("failed to get file %s: %s", errGetFile.Error(), opts.Source))
				return c.String(http.StatusNotFound, "file not found")
			}
			opts.AddTag(types.GetTagSourcePathHash(opts.Source))
			return SendStream(ctx, c, opts, content)
		}
		return c.String(http.StatusNotFound, "file not found")
	}
}
