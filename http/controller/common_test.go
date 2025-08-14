package controller

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetectFormatFromHeaderAccept(t *testing.T) {

	tests := []struct {
		name              string
		acceptHeaderValue string
		opts              *types.ResizeOption
		want              *types.ResizeOption
	}{
		{
			name:              "detectFormatAvifWithAutoAndGoodAcceptHeader",
			acceptHeaderValue: "image/avif,image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
		},
		{
			name:              "detectFormatAvifWithAvifAndGoodAcceptHeader",
			acceptHeaderValue: "image/avif,image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
		},
		{
			name:              "detectFormatAvifWithAvifAndWrongAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatAvifWithAutoAndWrongAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatWebWithAutoAndGoodAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatWebWithWebpAndGoodAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatWebWithWebpAndWrongAcceptHeader",
			acceptHeaderValue: "image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
		{
			name:              "detectFormatWebWithAutoAndWrongAcceptHeader",
			acceptHeaderValue: "image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
		{
			name:              "detectFormatPngWithAuto",
			acceptHeaderValue: "image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DetectFormatFromHeaderAccept(tt.acceptHeaderValue, tt.opts)
			assert.Equal(t, tt.want, tt.opts)
		})
	}
}

func TestSendStream(t *testing.T) {
	e := echo.New()
	opts := &types.ResizeOption{Format: types.TypePNG}
	content := []byte("hello")
	buffer := bytes.NewBuffer(content)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/"), nil)
	req.Host = "127.0.0.1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := SendStream(c, opts, buffer)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, types.MimeTypePNG, rec.Header().Get(echo.HeaderContentType))
	assert.Equal(t, content, rec.Body.Bytes())
}
