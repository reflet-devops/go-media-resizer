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

func Test_cloudflareTag_Type(t *testing.T) {
	purgeCache := &cloudflareTag{}
	assert.Equal(t, CloudflareTagKey, purgeCache.Type())
}

func Test_createCloudflareTagPurgeCache(t *testing.T) {
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
			name: "SuccessWithAuthEmail",
			cfg: config.PurgeCacheConfig{
				Type: CloudflareTagKey,
				Config: map[string]interface{}{
					"zone_id":    "zone_id",
					"auth_email": "example@example.com",
					"auth_key":   "secret",
				},
			},
			want: &cloudflareTag{ctx: ctx, projectCfg: projectCfg, cfg: ConfigCloudflare{ZoneId: "zone_id", AuthEmail: "example@example.com", AuthKey: "secret"}},
		},
		{
			name: "SuccessWithAuthToken",
			cfg: config.PurgeCacheConfig{
				Type: CloudflareTagKey,
				Config: map[string]interface{}{
					"zone_id":    "zone_id",
					"auth_token": "secret_token",
				},
			},
			want: &cloudflareTag{ctx: ctx, projectCfg: projectCfg, cfg: ConfigCloudflare{ZoneId: "zone_id", AuthToken: "secret_token"}},
		},
		{
			name: "FailDecodeCfg",
			cfg: config.PurgeCacheConfig{
				Type: CloudflareTagKey,
				Config: map[string]interface{}{
					"zone_id": []string{"zone_id"},
				},
			},
			wantErr:     true,
			errContains: "zone_id' expected type 'string', got unconvertible type '[]string'",
		},
		{
			name: "FailValidate",
			cfg: config.PurgeCacheConfig{
				Type: CloudflareTagKey,
				Config: map[string]interface{}{
					"zone_id": "zone_id",
				},
			},
			wantErr:     true,
			errContains: "Error:Field validation for 'AuthToken' failed on the 'required_without' tag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createCloudflareTagPurgeCache(ctx, projectCfg, tt.cfg)

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

func Test_cloudflareTag_Purge(t *testing.T) {

	tests := []struct {
		name       string
		cfg        ConfigCloudflare
		projectCfg *config.Project
		mockFn     func(mockClient *mockTypes.MockClient)
		events     types.Events
	}{
		{
			name:       "Success",
			cfg:        ConfigCloudflare{ZoneId: "zone_id", AuthToken: "secret_token"},
			projectCfg: &config.Project{Hostname: "example.com", PrefixPath: ""},
			events:     types.Events{{Type: types.EventTypePurge, Path: "test/text.txt"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.Body(), []byte(`{"tags":["source_path_hash_72860d7a91b9aa7fb0377f2b4c823c4f"]}`)) {
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
			name:       "SuccessWithPrefix",
			cfg:        ConfigCloudflare{ZoneId: "zone_id", AuthToken: "secret_token"},
			projectCfg: &config.Project{Hostname: "example.com", PrefixPath: "prefix"},
			events:     types.Events{{Type: types.EventTypePurge, Path: "test/text.txt"}},
			mockFn: func(mockClient *mockTypes.MockClient) {
				mockClient.EXPECT().DoTimeout(gomock.Cond(func(req *fasthttp.Request) bool {
					if !bytes.Equal(req.Body(), []byte(`{"tags":["source_path_hash_72860d7a91b9aa7fb0377f2b4c823c4f"]}`)) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			v := cloudflareTag{
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
