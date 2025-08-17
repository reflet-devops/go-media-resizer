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
		e, err := http.CreateServerHTTP(ctx)
		if err != nil {
			return err
		}

		httpConfig := ctx.Config.HTTP

		go func() {
			for {
				select {
				case sig := <-ctx.Signal():
					ctx.Logger.Info(fmt.Sprintf("%s signal received, exiting...", sig.String()))
					_ = e.Shutdown(stdContext.Background())
					ctx.Done()
				}
			}
		}()

		ctx.Logger.Info(fmt.Sprintf("starting server on %s", httpConfig.Listen))
		errStart := e.Start(httpConfig.Listen)
		if errStart != nil && errStart != buildinHttp.ErrServerClosed {
			panic(errStart)
		}

		return nil
	}
}
