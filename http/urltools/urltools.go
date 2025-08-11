package urltools

import "strings"

func GetHostname(url string) string {
	return strings.Split(url, ":")[0]
}
