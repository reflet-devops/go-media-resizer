package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http/route"
	"github.com/stretchr/testify/assert"
)

func Test_ConfigurePrometheusMiddleware(t *testing.T) {
	username := "foo"
	password := "bar"

	tests := []struct {
		name     string
		username string
		password string
		checkFn  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:     "successValidBasicAuth",
			username: username,
			password: password,
			checkFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.NotEmpty(t, rec.Body.String())
			},
		},
		{
			name:     "failedInvalidBasicAuthUsername",
			username: "wrong",
			password: password,
			checkFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
				assert.Equal(t, "basic realm=Restricted", rec.Header().Get("WWW-Authenticate"))
			},
		},
		{
			name:     "failedInvalidBasicAuthPassword",
			username: username,
			password: "wrong",
			checkFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
				assert.Equal(t, "basic realm=Restricted", rec.Header().Get("WWW-Authenticate"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			ctx.Config.HTTP.Metrics.Enable = true
			ctx.Config.HTTP.Metrics.BasicAuth.Username = username
			ctx.Config.HTTP.Metrics.BasicAuth.Password = password

			e := echo.New()
			e.HideBanner = true
			e.HidePort = true
			req := httptest.NewRequest(http.MethodGet, route.MetricsRoute, nil)
			req.Host = "127.0.0.1"
			req.SetBasicAuth(tt.username, tt.password)
			rec := httptest.NewRecorder()

			ConfigurePrometheusMiddleware(ctx, e)
			e.ServeHTTP(rec, req)
			tt.checkFn(t, rec)
		})
	}
}
func Test_basicAuthMidDefaultFunc(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/"), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := echo.HandlerFunc(func(c echo.Context) error { return nil })
	assert.Nil(t, basicAuthMidDefaultFunc(handler)(c))
}

func Test_prometheusSkipperFunc(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/"), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.Equalf(t, false, prometheusSkipperFunc(c), "normal case")

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1%s", route.HealthCheckPingRoute), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.Equalf(t, true, prometheusSkipperFunc(c), "healh check case")
}
