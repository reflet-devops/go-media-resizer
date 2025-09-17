package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/http/route"
	"github.com/stretchr/testify/assert"
)

func TestConfigureRequestIdMiddleware(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	ConfigureRequestIdMiddleware(e)
}

func Test_requestIdSkipperFunc(t *testing.T) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/"), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.Equalf(t, false, requestIdSkipperFunc(c), "normal case")

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1%s", route.HealthCheckPingRoute), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.Equalf(t, true, requestIdSkipperFunc(c), "healh check case")

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1"), nil)
	req.Header.Add(echo.HeaderXRequestID, "1234567890")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.Equalf(t, true, requestIdSkipperFunc(c), "request id header present case")

}
