package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateXXHashFromString_Success(t *testing.T) {
	got, err := GenerateXXHashFromString("test")
	assert.NoError(t, err)
	assert.Equal(t, "4fdcca5ddb678139", got)
}

func TestGenerateXXHashFromStringSuccess(t *testing.T) {
	got, err := GenerateXXHashFromBytes([]byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, "4fdcca5ddb678139", got)
}
