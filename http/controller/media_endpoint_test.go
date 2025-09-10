package controller

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
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

func (e *errorReader) Close() error {
	return nil
}

func Test_GetMedia(t *testing.T) {
	ctx := context.TestContext(nil)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	regexStr := "/wrong/(?<source>.*)"
	re, errReCompile := regexp.Compile(regexStr)
	assert.NoError(t, errReCompile)
	tests := []struct {
		name     string
		resource string
		prjConf  *config.Project
		mockFn   func(mockStorage *mockTypes.MockStorage)
		wantFn   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:     "success",
			resource: "resource.txt",
			prjConf: &config.Project{
				AcceptTypeFiles: []string{types.TypeText},
				Endpoints: []config.Endpoint{
					{
						Regex:             "",
						DefaultResizeOpts: types.ResizeOption{},
						CompiledRegex:     nil,
					},
				},
			},
			mockFn: func(mockStorage *mockTypes.MockStorage) {
				b := io.NopCloser(bytes.NewBufferString("hello world"))
				mockStorage.EXPECT().GetFile(gomock.Eq("resource.txt")).Times(1).Return(b, nil)
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "hello world", rec.Body.String())
			},
		},
		{
			name:     "success_EndpointNotMatch",
			resource: "resource.txt",
			prjConf: &config.Project{
				AcceptTypeFiles: []string{types.TypeText},
				Endpoints: []config.Endpoint{
					{
						Regex:             regexStr,
						DefaultResizeOpts: types.ResizeOption{},
						CompiledRegex:     re,
					},
				},
			},
			mockFn: func(mockStorage *mockTypes.MockStorage) {},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "file not found", rec.Body.String())
			},
		},
		{
			name:     "success_NoEndpoint",
			resource: "resource.txt",
			prjConf: &config.Project{
				AcceptTypeFiles: []string{types.TypeText},
				Endpoints:       []config.Endpoint{},
			},
			mockFn: func(mockStorage *mockTypes.MockStorage) {},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "file not found", rec.Body.String())
			},
		},
		{
			name:     "fail_fileTypeNotAcceptedFail",
			resource: "resource.txt",
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
			mockFn: func(mockStorage *mockTypes.MockStorage) {},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, []byte("file type not accepted"), rec.Body.Bytes())
			},
		},
		{
			name:     "fail_GetFile",
			resource: "resource.txt",
			prjConf: &config.Project{
				AcceptTypeFiles: []string{types.TypeText},
				Endpoints: []config.Endpoint{
					{
						Regex:             "",
						DefaultResizeOpts: types.ResizeOption{},
						CompiledRegex:     nil,
					},
				},
			},
			mockFn: func(mockStorage *mockTypes.MockStorage) {
				mockStorage.EXPECT().GetFile(gomock.Eq("resource.txt")).Times(1).Return(nil, errors.New("file not found"))
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "file not found", rec.Body.String())
			},
		},
		{
			name:     "fail_Copy",
			resource: "resource.txt",
			prjConf: &config.Project{
				AcceptTypeFiles: []string{types.TypeText},
				Endpoints: []config.Endpoint{
					{
						Regex:             "",
						DefaultResizeOpts: types.ResizeOption{},
						CompiledRegex:     nil,
					},
				},
			},
			mockFn: func(mockStorage *mockTypes.MockStorage) {
				mockStorage.EXPECT().GetFile(gomock.Eq("resource.txt")).Times(1).Return(&errorReader{r: bytes.NewBufferString("test")}, nil)
			},
			wantFn: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), types.MimeTypeText)
				assert.Equal(t, "buffer copy failed", rec.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStorage := mockTypes.NewMockStorage(ctrl)
			tt.mockFn(mockStorage)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.resource), nil)
			req.Host = "127.0.0.1"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath(fmt.Sprintf("/%s", tt.resource))

			err := GetMedia(ctx, tt.prjConf, mockStorage)(c)
			assert.NoError(t, err)
			tt.wantFn(t, rec)
		})
	}
}
