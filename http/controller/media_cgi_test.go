package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetMediaCGI(t *testing.T) {

	tests := []struct {
		name              string
		resource          string
		want              error
		acceptedFileTypes []string
	}{
		{
			name:              "FileTypeNotAcceptedFail",
			resource:          "resource.txt",
			want:              HTTPErrorFileTypeNotAccepted,
			acceptedFileTypes: []string{types.TypePNG},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/%s", tt.resource), nil)
			req.Host = "127.0.0.1"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath(fmt.Sprintf("/%s", tt.resource))

			fn := GetMediaCGI(ctx)

			err := fn(c)
			assert.Equal(t, tt.want, err)
		})
	}
}
