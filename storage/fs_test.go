package storage

import (
	"errors"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	mockAfero "github.com/reflet-devops/go-media-resizer/mocks/afero"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"testing"
)

func Test_fs_Type(t *testing.T) {
	stateStorage := &fs{}
	assert.Equal(t, FsKey, stateStorage.Type())
}

func Test_createFsStorage(t *testing.T) {
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
				Type: FsKey,
				Config: map[string]interface{}{
					"prefix_path": "/app",
				},
			},
			want: &fs{fs: ctx.Fs, cfg: ConfigFs{PrefixPath: "/app"}},
		},
		{
			name: "SuccessTrimPrefixPath",
			cfg: config.StorageConfig{
				Type: FsKey,
				Config: map[string]interface{}{
					"prefix_path": "/app/",
				},
			},
			want: &fs{fs: ctx.Fs, cfg: ConfigFs{PrefixPath: "/app"}},
		},
		{
			name: "FailDecodeCfg",
			cfg: config.StorageConfig{
				Type: FsKey,
				Config: map[string]interface{}{
					"prefix_path": []string{},
				},
			},
			wantErr:     true,
			errContains: "prefix_path' expected type 'string', got unconvertible type '[]string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createFsStorage(ctx, tt.cfg)

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

func Test_fs_GetFile(t *testing.T) {

	cfg := ConfigFs{PrefixPath: "/app"}
	aferoFs := afero.NewMemMapFs()
	file, _ := aferoFs.Create("/app/foo/bar.tx")
	tests := []struct {
		name    string
		path    string
		mockFn  func(fsMock *mockAfero.MockFs)
		want    io.Reader
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			path: "foo/bar.txt",
			mockFn: func(fsMock *mockAfero.MockFs) {
				fsMock.EXPECT().Open(gomock.Eq("/app/foo/bar.txt")).Return(file, nil)
			},
			want:    file,
			wantErr: assert.NoError,
		},
		{
			name: "SuccessWithPrefixSlash",
			path: "/foo/bar.txt",
			mockFn: func(fsMock *mockAfero.MockFs) {
				fsMock.EXPECT().Open(gomock.Eq("/app/foo/bar.txt")).Return(file, nil)
			},
			want:    file,
			wantErr: assert.NoError,
		},
		{
			name: "FailedOpen",
			path: "/foo/bar.txt",
			mockFn: func(fsMock *mockAfero.MockFs) {
				fsMock.EXPECT().Open(gomock.Eq("/app/foo/bar.txt")).Return(nil, errors.New("open error"))
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fsMock := mockAfero.NewMockFs(ctrl)
			tt.mockFn(fsMock)
			f := fs{
				fs:  fsMock,
				cfg: cfg,
			}
			got, err := f.GetFile(tt.path)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFile(%v)", tt.path)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetFile(%v)", tt.path)
		})
	}
}
