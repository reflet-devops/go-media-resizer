package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_GetMediaCGI_AcceptedType_Fail(t *testing.T) {

	resource := "resource.txt"
	want := HTTPErrorFileTypeNotAccepted
	acceptedFileTypes := []string{types.TypePNG}

	ctx := context.TestContext(nil)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	ctx.Config.AcceptTypeFiles = acceptedFileTypes

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1/%s", resource), nil)
	req.Host = "127.0.0.1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(fmt.Sprintf("/%s", resource))

	fn := GetMediaCGI(ctx)

	err := fn(c)
	assert.Equal(t, want, err)

}

func Test_GetMediaCGI_Success(t *testing.T) {
	ctx := context.TestContext(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx.Config.AcceptTypeFiles = []string{types.TypePNG}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/images.png", nil)
	req.Host = "127.0.0.1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(fmt.Sprintf("/images.png"))
	c.SetParamNames("source")
	c.SetParamValues("https://test.test/images.png")

	mockClient := mockTypes.NewMockClient(ctrl)

	ctx.HttpClient = mockClient

	mockClient.EXPECT().DoTimeout(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
			respFn.SetStatusCode(fasthttp.StatusOK)
			respFn.SetBody([]byte("hello world"))
			return nil
		},
	)

	fn := GetMediaCGI(ctx)

	err := fn(c)
	assert.Nil(t, err)
	header := rec.Result().Header
	body := rec.Body.String()
	assert.Equal(t, "image/png", header.Get("Content-Type"))
	assert.Equal(t, body, "hello world")
}

func Test_GetMediaCGI_Decode_Error(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.AcceptTypeFiles = []string{types.TypePNG}
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/images.png", nil)
	req.Host = "127.0.0.1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(fmt.Sprintf("/images.png"))
	c.SetParamNames("source", "options")
	c.SetParamValues("https://test.test/images.png", "height={,")

	fn := GetMediaCGI(ctx)

	err := fn(c)
	assert.Nil(t, err)
	body := rec.Body.String()
	assert.Contains(t, body, "decoding failed due to the following error(s):")
}

func Test_GetMediaCGI_fetchCGIResource_Fail(t *testing.T) {
	ctx := context.TestContext(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx.Config.AcceptTypeFiles = []string{types.TypePNG}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/images.png", nil)
	req.Host = "127.0.0.1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(fmt.Sprintf("/images.png"))
	c.SetParamNames("source")
	c.SetParamValues("https://test.test/images.png")

	mockClient := mockTypes.NewMockClient(ctrl)

	ctx.HttpClient = mockClient

	mockClient.EXPECT().DoTimeout(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("test error"))

	fn := GetMediaCGI(ctx)

	err := fn(c)
	assert.Nil(t, err)
	assert.Equal(t, fasthttp.StatusInternalServerError, rec.Code)

	body := rec.Body.String()
	assert.Contains(t, body, "fetchCGIResource: GET https://test.test/images.png: error with request: test error")
}

func Test_fetchCGIResource_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	timeOut := time.Millisecond * 100
	source := "http://image.com/image.png"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockTypes.NewMockClient(ctrl)

	ctx.Config.RequestTimeout = timeOut
	ctx.HttpClient = mockClient

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(source)

	mockClient.EXPECT().DoTimeout(gomock.Eq(req), gomock.Eq(resp), gomock.Eq(timeOut)).DoAndReturn(
		func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
			respFn.SetStatusCode(fasthttp.StatusOK)
			respFn.SetBody([]byte("hello world"))
			return nil
		},
	)
	body, err := fetchCGIResource(ctx, source)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(body))
}

func Test_fetchCGIResource_Error_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	timeOut := time.Millisecond * 100
	source := "http://image.com/image.png"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockTypes.NewMockClient(ctrl)

	ctx.Config.RequestTimeout = timeOut
	ctx.HttpClient = mockClient

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(source)

	mockClient.EXPECT().DoTimeout(gomock.Eq(req), gomock.Eq(resp), gomock.Eq(timeOut)).Return(fmt.Errorf("test error"))
	_, err := fetchCGIResource(ctx, source)
	assert.Error(t, err)
	assert.Equal(t, "fetchCGIResource: GET http://image.com/image.png: error with request: test error", err.Error())
}

func Test_fetchCGIResource_StatusCode_Error(t *testing.T) {
	ctx := context.TestContext(nil)
	timeOut := time.Millisecond * 100
	source := "http://image.com/image.png"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockTypes.NewMockClient(ctrl)

	ctx.Config.RequestTimeout = timeOut
	ctx.HttpClient = mockClient

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(source)

	mockClient.EXPECT().DoTimeout(gomock.Eq(req), gomock.Eq(resp), gomock.Eq(timeOut)).DoAndReturn(
		func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
			respFn.SetStatusCode(fasthttp.StatusForbidden)
			return nil
		},
	)
	_, err := fetchCGIResource(ctx, source)
	assert.Error(t, err)
	assert.Equal(t, "fetchCGIResource: GET http://image.com/image.png: invalid status code status code: 403", err.Error())
}
