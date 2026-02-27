package context

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/reflet-devops/go-media-resizer/config"
	appProm "github.com/reflet-devops/go-media-resizer/prometheus"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/spf13/afero"
)

type Context struct {
	Logger     *slog.Logger
	LogLevel   *slog.LevelVar
	WorkingDir string
	Fs         afero.Fs
	sigs       chan os.Signal
	done       chan bool
	HttpClient types.Client

	BufferPool     *sync.Pool
	OptsResizePool *sync.Pool

	Config *config.Config

	MetricsRegistry appProm.Registry
}

func (c *Context) GetFS() afero.Fs {
	return c.Fs
}

func (c *Context) GetLogger() *slog.Logger {
	return c.Logger
}

func (c *Context) GetLogLevel() *slog.LevelVar {
	return c.LogLevel
}

func (c *Context) GetWorkingDir() string {
	return c.WorkingDir
}

func (c *Context) Cancel() {
	close(c.done)
}

func (c *Context) Done() <-chan bool {
	return c.done
}

func (c *Context) Signal() chan os.Signal {
	return c.sigs
}

func DefaultContext() *Context {
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return &Context{
		Logger:          slog.New(slog.NewTextHandler(os.Stdout, opts)),
		LogLevel:        level,
		WorkingDir:      workingDir,
		Fs:              afero.NewOsFs(),
		done:            make(chan bool),
		sigs:            sigs,
		Config:          config.DefaultConfig(),
		MetricsRegistry: prometheus.NewRegistry(),
		BufferPool: &sync.Pool{
			New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 5*1024*1024)) },
		},
		OptsResizePool: &sync.Pool{
			New: func() interface{} { return &types.ResizeOption{} },
		},
	}
}

func TestContext(logBuffer io.Writer) *Context {
	if logBuffer == nil {
		logBuffer = io.Discard
	}
	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return &Context{
		Logger:          slog.New(slog.NewTextHandler(logBuffer, opts)),
		LogLevel:        level,
		Fs:              afero.NewMemMapFs(),
		done:            make(chan bool),
		sigs:            sigs,
		Config:          config.DefaultConfig(),
		MetricsRegistry: prometheus.NewRegistry(),
		BufferPool: &sync.Pool{
			New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 1024*1024)) },
		},
		OptsResizePool: &sync.Pool{
			New: func() interface{} { return &types.ResizeOption{} },
		},
	}
}
