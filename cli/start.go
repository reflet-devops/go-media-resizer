package cli

import (
	stdContext "context"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http"
	"github.com/spf13/cobra"
	buildinHttp "net/http"
)

func GetStartCmd(ctx *context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start server",
		RunE:  GetStartRunFn(ctx),
	}
}

func GetStartRunFn(ctx *context.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		e := http.CreateServerHTTP(ctx)

		httpConfig := ctx.Config.HTTP
		go func() {
			ctx.Logger.Info(fmt.Sprintf("starting server on %s", httpConfig.Listen))
			errStart := e.Start(httpConfig.Listen)
			if errStart != nil && errStart != buildinHttp.ErrServerClosed {
				panic(errStart)
			}
		}()

		for {
			select {
			case sig := <-ctx.Signal():
				//ctx.Cancel()
				ctx.Logger.Info(fmt.Sprintf("%s signal received, exiting...", sig.String()))
				return e.Shutdown(stdContext.Background())
			}
		}

	}
}
