package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/parser"
	"github.com/reflet-devops/go-media-resizer/types"
)

func GetMedia(ctx *context.Context, project *config.Project, storage types.Storage) func(c echo.Context) error {
	return func(c echo.Context) error {

		requestPath := strings.Replace(c.Request().RequestURI, fmt.Sprintf("/%s", project.PrefixPath), "", 1)
		for _, endpoint := range project.Endpoints {
			opts := ctx.OptsResizePool.Get().(*types.ResizeOption)
			found, errMatch := parser.ParseOption(&endpoint, project, requestPath, opts)
			if errMatch != nil {
				ctx.Logger.Debug(fmt.Sprintf("%s: %s", errMatch.Error(), requestPath))
				return c.String(http.StatusBadRequest, errMatch.Error())
			}

			if !found {
				continue
			}

			file, errGetFile := storage.GetFile(opts.Source)
			if errGetFile != nil {
				ctx.Logger.Debug(fmt.Sprintf("failed to get file %s: %s", errGetFile.Error(), opts.Source), addLogAttr(c)...)
				return c.String(http.StatusNotFound, "file not found")
			}

			buffer := ctx.BufferPool.Get().(*bytes.Buffer)
			_, errCopy := io.Copy(buffer, file)
			if errCopy != nil {
				resetBuffer(ctx, buffer)
				_ = file.Close()
				return c.String(http.StatusInternalServerError, "buffer copy failed")
			}
			_ = file.Close()

			opts.AddTag(types.GetTagSourcePathHash(opts.Source))

			return SendStream(ctx, c, opts, buffer)
		}
		return c.String(http.StatusNotFound, "file not found")
	}
}
