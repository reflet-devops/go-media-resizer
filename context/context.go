package context

import (
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/spf13/afero"
	"github.com/valyala/fasthttp"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Context struct {
	Logger     *slog.Logger
	LogLevel   *slog.LevelVar
	WorkingDir string
	Fs         afero.Fs
	sigs       chan os.Signal
	done       chan bool
	HttpClient types.Client

	Config *config.Config
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
	c.done <- true
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
		Logger:     slog.New(slog.NewTextHandler(os.Stdout, opts)),
		LogLevel:   level,
		WorkingDir: workingDir,
		Fs:         afero.NewOsFs(),
		done:       make(chan bool),
		sigs:       sigs,
		Config:     config.DefaultConfig(),
		HttpClient: &fasthttp.Client{},
	}
}

func (c *Context) Clone() *Context {
	newCtx := *c
	return &newCtx
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
		Logger:   slog.New(slog.NewTextHandler(logBuffer, opts)),
		LogLevel: level,
		Fs:       afero.NewMemMapFs(),
		done:     make(chan bool),
		sigs:     sigs,
		Config:   config.DefaultConfig(),
	}
}
