package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func Test_GetMedia(t *testing.T) {
	ctx := context.TestContext(nil)
	e := echo.New()

	tests := []struct {
		name     string
		resource string
		want     error
		prjConf  *config.Project
	}{
		{
			name:     "FileTypeNotAcceptedFail",
			resource: "resource.txt",
			want:     HTTPErrorFileTypeNotAccepted,
			prjConf: &config.Project{
				AcceptTypeFiles: []string{types.TypePNG},
				Endpoints: []config.Endpoint{
					{
						Regex:             "",
						DefaultResizeOpts: types.ResizeOption{},
						CompiledRegex:     nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/%s", tt.resource), nil)
			req.Host = "127.0.0.1"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath(fmt.Sprintf("/%s", tt.resource))

			fn := GetMedia(ctx, tt.prjConf)

			err := fn(c)
			assert.Equal(t, tt.want, err)
		})
	}
}

func Test_findMatch(t *testing.T) {

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
			projectCfg: &config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			path:       "/media/image.png",
			want:       &types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png"},
			wantErr:    assert.NoError,
		},
		{
			name:       "successWithRegex",
			endpoint:   &config.Endpoint{Regex: "(?<source>.*)"},
			projectCfg: &config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			path:       "/media/image.png",
			want:       &types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png"},
			wantErr:    assert.NoError,
		},
		{
			name:       "successWithRegexAndSizeOpts",
			endpoint:   &config.Endpoint{Regex: "\\/(?<width>[0-9]{1,4})\\/(?<height>[0-9]{1,4})(?<source>.*)"},
			projectCfg: &config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			path:       "/500/500/media/image.png",
			want:       &types.ResizeOption{OriginFormat: types.TypePNG, Width: 500, Height: 500, Source: "/media/image.png"},
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

			got, err := findMatch(tt.endpoint, tt.projectCfg, tt.path)
			if !tt.wantErr(t, err, fmt.Sprintf("findMatch(%v, %v, %v)", tt.endpoint, tt.projectCfg, tt.path)) {
				return
			}
			assert.Equalf(t, tt.want, got, "findMatch(%v, %v, %v)", tt.endpoint, tt.projectCfg, tt.path)
		})
	}
}
