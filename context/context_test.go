package context

import (
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestDefaultContext_Success(t *testing.T) {

	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	fs := afero.NewOsFs()
	workingDir, err := os.Getwd()
	assert.NoError(t, err)
	want := &Context{
		WorkingDir: workingDir,
		Logger:     logger,
		LogLevel:   level,
		Fs:         fs,
		Config:     config.DefaultConfig(),
	}
	got := DefaultContext()
	assert.NotNil(t, got.done)
	got.done = nil
	got.sigs = nil
	assert.Equal(t, want, got)
}

func TestDefaultContext_FailGetwd(t *testing.T) {
	dir, _ := os.MkdirTemp("", "")
	_ = os.Chdir(dir)
	_ = os.RemoveAll(dir)
	assert.Panics(t, func() {
		DefaultContext()
	})
}

func TestTestContext(t *testing.T) {

	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	logger := slog.New(slog.NewTextHandler(io.Discard, opts))
	fs := afero.NewMemMapFs()
	want := &Context{
		Logger:   logger,
		LogLevel: level,
		Fs:       fs,
		Config:   config.DefaultConfig(),
	}
	got := TestContext(nil)
	assert.NotNil(t, got.done)
	got.done = nil
	got.sigs = nil
	assert.Equal(t, want, got)
}

func TestTestContext_WithLogBuffer(t *testing.T) {

	level := &slog.LevelVar{}
	level.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{AddSource: false, Level: level}
	logger := slog.New(slog.NewTextHandler(io.Discard, opts))
	fs := afero.NewMemMapFs()
	want := &Context{
		Logger:   logger,
		LogLevel: level,
		Fs:       fs,
		Config:   config.DefaultConfig(),
	}
	got := TestContext(io.Discard)
	assert.NotNil(t, got.done)
	got.done = nil
	got.sigs = nil
	assert.Equal(t, want, got)
}

func TestContext_Cancel(t *testing.T) {
	ctx := &Context{}
	ctx.done = make(chan bool)
	running := true
	go func() {
		select {
		case <-ctx.done:
			running = false
		}
	}()
	ctx.Cancel()
	assert.Equal(t, false, running)
}

func TestContext_Done(t *testing.T) {
	ctx := &Context{}
	ctx.done = make(chan bool)
	running := true
	go func() {
		select {
		case <-ctx.Done():
			running = false
		}
	}()
	ctx.done <- true
	assert.Equal(t, false, running)
}

func TestContext_Signal(t *testing.T) {
	ctx := &Context{}
	ctx.sigs = make(chan os.Signal, 1)
	signal.Notify(ctx.sigs, syscall.SIGINT, syscall.SIGTERM)
	running := true
	go func() {
		select {
		case <-ctx.Signal():
			running = false
		}
	}()
	ctx.Signal() <- syscall.SIGINT
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, false, running)
}

func TestContext_GetFS(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := &Context{
		Fs: fs,
	}
	assert.Equalf(t, fs, c.GetFS(), "GetFS()")
}

func TestContext_GetLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{AddSource: false, Level: &slog.LevelVar{}}))
	c := &Context{
		Logger: logger,
	}
	assert.Equalf(t, logger, c.GetLogger(), "GetLogger()")
}

func TestContext_GetLogLevel(t *testing.T) {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelInfo)
	c := &Context{
		LogLevel: logLevel,
	}
	assert.Equalf(t, logLevel, c.GetLogLevel(), "GetLogLevel()")
}

func TestContext_GetWorkingDir(t *testing.T) {
	workingDir := "/app"
	c := &Context{
		WorkingDir: workingDir,
	}
	assert.Equalf(t, workingDir, c.GetWorkingDir(), "GetWorkingDir()")
}
