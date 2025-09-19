package urltools

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestJoinUri(t *testing.T) {
	tests := []struct {
		name string
		elem []string
		want string
	}{
		{
			name: "successWithEmpty",
			elem: []string{},
			want: "",
		},
		{
			name: "successWithEmptyFirstElement",
			elem: []string{"", "foo", "bar"},
			want: "foo/bar",
		},
		{
			name: "successWithFirstElement",
			elem: []string{"foo", "foo", "bar"},
			want: "foo/foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, JoinUri(tt.elem...), "JoinUri(%v)", tt.elem)
		})
	}
}

func TestRemoveProtocol(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "successWithHttp",
			url:  "http://hostname",
			want: "hostname",
		},
		{
			name: "successWithHttps",
			url:  "https://hostname",
			want: "hostname",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, RemoveProtocol(tt.url), "RemoveProtocol(%v)", tt.url)
		})
	}
}

func TestGetUri(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "successWithNotUri",
			url:  "http://hostname",
			want: "",
		},
		{
			name: "successWithUri",
			url:  "http://hostname/test",
			want: "test",
		},
		{
			name: "successWithUriMultiple",
			url:  "http://hostname/test/test",
			want: "test/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetUri(tt.url), "GetUri(%v)", tt.url)
		})
	}
}

func TestFormatPathWithPrefix(t *testing.T) {
	type args struct {
		prefixPath string
		path       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "successWithEmptyPrefix",
			args: args{
				prefixPath: "",
				path:       "test",
			},
			want: "test",
		},
		{
			name: "successWithEmptyPrefixAndSlashInPath",
			args: args{
				prefixPath: "",
				path:       "/test/test/",
			},
			want: "test/test",
		},
		{
			name: "successWithPrefixAndSlash",
			args: args{
				prefixPath: "/prefix",
				path:       "test",
			},
			want: "prefix/test",
		},
		{
			name: "successWithPrefixAndSlashPrefixAndPath",
			args: args{
				prefixPath: "/prefix/",
				path:       "/test/test/",
			},
			want: "prefix/test/test",
		},
		{
			name: "successWithPrefix",
			args: args{
				prefixPath: "prefix",
				path:       "test",
			},
			want: "prefix/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FormatPathWithPrefix(tt.args.prefixPath, tt.args.path), "FormatPathWithPrefix(%v, %v)", tt.args.prefixPath, tt.args.path)
		})
	}
}
