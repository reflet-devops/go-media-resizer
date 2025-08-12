package types

import (
	"io"
)

type Storage interface {
	Type() string
	GetPrefix() string
	GetFile(path string) (io.Reader, error)
}
