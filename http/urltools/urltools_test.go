package urltools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RemovePortNumber(t *testing.T) {
	addr := "hostname:port"
	assert.Equal(t, "hostname", GetHostname(addr))
}

func Test_GetHostname(t *testing.T) {

	want := "hostname"
	tests := []struct {
		name string
		arg  string
	}{
		{
			name: "http",
			arg:  "http://hostname",
		},
		{
			name: "https",
			arg:  "http://hostname",
		},
		{
			name: "with path",
			arg:  "http://hostname/path/to/resource",
		},
		{
			name: "port number",
			arg:  "https://hostname:433",
		},
		{
			name: "port number and path",
			arg:  "https://hostname:433/path/to/resource",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHostname(tt.arg)
			assert.Equal(t, want, got)
		})
	}
}

func Test_GetExtension(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "single level resource path",
			arg:  "http://hostname/filename.txt",
			want: ".txt",
		},
		{
			name: "multi-level resource path",
			arg:  "http://hostname/path/to/filename.txt",
			want: ".txt",
		},
		{
			name: "multi-level resource path",
			arg:  "http://hostname/path/to/filename.png.pdf.txt",
			want: ".txt",
		},
		{
			name: "no extension",
			arg:  "http://hostname/filename",
			want: "",
		},
		{
			name: "empty resource path",
			arg:  "http://hostname/",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetExtension(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}
