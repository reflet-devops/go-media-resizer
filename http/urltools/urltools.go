package urltools

import (
	"path/filepath"
	"strings"
)

func RemoveProtocol(url string) string {
	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "https://", "", 1)
	return url
}

func GetUri(url string) string {
	urlArr := strings.Split(RemoveProtocol(url), "/")
	if len(urlArr) > 1 {
		return strings.Join(urlArr[1:], "/")
	}
	return ""
}

func GetHostname(url string) string {
	url = RemoveProtocol(url)
	hostname := strings.Split(url, "/")[0]
	return RemovePortNumber(hostname)
}

func RemovePortNumber(host string) string {
	return strings.Split(host, ":")[0]
}

func GetExtension(url string) string {
	return filepath.Ext(url)
}

func JoinUri(elem ...string) string {
	if len(elem) == 0 {
		return ""
	} else if elem[0] == "" {
		elem = elem[1:]
	}
	return strings.Join(elem, "/")
}

func FormatPathWithPrefix(prefixPath, path string) string {
	return JoinUri(strings.Trim(prefixPath, "/"), strings.Trim(path, "/"))
}
