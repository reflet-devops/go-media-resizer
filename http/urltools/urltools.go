package urltools

import "strings"

func GetHostname(url string) string {

	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "https://", "", 1)
	hostname := strings.Split(url, "/")[0]
	return RemovePortNumber(hostname)
}

func RemovePortNumber(host string) string {
	return strings.Split(host, ":")[0]
}
