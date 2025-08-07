package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	got := DefaultConfig()
	assert.Equal(t,
		&Config{HTTP: HTTPConfig{Listen: "127.0.0.1:8080"}},
		got,
	)
}
