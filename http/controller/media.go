package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"net/http"
	"strings"
)

func GetMedia(ctx *context.Context, project *config.Project) func(c echo.Context) error {
	return func(c echo.Context) error {
		requestPath := strings.Replace(c.Request().RequestURI, project.PrefixPath, "", 1)
		for _, endpoint := range project.Endpoints {
			opts, errMatch := findMatch(&endpoint, requestPath)
			if errMatch != nil {
				ctx.Logger.Debug(errMatch.Error())
				return c.NoContent(http.StatusNotFound)
			}

			if opts == nil {
				continue
			}

			return c.String(http.StatusNotImplemented, "Not Implemented")
		}
		return c.NoContent(http.StatusNotFound)
	}
}

func findMatch(endpoint *config.Endpoint, path string) (*types.ResizeOption, error) {
	if endpoint.CompiledRegex == nil {
		return &types.ResizeOption{Source: path}, nil
	}

	re := endpoint.CompiledRegex
	if !re.MatchString(path) {
		return nil, nil
	}

	matches := re.FindStringSubmatch(path)
	groupNames := re.SubexpNames()

	params := make(map[string]string)
	for i, match := range matches {
		if i != 0 {
			if groupNames[i] != "" && match != "" {
				params[groupNames[i]] = match
			}
		}
	}

	opts := endpoint.DefaultResizeOpts
	err := mapstructure.Decode(params, &opts)

	return &opts, err
}
