package parser

import (
	"errors"
	"strings"

	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/http/urltools"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
)

func ParseOption(endpoint *config.Endpoint, projectCfg *config.Project, uri string, opts *types.ResizeOption) (bool, error) {
	path := strings.Split(uri, "?")[0]
	originExt := urltools.GetExtension(path)
	originType := types.GetType(originExt)

	fileTypeIsValid := types.ValidateType(originType, projectCfg.AcceptTypeFiles)
	if !fileTypeIsValid {
		return false, errors.New("file type not accepted")
	}

	if endpoint.CompiledRegex == nil {
		opts.Source = path
		opts.OriginFormat = originType
		for k, v := range projectCfg.Headers {
			opts.AddHeader(k, v)
		}

		return true, nil
	}

	re := endpoint.CompiledRegex
	if !re.MatchString(uri) {
		return false, nil
	}

	matches := re.FindStringSubmatch(uri)
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
	for k, v := range projectCfg.Headers {
		opts.AddHeader(k, v)
	}

	err := mapstructure.Decode(params, &opts)
	if err != nil {
		return true, err
	}

	opts.Source = strings.Split(opts.Source, "?")[0]

	return true, nil
}
