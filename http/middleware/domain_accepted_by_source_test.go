package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDomainAcceptedBySource_Validate(t *testing.T) {
	tests := []struct {
		name              string
		source            string
		server            string
		acceptedHostnames []string
		allowSelfDomain   bool
		want              bool
	}{
		{
			name:              "AllowSelfDomainTrue",
			source:            "https://self/my/resource",
			server:            "https://self",
			acceptedHostnames: []string{},
			allowSelfDomain:   true,
			want:              true,
		},
		{
			name:              "AllowSelfDomainFalse",
			source:            "https://not.self/my/resource",
			server:            "https://self",
			acceptedHostnames: []string{},
			allowSelfDomain:   true,
			want:              false,
		},
		{
			name:              "AllowDomainTrue",
			source:            "https://accepted/my/resource",
			server:            "https://self",
			acceptedHostnames: []string{"accepted"},
			allowSelfDomain:   false,
			want:              true,
		},
		{
			name:              "AllowDomainFalse",
			source:            "https://refused/my/resource",
			server:            "https://self",
			acceptedHostnames: []string{"accepted"},
			allowSelfDomain:   false,
			want:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			ctx.Config.ResizeCGI.AllowSelfDomain = tt.allowSelfDomain
			ctx.Config.ResizeCGI.AllowDomains = tt.acceptedHostnames
			validator := NewDomainAcceptedBySource(ctx)

			got := validator.Validate(tt.source, tt.server)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDomainAcceptedBySource_Handler(t *testing.T) {
	ctx := context.TestContext(nil)
	domainAcceptedBySource := NewDomainAcceptedBySource(ctx)
	e := echo.New()

	tests := []struct {
		name              string
		source            string
		want              error
		allowSelfDomain   bool
		acceptedHostnames []string
	}{
		{
			name:              "AllowSelfDomainSuccess",
			source:            "https://127.0.0.1/resource",
			want:              nil,
			allowSelfDomain:   true,
			acceptedHostnames: []string{},
		},
		{
			name:              "AllowSelfDomainFail",
			source:            "https://127.0.0.2/resource",
			want:              &echo.HTTPError{Code: http.StatusForbidden, Message: "domain not allowed: https://127.0.0.2/resource"},
			allowSelfDomain:   true,
			acceptedHostnames: []string{},
		},
		{
			name:              "AllowSelfDomainEmptySourceFail",
			source:            "",
			want:              &echo.HTTPError{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)},
			allowSelfDomain:   true,
			acceptedHostnames: []string{},
		},
		{
			name:              "AllowDomainSuccess",
			source:            "https://127.0.0.1/resource",
			acceptedHostnames: []string{"127.0.0.1"},
			allowSelfDomain:   false,
			want:              nil,
		},
		{
			name:              "AllowDomainFail",
			source:            "https://127.0.0.2/resource",
			acceptedHostnames: []string{"127.0.0.1"},
			allowSelfDomain:   false,
			want:              &echo.HTTPError{Code: http.StatusForbidden, Message: "domain not allowed: https://127.0.0.2/resource"},
		},
		{
			name:              "AllowDomainFailEmptyList",
			source:            "https://127.0.0.2/resource",
			acceptedHostnames: []string{""},
			allowSelfDomain:   false,
			want:              &echo.HTTPError{Code: http.StatusForbidden, Message: "domain not allowed: https://127.0.0.2/resource"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.Config.ResizeCGI.AllowSelfDomain = tt.allowSelfDomain
			ctx.Config.ResizeCGI.AllowDomains = tt.acceptedHostnames

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.source), nil)
			req.Host = "127.0.0.1"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:source")
			c.SetParamNames("source")
			c.SetParamValues(tt.source)
			handler := echo.HandlerFunc(func(c echo.Context) error { return nil })
			err := domainAcceptedBySource.Handler(handler)(c)
			assert.Equal(t, tt.want, err)
		})
	}
}
