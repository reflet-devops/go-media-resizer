package cli

import (
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
