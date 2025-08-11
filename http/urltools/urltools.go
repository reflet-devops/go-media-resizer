package urltools

import "strings"

func GetHostname(url string) string {

	strings.Replace(url, "http://", "", 1)
	strings.Replace(url, "https://", "", 1)
	hostname := strings.Split(url, "/")[0]
	return strings.Split(hostname, ":")[0]
}
