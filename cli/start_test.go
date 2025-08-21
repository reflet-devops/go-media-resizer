package cli

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"syscall"
	"testing"
	"time"
)

func TestGetStartRunFn_SuccessOnlyListenHTTP(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.HTTP.Listen = "127.0.0.1:0"
	viper.Reset()
	viper.SetFs(ctx.Fs)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cmd := GetStartCmd(ctx)
	go func() {
		err := GetStartRunFn(ctx)(cmd, []string{})
		assert.NoError(t, err)
	}()
	time.Sleep(time.Millisecond * 500)
	ctx.Signal() <- syscall.SIGINT
}

func TestGetStartRunFn_FailPortAlreadyBind(t *testing.T) {

	ctx := context.TestContext(nil)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	go func() {
		_ = e.Start("127.0.0.1:0")
	}()
	time.Sleep(time.Millisecond * 500)

	ctx.Config.HTTP.Listen = e.Listener.Addr().String()
	cmd := GetStartCmd(ctx)

	assert.Panics(t, func() {
		_ = GetStartRunFn(ctx)(cmd, []string{})
	})
	_ = e.Close()
}

func TestGetStartRunFn_FailCreateServerHTTP(t *testing.T) {

	ctx := context.TestContext(nil)
	ctx.Config.Projects = []config.Project{
		{
			ID:         "id",
			Hostname:   "example.com",
			PrefixPath: "prefix",
			Storage:    config.StorageConfig{Type: "wrong"},
		},
	}
	cmd := GetStartCmd(ctx)

	err := GetStartRunFn(ctx)(cmd, []string{})
	assert.Error(t, err)

}

func TestGetStartRunFn_FailCreatePID(t *testing.T) {

	ctx := context.TestContext(nil)
	ctx.Config.Projects = []config.Project{
		{
			ID:         "id",
			Hostname:   "example.com",
			PrefixPath: "prefix",
			Storage:    config.StorageConfig{Type: "fs"},
		},
	}

	errPid := afero.WriteFile(ctx.Fs, ctx.Config.PidPath, []byte("1234"), 0644)
	assert.NoError(t, errPid)
	cmd := GetStartCmd(ctx)
	err := GetStartRunFn(ctx)(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pid file already exist with pid: ")

}

func Test_managePidFile(t *testing.T) {

	tests := []struct {
		name            string
		mockFn          func(ctx *context.Context)
		wantErr         bool
		wantErrContains string
	}{
		{
			name:    "success",
			mockFn:  func(ctx *context.Context) {},
			wantErr: false,
		},
		{
			name: "failedPidExist",
			mockFn: func(ctx *context.Context) {
				_ = afero.WriteFile(ctx.Fs, ctx.Config.PidPath, []byte("1234"), 0644)
			},
			wantErr:         true,
			wantErrContains: "pid file already exist with pid",
		},
		{
			name: "failedWriteFile",
			mockFn: func(ctx *context.Context) {
				ctx.Fs = afero.NewReadOnlyFs(ctx.Fs)
			},
			wantErr:         true,
			wantErrContains: "pid file could not be written",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			tt.mockFn(ctx)
			err := managePidFile(ctx, ctx.Config.PidPath)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				assert.NoError(t, err)
				time.Sleep(time.Millisecond * 100)
				ctx.Cancel()
			}
		})
	}
}

func Test_managePidFile_FailedRemoveOnDone(t *testing.T) {
	b := bytes.NewBufferString("")
	ctx := context.TestContext(b)
	err := managePidFile(ctx, ctx.Config.PidPath)
	assert.NoError(t, err)
	_ = ctx.Fs.Remove(ctx.Config.PidPath)
	ctx.Cancel()
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, b.String(), "failed to remove pid file")
}

func Test_removePidFIle(t *testing.T) {
	ctx := context.TestContext(nil)
	_ = afero.WriteFile(ctx.Fs, ctx.Config.PidPath, []byte("1234"), 0644)
	err := removePidFIle(ctx, ctx.Config.PidPath)
	assert.NoError(t, err)
}
