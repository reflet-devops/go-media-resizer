package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	buildinHttp "net/http"
	"strings"
)

func GetMediaCGI(ctx *context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		source := c.Param("source")
		opts := &types.ResizeOption{}

		fileExtension := urltools.GetExtension(source)
		fileType := types.GetType(fileExtension)
		opts.OriginFormat = fileType

		optMap := map[string]interface{}{}
		optRaw := strings.Split(c.Param("options"), ",")

		fileTypeIsValid := types.ValidateType(fileType, ctx.Config.AcceptTypeFiles)
		if !fileTypeIsValid {
			ctx.Logger.Error(fmt.Sprintf("GetMediaCGI: file type not accepted: %s", source))
			return HTTPErrorFileTypeNotAccepted
		}

		for _, optStr := range optRaw {
			optSplit := strings.Split(optStr, "=")
			if len(optSplit) == 2 {
				optMap[optSplit[0]] = optSplit[1]
			}
		}

		err := mapstructure.Decode(optMap, opts)
		if err != nil {
			return c.JSON(buildinHttp.StatusInternalServerError, err.Error())
		}

		return c.JSON(buildinHttp.StatusNotImplemented, fmt.Sprintf("opts: %v source: %s", opts, source))
	}

}
