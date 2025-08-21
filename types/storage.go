package types

import (
	"io"
)

type Storage interface {
	//Type() string
	GetFile(path string) (io.Reader, error)
	NotifyFileChange(chanEvent chan Events)
}
