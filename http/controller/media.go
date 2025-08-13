package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	"net/http"
	"strings"
)

func GetMedia(ctx *context.Context, project *config.Project) func(c echo.Context) error {
	return func(c echo.Context) error {
		requestPath := strings.Replace(c.Request().RequestURI, project.PrefixPath, "", 1)
		for _, endpoint := range project.Endpoints {
			opts, errMatch := findMatch(&endpoint, project, requestPath)
			if errMatch != nil {
				ctx.Logger.Debug(fmt.Sprintf("%s: %s", errMatch.Error(), requestPath))
				return errMatch
			}

			if opts == nil {
				continue
			}

			return c.String(http.StatusNotImplemented, "Not Implemented")
		}
		return c.NoContent(http.StatusNotFound)
	}
}

func findMatch(endpoint *config.Endpoint, projectCfg *config.Project, path string) (*types.ResizeOption, error) {
	originExt := urltools.GetExtension(path)
	originType := types.GetType(originExt)

	fileTypeIsValid := types.ValidateType(originType, projectCfg.AcceptTypeFiles)
	if !fileTypeIsValid {
		return nil, HTTPErrorFileTypeNotAccepted
	}

	if endpoint.CompiledRegex == nil {
		return &types.ResizeOption{Source: path, OriginFormat: originType}, nil
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
	opts.OriginFormat = originType

	err := mapstructure.Decode(params, &opts)

	return &opts, err
}
