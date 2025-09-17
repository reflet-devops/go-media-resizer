package types

import (
	"io"
)

type Storage interface {
	//Type() string
	GetFile(path string) (io.ReadCloser, error)
	NotifyFileChange(chanEvent chan Events)
}
