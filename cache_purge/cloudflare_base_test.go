package cache_purge

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
	"strings"
	"testing"
	"time"
)

func TestCloudflareDoRequest(t *testing.T) {
	tests := []struct {
		name   string
		cfg    ConfigCloudflare
		opts   CloudflareCachePurge
		mockFn func(mockClient *mockTypes.MockClient)
	}{
		{
			name: "SuccessWithAuthToken",
			cfg:  ConfigCloudflare{ZoneId: "zone_id", AuthToken: "auth_token"},
			opts: CloudflareCachePurge{Files: []string{"http://example.com/path"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.RequestURI(), []byte("https://api.cloudflare.com/client/v4/zones/zone_id/purge_cache")) {
						return false
					}
					if !bytes.Equal(req.Header.Method(), []byte("POST")) {
						return false
					}
					if !strings.Contains(req.Header.String(), fmt.Sprintf("%s: Bearer auth_token", fasthttp.HeaderAuthorization)) {
						return false
					}
					if !bytes.Equal(req.Body(), []byte(`{"files":["http://example.com/path"]}`)) {
						return false
					}
					return true
				}), gomock.Any(), gomock.Any()).DoAndReturn(
					func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
						respFn.SetStatusCode(fasthttp.StatusOK)
						respFn.SetBody([]byte(`{"success":true}`))
						return nil
					},
				)
			},
		},
		{
			name: "SuccessWithAuthEmail",
			cfg:  ConfigCloudflare{ZoneId: "zone_id", AuthEmail: "email", AuthKey: "key"},
			opts: CloudflareCachePurge{Files: []string{"http://example.com/path"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.RequestURI(), []byte("https://api.cloudflare.com/client/v4/zones/zone_id/purge_cache")) {
						return false
					}
					if !bytes.Equal(req.Header.Method(), []byte("POST")) {
						return false
					}
					if !strings.Contains(req.Header.String(), "X-Auth-Email: email") {
						return false
					}
					if !strings.Contains(req.Header.String(), "X-Auth-Key: key") {
						return false
					}
					if !bytes.Equal(req.Body(), []byte(`{"files":["http://example.com/path"]}`)) {
						return false
					}
					return true
				}), gomock.Any(), gomock.Any()).DoAndReturn(
					func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
						respFn.SetStatusCode(fasthttp.StatusOK)
						respFn.SetBody([]byte(`{"success":true}`))
						return nil
					},
				)
			},
		},
		{
			name: "SuccessWithTags",
			cfg:  ConfigCloudflare{ZoneId: "zone_id", AuthToken: "auth_token"},
			opts: CloudflareCachePurge{Tags: []string{"tag1", "tag2"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.RequestURI(), []byte("https://api.cloudflare.com/client/v4/zones/zone_id/purge_cache")) {
						return false
					}
					if !bytes.Equal(req.Header.Method(), []byte("POST")) {
						return false
					}
					if !strings.Contains(req.Header.String(), fmt.Sprintf("%s: Bearer auth_token", fasthttp.HeaderAuthorization)) {
						return false
					}
					if !bytes.Equal(req.Body(), []byte(`{"tags":["tag1","tag2"]}`)) {
						return false
					}
					return true
				}), gomock.Any(), gomock.Any()).DoAndReturn(
					func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
						respFn.SetStatusCode(fasthttp.StatusOK)
						respFn.SetBody([]byte(`{"success":true}`))
						return nil
					},
				)
			},
		},
		{
			name: "FailDoRequest",
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			name: "FailResponseCode500",
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
			CloudflareDoRequest(ctx, tt.cfg, tt.opts)
		})
	}
}
