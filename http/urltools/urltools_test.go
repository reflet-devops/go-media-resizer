package urltools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetHostname(t *testing.T) {
	addr := "hostname:port"
	assert.Equal(t, "hostname", GetHostname(addr))
}
