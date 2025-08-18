package controller

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type errorReader struct {
	r     io.Reader
	limit int
	count int
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	if e.count >= e.limit {
		return 0, fmt.Errorf("simulated read error")
	}
	n, err = e.r.Read(p)
	e.count += n
	return
}

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
	ctx := context.TestContext(nil)

	tests := []struct {
		name         string
		opts         *types.ResizeOption
		headerAccept string
		contentFn    func() io.Reader

		wantErr assert.ErrorAssertionFunc
		wantFn  func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:         "successWithPlainText",
			opts:         &types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypeText, Source: "/text.txt"},
			headerAccept: "text/plain",
			contentFn: func() io.Reader {
				return bytes.NewReader([]byte("hello"))
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, types.MimeTypeText, rec.Header().Get(echo.HeaderContentType))
				assert.Equal(t, []byte("hello"), rec.Body.Bytes())
			},
			wantErr: assert.NoError,
		},
		{
			name:         "successWithResizeFormat",
			opts:         &types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypePNG, Source: "/paysage.png", Width: 500},
			headerAccept: "image/avif,image/webp,image/png",
			contentFn: func() io.Reader {
				file, errOpen := os.Open("../../fixtures/paysage.png")
				assert.NoError(t, errOpen)
				return file
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, types.MimeTypeAVIF, rec.Header().Get(echo.HeaderContentType))
				assert.NotEmpty(t, rec.Body.Bytes())
			},
			wantErr: assert.NoError,
		},
		{
			name:         "failedResize",
			opts:         &types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypePNG, Source: "/paysage.png", Width: 500},
			headerAccept: "image/avif,image/webp,image/png",
			contentFn: func() io.Reader {
				return &errorReader{r: bytes.NewBufferString("")}
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "failed to resize /paysage.png", rec.Body.String())
			},
			wantErr: assert.NoError,
		},
		{
			name:         "failedFormat",
			opts:         &types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypePNG, Source: "/paysage.png"},
			headerAccept: "image/avif,image/webp,image/png",
			contentFn: func() io.Reader {
				return &errorReader{r: bytes.NewBufferString("")}
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "failed to format /paysage.png", rec.Body.String())
			},
			wantErr: assert.NoError,
		},
		{
			name:         "failedReadData",
			opts:         &types.ResizeOption{Format: types.TypePNG, OriginFormat: types.TypePNG, Source: "/paysage.png"},
			headerAccept: "image/avif,image/webp,image/png",
			contentFn: func() io.Reader {
				return &errorReader{r: bytes.NewBufferString("")}
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "failed to read data /paysage.png", rec.Body.String())
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
