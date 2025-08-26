package cli

import (
	stdContext "context"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	buildinHttp "net/http"
	"os"
	"strconv"
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
					ctx.Cancel()
				}
			}
		}()

		errPid := managePidFile(ctx, ctx.Config.PidPath)
		if errPid != nil {
			return errPid
		}

		ctx.Logger.Info(fmt.Sprintf("starting server on %s", httpConfig.Listen))
		errStart := e.Start(httpConfig.Listen)
		if errStart != nil && errStart != buildinHttp.ErrServerClosed {
			_ = removePidFIle(ctx, ctx.Config.PidPath)
			panic(errStart)
		}

		return nil
	}
}

func managePidFile(ctx *context.Context, pidFile string) error {

	if _, err := ctx.Fs.Stat(pidFile); err == nil {
		data, _ := afero.ReadFile(ctx.Fs, pidFile)
		oldPid := string(data)
		return fmt.Errorf("pid file already exist with pid: %s", oldPid)
	}

	pid := os.Getpid()
	if err := afero.WriteFile(ctx.Fs, pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("pid file could not be written: %v", err)
	}

	go func() {
		<-ctx.Done()
		err := removePidFIle(ctx, pidFile)
		if err != nil {
			ctx.Logger.Error(fmt.Sprintf("failed to remove pid file: %v", err))
		}
	}()

	return nil
}

func removePidFIle(ctx *context.Context, pidFile string) error {
	return ctx.Fs.Remove(pidFile)
}
