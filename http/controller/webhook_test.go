package controller

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWebhook(t *testing.T) {
	ctx := context.TestContext(nil)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	tests := []struct {
		name    string
		headers map[string]string
		body    []byte
		project *config.Project
		wantFn  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:    "successNoAuth",
			headers: map[string]string{"Content-Type": "application/json"},
			body:    []byte(`[{"type": "purge", "path": "app/text.txt"}]`),
			project: &config.Project{},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, rec.Code)
			},
		},
		{
			name:    "successAuth",
			headers: map[string]string{"Content-Type": "application/json", "Authorization": "Bearer token"},
			body:    []byte(`[{"type": "purge", "path": "app/text.txt"}]`),
			project: &config.Project{WebhookToken: "token"},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, rec.Code)
			},
		},
		{
			name:    "failedUnauthorized",
			headers: map[string]string{"Content-Type": "application/json"},
			body:    []byte(`{}`),
			project: &config.Project{WebhookToken: "token"},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			},
		},
		{
			name:    "failedBindBody",
			headers: map[string]string{"Content-Type": "application/json"},
			body:    []byte(`{}`),
			project: &config.Project{},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name:    "failedValidateEventsEmpty",
			headers: map[string]string{"Content-Type": "application/json"},
			body:    []byte(`[]`),
			project: &config.Project{},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			name:    "failedValidateEventMissingField",
			headers: map[string]string{"Content-Type": "application/json"},
			body:    []byte(`[{"type": "purge"}]`),
			project: &config.Project{},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chanEvents := make(chan types.Events, 2024)
			req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(tt.body))
			for key, val := range tt.headers {
				req.Header.Set(key, val)
			}
			req.Host = "127.0.0.1"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/webhook")
			err := GetWebhook(ctx, chanEvents, tt.project)(c)
			assert.NoError(t, err)
			tt.wantFn(t, rec)
		})
	}
}
