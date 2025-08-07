package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	buildinHttp "net/http"
	"strings"
)

func MediaCGI(c echo.Context) (err error) {
	optMap := map[string]interface{}{}
	optRaw := strings.Split(c.Param("options"), ",")
	for _, optStr := range optRaw {
		optSplit := strings.Split(optStr, "=")
		if len(optSplit) == 2 {
			optMap[optSplit[0]] = optSplit[1]
		}
	}
	opts := &types.ResizeOption{}
	err = mapstructure.Decode(optMap, opts)
	if err != nil {
		return c.JSON(buildinHttp.StatusInternalServerError, err.Error())
	}
	source := c.Param("source")

	return c.JSON(buildinHttp.StatusNotImplemented, fmt.Sprintf("opts: %v source: %s", opts, source))
}
