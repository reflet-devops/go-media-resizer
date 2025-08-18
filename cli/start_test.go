package cli

import (
	"github.com/labstack/echo/v4"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
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
