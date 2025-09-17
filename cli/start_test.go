package cli

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
	pidPath := "/tmp/test.pid"
	ctx := context.TestContext(nil)
	ctx.Config.PidPath = pidPath
	ctx.Config.Projects = []config.Project{
		{
			ID:         "id",
			Hostname:   "example.com",
			PrefixPath: "prefix",
			Storage:    config.StorageConfig{Type: "fs"},
		},
	}

	errPid := afero.WriteFile(ctx.Fs, ctx.Config.PidPath, []byte(strconv.Itoa(os.Getpid())), 0644)
	assert.NoError(t, errPid)
	cmd := GetStartCmd(ctx)
	err := GetStartRunFn(ctx)(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("pid file (%s) already exist with pid: ", pidPath))

}

func Test_managePidFile(t *testing.T) {
	pidPath := "/tmp/test.pid"
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
				_ = afero.WriteFile(ctx.Fs, pidPath, []byte(strconv.Itoa(os.Getpid())), 0644)
			},
			wantErr:         true,
			wantErrContains: fmt.Sprintf("pid file (%s) already exist with pid: ", pidPath),
		},
		{
			name: "failedPidNotValid",
			mockFn: func(ctx *context.Context) {
				_ = afero.WriteFile(ctx.Fs, pidPath, []byte("wrong"), 0644)
			},
			wantErr:         true,
			wantErrContains: fmt.Sprintf("pid in %s is not valid: pid=wrong", pidPath),
		},
		{
			name: "failedWriteFile",
			mockFn: func(ctx *context.Context) {
				ctx.Fs = afero.NewReadOnlyFs(ctx.Fs)
			},
			wantErr:         true,
			wantErrContains: fmt.Sprintf("pid file (%s) could not be written: ", pidPath),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			ctx.Config.PidPath = pidPath
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

func Test_removePidFIle(t *testing.T) {
	ctx := context.TestContext(nil)
	_ = afero.WriteFile(ctx.Fs, ctx.Config.PidPath, []byte("1234"), 0644)
	err := removePidFIle(ctx, ctx.Config.PidPath)
	assert.NoError(t, err)
}

func Test_isProcessRunning_SuccessRunning(t *testing.T) {
	assert.True(t, isProcessRunning(os.Getpid()))
}

func Test_isProcessRunning_SuccessNotRunning(t *testing.T) {
	assert.False(t, isProcessRunning(99999999999))
}
