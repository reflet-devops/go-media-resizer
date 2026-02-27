package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*func TestGenerateMD5_Success(t *testing.T) {
	w := bytes.NewBufferString("test")
	got, err := GenerateMD5(w)
	assert.NoError(t, err)
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", got)
}

func TestGenerateMD5_Fail(t *testing.T) {
	w := &errorReader{r: bytes.NewBufferString("test")}
	got, err := GenerateMD5(w)
	assert.Error(t, err)
	assert.Equal(t, got, "")
}*/

func TestGenerateMD5FromString_Success(t *testing.T) {
	got, err := GenerateMD5FromString("test")
	assert.NoError(t, err)
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", got)
}

func TestGenerateMD5FromBytesSuccess(t *testing.T) {
	got, err := GenerateMD5FromBytes([]byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", got)
}
