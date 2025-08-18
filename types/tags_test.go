package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatTag(t *testing.T) {
	got := FormatTag("key", "value")
	assert.Equal(t, "key_value", got)
}

func TestGetTagSourcePathHash(t *testing.T) {
	got := GetTagSourcePathHash("app/text.txt")
	assert.Equal(t, fmt.Sprintf("%s_6ef8503f220bb4c34d072bddc808f136", TagSourcePathHash), got)
}
