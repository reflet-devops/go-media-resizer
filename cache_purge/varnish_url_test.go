package cache_purge

import (
	"bytes"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_varnishUrl_Type(t *testing.T) {
	purgeCache := &varnishUrl{}
	assert.Equal(t, VarnishUrlKey, purgeCache.Type())
}

func Test_createVarnishUrlPurgeCache(t *testing.T) {
	ctx := context.TestContext(nil)
	projectCfg := &config.Project{}
	tests := []struct {
		name        string
		cfg         config.PurgeCacheConfig
		want        types.PurgeCache
		wantErr     bool
		errContains string
	}{
		{
			name: "Success",
			cfg: config.PurgeCacheConfig{
				Type: VarnishUrlKey,
				Config: map[string]interface{}{
					"server": "127.0.0.1",
				},
			},
			want: &varnishUrl{ctx: ctx, projectCfg: projectCfg, cfg: ConfigVarnish{Server: "127.0.0.1"}},
		},
		{
			name: "FailDecodeCfg",
			cfg: config.PurgeCacheConfig{
				Type: VarnishUrlKey,
				Config: map[string]interface{}{
					"server": []string{"127.0.0.1"},
				},
			},
			wantErr:     true,
			errContains: "server' expected type 'string', got unconvertible type '[]string'",
		},
		{
			name: "FailValidate",
			cfg: config.PurgeCacheConfig{
				Type:   VarnishUrlKey,
				Config: map[string]interface{}{},
			},
			wantErr:     true,
			errContains: "Error:Field validation for 'Server' failed on the 'required' tag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createVarnishUrlPurgeCache(ctx, projectCfg, tt.cfg)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_varnishUrl_Purge(t *testing.T) {

	tests := []struct {
		name       string
		cfg        ConfigVarnish
		projectCfg *config.Project
		mockFn     func(mockClient *mockTypes.MockClient)
		events     types.Events
	}{
		{
			name:       "Success",
			cfg:        ConfigVarnish{Server: "http://127.0.0.1"},
			projectCfg: &config.Project{PrefixPath: ""},
			events:     types.Events{{Type: types.EventTypePurge, Path: "test/text.txt"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.RequestURI(), []byte("http://127.0.0.1/test/text.txt")) {
						return false
					}
					if !bytes.Equal(req.Header.Method(), []byte("PURGE")) {
						return false
					}
					return true
				}), gomock.Any(), gomock.Any()).DoAndReturn(
					func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
						respFn.SetStatusCode(fasthttp.StatusOK)
						respFn.SetBody([]byte("hello world"))
						return nil
					},
				)
			},
		},
		{
			name:       "SuccessWithPrefix",
			cfg:        ConfigVarnish{Server: "http://127.0.0.1"},
			projectCfg: &config.Project{PrefixPath: "prefix"},
			events:     types.Events{{Type: types.EventTypePurge, Path: "test/text.txt"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.RequestURI(), []byte("http://127.0.0.1/prefix/test/text.txt")) {
						return false
					}
					if !bytes.Equal(req.Header.Method(), []byte("PURGE")) {
						return false
					}
					return true
				}), gomock.Any(), gomock.Any()).DoAndReturn(
					func(req *fasthttp.Request, respFn *fasthttp.Response, timeout time.Duration) error {
						respFn.SetStatusCode(fasthttp.StatusOK)
						respFn.SetBody([]byte("hello world"))
						return nil
					},
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			v := varnishUrl{
				ctx:        ctx,
				cfg:        tt.cfg,
				projectCfg: tt.projectCfg,
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockClient := mockTypes.NewMockClient(ctrl)
			tt.mockFn(mockClient)
			ctx.HttpClient = mockClient
			v.Purge(tt.events)
		})
	}
}
