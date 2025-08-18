package types

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/hash"
)

const (
	TagSourcePathHash = "source_path_hash"
)

func GetTagSourcePathHash(source string) string {
	sourcePathHash, _ := hash.GenerateMD5FromString(source)
	return FormatTag(TagSourcePathHash, sourcePathHash)
}

func FormatTag(key, value string) string {
	return fmt.Sprintf("%s_%s", key, value)
}
