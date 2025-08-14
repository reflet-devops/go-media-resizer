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
	"testing"
)

func Test_GetMedia(t *testing.T) {
	ctx := context.TestContext(nil)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	tests := []struct {
		name     string
		resource string
		wantFn   func(t *testing.T, rec *httptest.ResponseRecorder)
		prjConf  *config.Project
	}{
		{
			name:     "FileTypeNotAcceptedFail",
			resource: "resource.txt",
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, []byte("file type not accepted"), rec.Body.Bytes())
			},
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

			err := GetMedia(ctx, tt.prjConf)(c)
			assert.NoError(t, err)
			tt.wantFn(t, rec)
		})
	}
}
