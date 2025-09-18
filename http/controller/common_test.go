package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
)

func TestDetectFormatFromHeaderAccept(t *testing.T) {

	tests := []struct {
		name                 string
		acceptHeaderValue    string
		enableFormatAutoAVIF bool
		opts                 *types.ResizeOption
		want                 *types.ResizeOption
	}{
		{
			name:                 "detectFormatAvifWithAutoAndGoodAcceptHeader",
			acceptHeaderValue:    "image/avif,image/webp,image/png",
			enableFormatAutoAVIF: true,
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
		},
		{
			name:                 "detectFormatWebpWithAutoAndAvifAcceptHeaderAndDisabledFormatAutoAVIF",
			acceptHeaderValue:    "image/avif,image/webp,image/png",
			enableFormatAutoAVIF: false,
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:                 "detectFormatAvifWithAvifAndGoodAcceptHeader",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/avif,image/webp,image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
		},
		{
			name:                 "detectFormatAvifWithAvifAndWrongAcceptHeader",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/webp,image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:                 "detectFormatAvifWithAutoAndWrongAcceptHeader",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/webp,image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:                 "detectFormatWebWithAutoAndGoodAcceptHeader",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/webp,image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:                 "detectFormatWebWithWebpAndGoodAcceptHeader",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/webp,image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:                 "detectFormatWebWithWebpAndWrongAcceptHeader",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
		{
			name:                 "detectFormatWebWithAutoAndWrongAcceptHeader",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
		{
			name:                 "detectFormatPngWithAuto",
			enableFormatAutoAVIF: true,
			acceptHeaderValue:    "image/png",
			opts:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:                 &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			ctx.Config.EnableFormatAutoAVIF = tt.enableFormatAutoAVIF
			DetectFormatFromHeaderAccept(ctx, tt.acceptHeaderValue, tt.opts)
			assert.Equal(t, tt.want, tt.opts)
		})
	}
}

func TestSendStream(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.EnableFormatAutoAVIF = true

	tests := []struct {
		name         string
		opts         *types.ResizeOption
		headerAccept string
		contentFn    func() *bytes.Buffer

		wantErr assert.ErrorAssertionFunc
		wantFn  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:         "successWithPlainText",
			opts:         &types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypeText, Source: "/text.txt", Headers: types.Headers{"X-Custom": "foo"}, Tags: []string{"tag1"}},
			headerAccept: "text/plain",
			contentFn: func() *bytes.Buffer {
				buff := ctx.BufferPool.Get().(*bytes.Buffer)
				buff.WriteString("hello")
				return buff
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, types.MimeTypeText, rec.Header().Get(echo.HeaderContentType))
				assert.Equal(t, "foo", rec.Header().Get("X-Custom"))
				assert.Equal(t, []byte("hello"), rec.Body.Bytes())
			},
			wantErr: assert.NoError,
		},
		{
			name:         "successWithResizeFormat",
			opts:         &types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypePNG, Source: "/paysage.png", Headers: types.Headers{"X-Custom": "foo"}, Width: 500},
			headerAccept: "image/avif,image/webp,image/png",
			contentFn: func() *bytes.Buffer {
				file, errOpen := os.Open("../../fixtures/paysage.png")
				assert.NoError(t, errOpen)
				buff := ctx.BufferPool.Get().(*bytes.Buffer)
				_, _ = io.Copy(buff, file)
				return buff
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, types.MimeTypeAVIF, rec.Header().Get(echo.HeaderContentType))
				assert.Equal(t, "foo", rec.Header().Get("X-Custom"))
				assert.NotEmpty(t, rec.Body.Bytes())
			},
			wantErr: assert.NoError,
		},
		{
			name:         "failedTransform",
			opts:         &types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypePNG, Source: "/paysage.png", Width: 500},
			headerAccept: "image/avif,image/webp,image/png",
			contentFn: func() *bytes.Buffer {
				buff := bytes.NewBuffer([]byte("test"))
				buff.Reset()
				return buff
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "failed to transform image /paysage.png", rec.Body.String())
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.HideBanner = true
			e.HidePort = true

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/"), nil)
			req.Host = "127.0.0.1"
			req.Header.Set(echo.HeaderAccept, tt.headerAccept)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := SendStream(ctx, c, tt.opts, tt.contentFn())

			tt.wantErr(t, err)
			tt.wantFn(t, rec)
		})
	}
}
