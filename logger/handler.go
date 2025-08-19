package logger

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var _ io.Writer = &HandlerRotateWriter{}

type HandlerRotateWriter struct {
	fs       afero.Fs
	file     io.Writer
	fallback io.Writer
	filepath string
	mutex    sync.Mutex
	done     <-chan bool
	sigs     chan os.Signal
}

func NewHandlerRotateWriter(fs afero.Fs, filename string, done <-chan bool) (*HandlerRotateWriter, error) {
	lw := &HandlerRotateWriter{
		fs:       fs,
		filepath: filename,
		mutex:    sync.Mutex{},
		done:     done,
		fallback: os.Stdout,
	}

	lw.sigs = make(chan os.Signal, 1)
	signal.Notify(lw.sigs, syscall.SIGUSR1, syscall.SIGHUP)

	_ = fs.MkdirAll(filepath.Dir(filename), 0755) // open will fail if MkdirAll fail
	err := lw.open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	go lw.ListenSignal()
	return lw, err
}

func (lw *HandlerRotateWriter) open() error {
	var err error
	lw.mutex.Lock()
	defer lw.mutex.Unlock()
	if lw.filepath == "" {
		lw.file = os.Stdout
		return nil
	}

	lw.file, err = lw.fs.OpenFile(lw.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	return err
}

func (lw *HandlerRotateWriter) ListenSignal() {
	for {
		select {
		case sigs := <-lw.sigs:
			err := lw.open()
			if err != nil {
				lw.mutex.Lock()
				lw.file = lw.fallback
				lw.mutex.Unlock()
				_, _ = fmt.Fprintf(lw.file, "fail to reopen logfile with sigs %s: %v. Switching to STDOUT", sigs.String(), err)
			}
		case <-lw.done:
			return
		}
	}
}

func (lw *HandlerRotateWriter) Write(p []byte) (n int, err error) {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()
	return lw.file.Write(p)
}
