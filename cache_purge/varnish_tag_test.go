package cache_purge

import (
	"bytes"
	"testing"
	"time"

	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

func Test_varnishTag_Type(t *testing.T) {
	purgeCache := &varnishTag{}
	assert.Equal(t, VarnishTagKey, purgeCache.Type())
}

func Test_createVarnishTagPurgeCache(t *testing.T) {
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
				Type: VarnishTagKey,
				Config: map[string]interface{}{
					"server": "127.0.0.1",
				},
			},
			want: &varnishTag{ctx: ctx, projectCfg: projectCfg, cfg: ConfigVarnish{Server: "127.0.0.1"}},
		},
		{
			name: "FailDecodeCfg",
			cfg: config.PurgeCacheConfig{
				Type: VarnishTagKey,
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
				Type:   VarnishTagKey,
				Config: map[string]interface{}{},
			},
			wantErr:     true,
			errContains: "Error:Field validation for 'Server' failed on the 'required' tag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createVarnishTagPurgeCache(ctx, projectCfg, tt.cfg)

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

func Test_varnishTag_Purge(t *testing.T) {

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
			projectCfg: &config.Project{ID: "project-id", PrefixPath: ""},
			events:     types.Events{{Type: types.EventTypePurge, Path: "test/text.txt"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.RequestURI(), []byte("http://127.0.0.1/test/text.txt")) {
						return false
					}
					if !bytes.Equal(req.Header.Method(), []byte("BAN")) {
						return false
					}

					if string(req.Header.Peek(types.HeaderCachePurge)) != "(source_path_hash_c6047f34708c36213118502d98d05466)" {
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
			projectCfg: &config.Project{ID: "project-id", PrefixPath: "prefix"},
			events:     types.Events{{Type: types.EventTypePurge, Path: "test/text.txt"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.RequestURI(), []byte("http://127.0.0.1/prefix/test/text.txt")) {
						return false
					}
					if !bytes.Equal(req.Header.Method(), []byte("BAN")) {
						return false
					}

					if string(req.Header.Peek(types.HeaderCachePurge)) != "(source_path_hash_b7c209b0abcb1d2bc96952a8a512ee12)" {
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
			v := varnishTag{
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
