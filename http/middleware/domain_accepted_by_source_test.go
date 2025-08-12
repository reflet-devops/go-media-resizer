package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestDomainAcceptedBySource_Handler_AllowSelfDomain(t *testing.T) {
	ctx := context.TestContext(nil)
	acceptedHostnames := []string{}
	ctx.Config.ResizeCGI.AllowSelfDomain = true
	ctx.Config.ResizeCGI.AllowDomains = acceptedHostnames
	domainAcceptedBySource := NewDomainAcceptedBySource(ctx)
	e := echo.New()

	tests := []struct {
		name   string
		source string
		want   int
	}{
		{
			name:   "AllowSelfDomainSuccess",
			source: "https://127.0.0.1/resource",
			want:   http.StatusOK,
		},
		{
			name:   "AllowSelfDomainFail",
			source: "https://127.0.0.2/resource",
			want:   http.StatusForbidden,
		},
		{
			name:   "AllowSelfDomainEmptySourceFail",
			source: "",
			want:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tt.source, nil)
			req.Host = "127.0.0.1"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("source")
			c.SetParamValues(tt.source)
			c.SetPath("/:source")

			_ = domainAcceptedBySource.Handler(c.Handler())(c)
			assert.Equal(t, tt.want, rec.Code)
		})
	}
}

func TestDomainAcceptedBySource_Handler_AllowDomains(t *testing.T) {
	ctx := context.TestContext(nil)
	path := "/:source"
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	domainAcceptedBySource := NewDomainAcceptedBySource(ctx)

	e.GET(path, func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Use(domainAcceptedBySource.Handler)
	go func() {
		_ = e.Start("127.0.0.1:8080")
	}()
	time.Sleep(time.Millisecond * 500)

	tests := []struct {
		name              string
		source            string
		acceptedHostnames []string
		want              int
	}{
		{
			name:              "Success",
			source:            "https://127.0.0.1/resource",
			acceptedHostnames: []string{"127.0.0.1"},
			want:              http.StatusOK,
		},
		{
			name:              "Fail",
			source:            "https://127.0.0.2/resource",
			acceptedHostnames: []string{"127.0.0.1"},
			want:              http.StatusForbidden,
		},
		{
			name:              "FailEmptySource",
			source:            "https://127.0.0.2/resource",
			acceptedHostnames: []string{""},
			want:              http.StatusForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.Config.ResizeCGI.AllowSelfDomain = false
			ctx.Config.ResizeCGI.AllowDomains = tt.acceptedHostnames
			resp, _ := http.Get(fmt.Sprintf("http://%s/%s", e.Server.Addr, tt.source))
			assert.Equal(t, tt.want, resp.StatusCode)
		})
	}
	_ = e.Close()

}
