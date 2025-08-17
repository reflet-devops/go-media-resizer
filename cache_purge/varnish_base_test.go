package cache_purge

import (
	"errors"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestVarnishDoRequest(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		uri     string
		headers map[string]string
		mockFn  func(mockClient *mockTypes.MockClient)
	}{
		{
			name:    "Success",
			method:  "PURGE",
			uri:     "http://varnish/image.png",
			headers: map[string]string{"test": "test"},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
						respFn.SetStatusCode(fasthttp.StatusOK)
						respFn.SetBody([]byte("hello world"))
						return nil
					},
				)
			},
		},
		{
			name:    "FailDoRequest",
			method:  "PURGE",
			uri:     "http://varnish/image.png",
			headers: map[string]string{},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			name:    "FailResponseCode500",
			method:  "PURGE",
			uri:     "http://varnish/image.png",
			headers: map[string]string{},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
						respFn.SetStatusCode(fasthttp.StatusInternalServerError)
						respFn.SetBody([]byte("error"))
						return nil
					},
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockClient := mockTypes.NewMockClient(ctrl)
			ctx.HttpClient = mockClient
			tt.mockFn(mockClient)
			VarnishDoRequest(ctx, tt.method, tt.uri, tt.headers)
		})
	}
}
