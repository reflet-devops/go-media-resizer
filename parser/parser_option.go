package parser

import (
	"errors"
	"strings"

	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
)

func ParseOption(endpoint *config.Endpoint, projectCfg *config.Project, path string, opts *types.ResizeOption) (bool, error) {
	originExt := urltools.GetExtension(path)
	originType := types.GetType(originExt)

	fileTypeIsValid := types.ValidateType(originType, projectCfg.AcceptTypeFiles)
	if !fileTypeIsValid {
		return false, errors.New("file type not accepted")
	}

	if endpoint.CompiledRegex == nil {
		opts.Source = path
		opts.OriginFormat = originType
		opts.Headers = projectCfg.Headers
		return true, nil
	}

	re := endpoint.CompiledRegex
	if !re.MatchString(path) {

		return false, nil
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
	opts.ResetToDefaults(&endpoint.DefaultResizeOpts)
	opts.OriginFormat = originType
	opts.Headers = types.Headers{}
	for k, v := range projectCfg.Headers {
		opts.Headers[k] = v
	}
	opts.Source = strings.Trim(opts.Source, "/")

	err := mapstructure.Decode(params, &opts)

	return true, err
}
