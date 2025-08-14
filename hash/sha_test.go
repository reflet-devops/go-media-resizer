package hash

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type errorReader struct {
	r     io.Reader
	limit int
	count int
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	if e.count >= e.limit {
		return 0, fmt.Errorf("simulated read error")
	}
	n, err = e.r.Read(p)
	e.count += n
	return
}

func TestGenerateSHA256_Success(t *testing.T) {
	w := bytes.NewBufferString("test")
	got, err := GenerateSHA256(w)
	assert.NoError(t, err)
	assert.Equal(t, got, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
}

func TestGenerateSHA256_Fail(t *testing.T) {
	w := &errorReader{r: bytes.NewBufferString("test")}
	got, err := GenerateSHA256(w)
	assert.Error(t, err)
	assert.Equal(t, got, "")
}

func TestGenerateSHA256FromString_Success(t *testing.T) {
	got, err := GenerateSHA256FromString("test")
	assert.NoError(t, err)
	assert.Equal(t, got, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
}
