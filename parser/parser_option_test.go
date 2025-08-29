package parser

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_ParseOption(t *testing.T) {

	tests := []struct {
		name       string
		endpoint   *config.Endpoint
		projectCfg *config.Project
		path       string

		want    *types.ResizeOption
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:       "successWithNoRegex",
			endpoint:   &config.Endpoint{},
			projectCfg: &config.Project{AcceptTypeFiles: []string{types.TypePNG}, Headers: types.Headers{"X-Custom": "foo"}},
			path:       "/media/image.png",
			want:       &types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png", Headers: types.Headers{"X-Custom": "foo"}},
			wantErr:    assert.NoError,
		},
		{
			name:       "successWithRegex",
			endpoint:   &config.Endpoint{Regex: "(?<source>.*)"},
			projectCfg: &config.Project{AcceptTypeFiles: []string{types.TypePNG}, Headers: types.Headers{"X-Custom": "foo"}},
			path:       "/media/image.png",
			want:       &types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png", Headers: types.Headers{"X-Custom": "foo"}},
			wantErr:    assert.NoError,
		},
		{
			name:       "successWithRegexAndSizeOpts",
			endpoint:   &config.Endpoint{Regex: "\\/(?<width>[0-9]{1,4})\\/(?<height>[0-9]{1,4})(?<source>.*)"},
			projectCfg: &config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			path:       "/500/500/media/image.png",
			want:       &types.ResizeOption{OriginFormat: types.TypePNG, Width: 500, Height: 500, Source: "/media/image.png", Headers: types.Headers{}},
			wantErr:    assert.NoError,
		},
		{
			name:       "failedWithFileTypeNotAccepted",
			endpoint:   &config.Endpoint{},
			projectCfg: &config.Project{AcceptTypeFiles: []string{}},
			path:       "/media/image.png",
			wantErr:    assert.Error,
		},
		{
			name:       "failedWithEndpointNotMatch",
			endpoint:   &config.Endpoint{Regex: "/media/wrong/(?<source>.*)"},
			projectCfg: &config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			path:       "/media/image.png",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.endpoint.Regex != "" {
				compiledRegex, errReCompile := regexp.Compile(tt.endpoint.Regex)
				assert.NoError(t, errReCompile)
				tt.endpoint.CompiledRegex = compiledRegex
			}

			got, err := ParseOption(tt.endpoint, tt.projectCfg, tt.path)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseOption(%v, %v, %v)", tt.endpoint, tt.projectCfg, tt.path)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseOption(%v, %v, %v)", tt.endpoint, tt.projectCfg, tt.path)
		})
	}
}
