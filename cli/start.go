package cli

import (
	stdContext "context"
	"fmt"
	buildinHttp "net/http"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/http"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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
					ctx.Cancel()
					ctx.Logger.Info(fmt.Sprintf("gracefull shutdown completed"))
					_ = e.Shutdown(stdContext.Background())
				}
			}
		}()

		errPid := managePidFile(ctx, ctx.Config.PidPath)
		if errPid != nil {
			return errPid
		}

		defer func() {
			errRemovePid := removePidFIle(ctx, ctx.Config.PidPath)
			if errRemovePid != nil {
				ctx.Logger.Error(fmt.Sprintf("failed to remove pid file: %v", errRemovePid))
			}
		}()

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
		strPid := strings.TrimSpace(string(data))
		if strPid != "" {
			oldPid, errConv := strconv.Atoi(strPid)
			if errConv != nil {
				return fmt.Errorf("pid in %s is not valid: pid=%s", pidFile, string(data))
			}

			if isProcessRunning(oldPid) {
				return fmt.Errorf("pid file (%s) already exist with pid: %d", pidFile, oldPid)
			}
		}
	}

	pid := os.Getpid()
	if err := afero.WriteFile(ctx.Fs, pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("pid file (%s) could not be written: %v", pidFile, err)
	}

	return nil
}

func removePidFIle(ctx *context.Context, pidFile string) error {
	return ctx.Fs.Remove(pidFile)
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}
