package storage

import (
	builtinCtx "context"
	"errors"
	"fmt"
	"github.com/jonboulle/clockwork"
	libMinio "github.com/minio/minio-go/v7"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	mockTypes "github.com/reflet-devops/go-media-resizer/mocks/types"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_minio_Type(t *testing.T) {
	storage := &minio{}
	assert.Equal(t, MinioKey, storage.Type())
}

func Test_minio_GetPrefix(t *testing.T) {
	storage := &minio{cfg: ConfigMinio{PrefixPath: "/app"}}
	assert.Equal(t, "/app", storage.GetPrefix())
}

func Test_minio_getClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	minioMock := mockTypes.NewMockMinioClient(ctrl)
	storage := &minio{currentClient: minioMock}
	assert.Equal(t, minioMock, storage.getClient())
}

func Test_minio_createMinioClient(t *testing.T) {
	minioClient, err := createMinioClient(ConfigClientMinio{
		Endpoint:   "localhost",
		BucketName: "foo",
		AccessKey:  "bar",
		SecretKey:  "bar",
		UseSSL:     false,
	})
	assert.NoError(t, err)
	assert.IsType(t, &libMinio.Client{}, minioClient)
}

func Test_minio_GetFile(t *testing.T) {
	ctx := context.TestContext(nil)
	tests := []struct {
		name    string
		path    string
		cfg     ConfigMinio
		mockFn  func(minioMock *mockTypes.MockMinioClient)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			path: "foo/bar.txt",
			cfg:  ConfigMinio{PrefixPath: "/app", ConfigClientMinio: ConfigClientMinio{BucketName: "bucket"}},
			mockFn: func(minioMock *mockTypes.MockMinioClient) {
				minioMock.EXPECT().GetObject(gomock.Any(), gomock.Eq("bucket"), gomock.Eq("/app/foo/bar.txt"), gomock.Any()).Times(1).Return(nil, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "SuccessWithSlash",
			path: "/foo/bar.txt",
			cfg:  ConfigMinio{PrefixPath: "/app", ConfigClientMinio: ConfigClientMinio{BucketName: "bucket"}},
			mockFn: func(minioMock *mockTypes.MockMinioClient) {
				minioMock.EXPECT().GetObject(gomock.Any(), gomock.Eq("bucket"), gomock.Eq("/app/foo/bar.txt"), gomock.Any()).Times(1).Return(nil, nil)
			},
			wantErr: assert.NoError,
		},
		{
			name: "FailGetObject",
			path: "foo/bar.txt",
			cfg:  ConfigMinio{PrefixPath: "/app", ConfigClientMinio: ConfigClientMinio{BucketName: "bucket"}},
			mockFn: func(minioMock *mockTypes.MockMinioClient) {
				minioMock.EXPECT().GetObject(gomock.Any(), gomock.Eq("bucket"), gomock.Eq("/app/foo/bar.txt"), gomock.Any()).Times(1).Return(nil, errors.New("GetObject failed"))
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			minioMock := mockTypes.NewMockMinioClient(ctrl)
			defer ctrl.Finish()
			tt.mockFn(minioMock)
			m := &minio{
				currentClient: minioMock,
				primaryClient: minioMock,
				cfg:           tt.cfg,
				ctx:           ctx,
			}
			_, err := m.GetFile(tt.path)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFile(%v)", tt.path)) {
				return
			}
		})
	}
}

func Test_createMinioStorage(t *testing.T) {
	ctx := context.TestContext(nil)

	tests := []struct {
		name        string
		cfg         config.StorageConfig
		want        types.Storage
		wantErr     bool
		errContains string
	}{
		{
			name: "Success",
			cfg: config.StorageConfig{
				Type: MinioKey,
				Config: map[string]interface{}{
					"prefix_path":           "/app",
					"health_check_interval": "2s",
					"endpoint":              "localhost",
					"bucket":                "bucket",
					"access_key":            "access",
					"secret_key":            "secret",
					"use_ssl":               true,
				},
			},
			want: &minio{ctx: ctx, clock: clockwork.NewRealClock(), cfg: ConfigMinio{
				ConfigClientMinio: ConfigClientMinio{
					Endpoint:   "localhost",
					BucketName: "bucket",
					AccessKey:  "access",
					SecretKey:  "secret",
					UseSSL:     true,
				},
				HealthCheckInterval: time.Second * 2,
				PrefixPath:          "/app",
			}},
		},
		{
			name: "SuccessTrimPrefixPath",
			cfg: config.StorageConfig{
				Type: MinioKey,
				Config: map[string]interface{}{
					"prefix_path": "/app/",
					"endpoint":    "localhost",
					"bucket":      "bucket",
					"access_key":  "access",
					"secret_key":  "secret",
					"use_ssl":     true,
				},
			},
			want: &minio{ctx: ctx, clock: clockwork.NewRealClock(), cfg: ConfigMinio{
				ConfigClientMinio: ConfigClientMinio{
					Endpoint:   "localhost",
					BucketName: "bucket",
					AccessKey:  "access",
					SecretKey:  "secret",
					UseSSL:     true,
				},
				HealthCheckInterval: time.Second * 5,
				PrefixPath:          "/app",
			}},
		},
		{
			name: "SuccessWithFallbackMinio",
			cfg: config.StorageConfig{
				Type: MinioKey,
				Config: map[string]interface{}{
					"prefix_path": "/app",
					"endpoint":    "localhost",
					"bucket":      "bucket",
					"access_key":  "access",
					"secret_key":  "secret",
					"use_ssl":     true,
					"fallback": map[string]interface{}{
						"endpoint":   "fallback",
						"bucket":     "bucket",
						"access_key": "access",
						"secret_key": "secret",
						"use_ssl":    true,
					},
				},
			},
			want: &minio{ctx: ctx, clock: clockwork.NewRealClock(), cfg: ConfigMinio{
				ConfigClientMinio: ConfigClientMinio{
					Endpoint:   "localhost",
					BucketName: "bucket",
					AccessKey:  "access",
					SecretKey:  "secret",
					UseSSL:     true,
				},
				FallbackMinio: &ConfigClientMinio{
					Endpoint:   "fallback",
					BucketName: "bucket",
					AccessKey:  "access",
					SecretKey:  "secret",
					UseSSL:     true,
				},
				HealthCheckInterval: time.Second * 5,
				PrefixPath:          "/app",
			}},
		},
		{
			name: "FailDecodeCfg",
			cfg: config.StorageConfig{
				Type: MinioKey,
				Config: map[string]interface{}{
					"prefix_path": []string{},
				},
			},
			wantErr:     true,
			errContains: "prefix_path' expected type 'string', got unconvertible type '[]string",
		},
		{
			name: "FailValidateCfg",
			cfg: config.StorageConfig{
				Type: MinioKey,
				Config: map[string]interface{}{
					"prefix_path": "/app",
				},
			},
			wantErr:     true,
			errContains: "Error:Field validation for 'Endpoint' failed on the 'required' tag",
		},
		{
			name: "FailCreateMinioClient",
			cfg: config.StorageConfig{
				Type: MinioKey,
				Config: map[string]interface{}{
					"prefix_path": "/app",
					"endpoint":    "localhost:9000:9",
					"bucket":      "bucket",
					"access_key":  "access",
					"secret_key":  "secret",
					"use_ssl":     true,
				},
			},
			wantErr:     true,
			errContains: "Endpoint: localhost:9000:9 does not follow ip address or domain name standards.",
		},
		{
			name: "FailedWithFallbackMinio",
			cfg: config.StorageConfig{
				Type: MinioKey,
				Config: map[string]interface{}{
					"prefix_path": "/app",
					"endpoint":    "localhost",
					"bucket":      "bucket",
					"access_key":  "access",
					"secret_key":  "secret",
					"use_ssl":     true,
					"fallback": map[string]interface{}{
						"endpoint":   "fallback:9000:9",
						"bucket":     "bucket",
						"access_key": "access",
						"secret_key": "secret",
						"use_ssl":    true,
					},
				},
			},
			wantErr:     true,
			errContains: "Endpoint: fallback:9000:9 does not follow ip address or domain name standards.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createMinioStorage(ctx, tt.cfg)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				gotMinio := got.(*minio)
				if gotMinio.secondaryClient != nil {
					ctx.Cancel()
					assert.IsType(t, &libMinio.Client{}, gotMinio.secondaryClient)
					gotMinio.secondaryClient = nil
				}
				assert.IsType(t, &libMinio.Client{}, gotMinio.currentClient)
				assert.IsType(t, &libMinio.Client{}, gotMinio.primaryClient)
				gotMinio.currentClient = nil
				gotMinio.primaryClient = nil
				assert.NoError(t, err)
				assert.Equal(t, tt.want, gotMinio)
			}
		})
	}
}

func Test_minio_startFallback(t *testing.T) {
	ctx := context.TestContext(nil)
	cfg := ConfigMinio{HealthCheckInterval: time.Millisecond * 100}
	tests := []struct {
		name   string
		mockFn func(primaryMock *mockTypes.MockMinioClient)
		testFn func(t *testing.T, m *minio, clock *clockwork.FakeClock)
	}{
		{
			name: "Success",
			mockFn: func(primaryMock *mockTypes.MockMinioClient) {
				primaryMock.EXPECT().HealthCheck(gomock.Any()).Return(nil, nil)
				primaryMock.EXPECT().IsOnline().Return(true)
			},
			testFn: func(t *testing.T, m *minio, clock *clockwork.FakeClock) {
				assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")
				_ = clock.BlockUntilContext(builtinCtx.Background(), 1)
				clock.Advance(cfg.HealthCheckInterval)
				time.Sleep(100 * time.Millisecond)
				assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")
			},
		},
		{
			name: "FailedHealthCheck",
			mockFn: func(primaryMock *mockTypes.MockMinioClient) {
				primaryMock.EXPECT().HealthCheck(gomock.Any()).Return(func() {}, errors.New("health check failed"))
				primaryMock.EXPECT().IsOnline().Return(true)
			},
			testFn: func(t *testing.T, m *minio, clock *clockwork.FakeClock) {
				assert.True(t, m.primaryClient == m.currentClient, "current Client must be primary")
				_ = clock.BlockUntilContext(builtinCtx.Background(), 1)
				clock.Advance(cfg.HealthCheckInterval)
				time.Sleep(100 * time.Millisecond)
				assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")
			},
		},
		{
			name: "SuccessSwitchToSecondary",
			mockFn: func(primaryMock *mockTypes.MockMinioClient) {
				primaryMock.EXPECT().HealthCheck(gomock.Any()).Times(1).Return(nil, nil)
				primaryMock.EXPECT().IsOnline().Times(3).Return(false)
			},
			testFn: func(t *testing.T, m *minio, clock *clockwork.FakeClock) {
				assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")

				for i := 0; i < 2; i++ {
					_ = clock.BlockUntilContext(builtinCtx.Background(), 1)
					clock.Advance(cfg.HealthCheckInterval)
					time.Sleep(100 * time.Millisecond)
					assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")
				}

				_ = clock.BlockUntilContext(builtinCtx.Background(), 1)
				clock.Advance(cfg.HealthCheckInterval)
				time.Sleep(100 * time.Millisecond)
				assert.Truef(t, m.secondaryClient == m.currentClient, "current Client must be secondary")
			},
		},
		{
			name: "SuccessWithSwitchToPrimary",
			mockFn: func(primaryMock *mockTypes.MockMinioClient) {
				gomock.InOrder(
					primaryMock.EXPECT().HealthCheck(gomock.Any()).Times(1).Return(nil, nil),
					primaryMock.EXPECT().IsOnline().Times(3).Return(false),
					primaryMock.EXPECT().IsOnline().Times(1).Return(true),
				)
			},
			testFn: func(t *testing.T, m *minio, clock *clockwork.FakeClock) {
				assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")

				for i := 0; i < 2; i++ {
					_ = clock.BlockUntilContext(builtinCtx.Background(), 1)
					clock.Advance(cfg.HealthCheckInterval)
					time.Sleep(100 * time.Millisecond)
					assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")
				}

				_ = clock.BlockUntilContext(builtinCtx.Background(), 1)
				clock.Advance(cfg.HealthCheckInterval)
				time.Sleep(100 * time.Millisecond)
				assert.Truef(t, m.secondaryClient == m.currentClient, "current Client must be secondary")

				_ = clock.BlockUntilContext(builtinCtx.Background(), 1)
				clock.Advance(cfg.HealthCheckInterval)
				time.Sleep(100 * time.Millisecond)
				assert.Truef(t, m.primaryClient == m.currentClient, "current Client must be primary")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			primaryMinioMock := mockTypes.NewMockMinioClient(ctrl)
			tt.mockFn(primaryMinioMock)
			secondaryMinioMock := mockTypes.NewMockMinioClient(ctrl)
			fakeClock := clockwork.NewFakeClock()
			m := &minio{
				currentClient:   primaryMinioMock,
				primaryClient:   primaryMinioMock,
				secondaryClient: secondaryMinioMock,
				cfg:             cfg,
				ctx:             ctx,
				clock:           fakeClock,
			}
			go m.startFallback()
			tt.testFn(t, m, fakeClock)
			ctx.Cancel()

		})
	}
}
