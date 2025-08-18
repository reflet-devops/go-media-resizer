package parser

import (
	"errors"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
)

func ParseOption(endpoint *config.Endpoint, projectCfg *config.Project, path string) (*types.ResizeOption, error) {
	originExt := urltools.GetExtension(path)
	originType := types.GetType(originExt)

	fileTypeIsValid := types.ValidateType(originType, projectCfg.AcceptTypeFiles)
	if !fileTypeIsValid {
		return nil, errors.New("file type not accepted")
	}

	if endpoint.CompiledRegex == nil {
		return &types.ResizeOption{Source: path, OriginFormat: originType, Headers: projectCfg.Headers}, nil
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
	opts.Headers = projectCfg.Headers

	err := mapstructure.Decode(params, &opts)

	return &opts, err
}
