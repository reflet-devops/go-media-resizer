package logger

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/context"
	mockAfero "github.com/reflet-devops/go-media-resizer/mocks/afero"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_NewHandlerRotateWriter_Success(t *testing.T) {
	ctx := context.TestContext(nil)

	filepath := "/my/path.txt"

	got, err := NewHandlerRotateWriter(ctx.Fs, filepath, ctx.Done())

	logFile, _ := got.fs.Open(filepath)

	assert.Nil(t, err)
	assert.Equal(t, ctx.Fs, got.fs)
	assert.Equal(t, ctx.Done(), got.done)
	assert.Equal(t, filepath, got.filepath)
	assert.Equal(t, filepath, logFile.Name())
	ctx.Cancel()
}

func Test_NewHandlerRotateWriter_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Fs = afero.NewReadOnlyFs(ctx.Fs)
	filepath := "/my/path.txt"

	_, err := NewHandlerRotateWriter(ctx.Fs, filepath, ctx.Done())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file: ")
}

func Test_NewHandlerRotateWriter_Stdout_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	filepath := ""

	lw, err := NewHandlerRotateWriter(ctx.Fs, filepath, ctx.Done())
	assert.Nil(t, err)
	assert.Equal(t, lw.file, os.Stdout)
}

func Test_Write_Fail(t *testing.T) {
	ctx := context.TestContext(nil)

	filepath := "/my/path.txt"

	lw, _ := NewHandlerRotateWriter(ctx.Fs, filepath, ctx.Done())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buff := []byte("hello world")

	mockFile := mockAfero.NewMockFile(ctrl)
	mockFile.EXPECT().Write(gomock.Eq(buff)).Return(10, fmt.Errorf("mock error"))
	lw.file = mockFile

	n, err := lw.Write(buff)
	assert.Equal(t, 10, n)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock error")
}

func Test_ListenSignal_Sig(t *testing.T) {
	ctx := context.TestContext(nil)

	filepath := "/my/path.txt"
	lw, _ := NewHandlerRotateWriter(ctx.Fs, filepath, ctx.Done())
	lw.fallback = io.Discard

	time.Sleep(time.Millisecond * 100)
	_ = lw.fs.Remove(filepath)

	rwFs := lw.fs
	lw.fs = afero.NewReadOnlyFs(lw.fs)

	lw.sigs <- syscall.SIGUSR1
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, lw.file, lw.fallback)
	time.Sleep(time.Millisecond * 100)
	lw.fs = rwFs

	lw.sigs <- syscall.SIGHUP
	time.Sleep(time.Millisecond * 100)

	logFile, _ := lw.fs.Open(filepath)
	assert.Equal(t, logFile.Name(), "/my/path.txt")
	time.Sleep(time.Millisecond * 100)

}
