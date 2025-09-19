package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTag(t *testing.T) {
	got := FormatTag("key", "value")
	assert.Equal(t, "key_value", got)
}

func TestGetTagSourcePathHash(t *testing.T) {
	got := GetTagSourcePathHash("app/text.txt")
	assert.Equal(t, fmt.Sprintf("%s_6ef8503f220bb4c34d072bddc808f136", TagSourcePathHash), got)
}

func TestFormatProjectPathHash(t *testing.T) {

	tests := []struct {
		name      string
		projectId string
		source    string
		want      string
	}{
		{
			name:      "success",
			projectId: "123",
			source:    "app/text.txt",
			want:      "123_app/text.txt",
		},
		{
			name:      "successWithSpace",
			projectId: " 123 ",
			source:    "app/text.txt",
			want:      "123_app/text.txt",
		},
		{
			name:      "successWithSlashInPath",
			projectId: " 123 ",
			source:    "/app/text.txt/",
			want:      "123_app/text.txt",
		},
		{
			name:      "successWithSpaceAndSlashInPath",
			projectId: " 123 ",
			source:    "/app/text.txt/",
			want:      "123_app/text.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FormatProjectPathHash(tt.projectId, tt.source), "FormatProjectPathHash(%v, %v)", tt.projectId, tt.source)
		})
	}
}
