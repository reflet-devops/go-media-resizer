package types

import (
	"fmt"
	"strings"

	"github.com/reflet-devops/go-media-resizer/hash"
)

const (
	TagSourcePathHash = "source_path_hash"
)

func GetTagSourcePathHash(value string) string {
	sourcePathHash, _ := hash.GenerateXXHashFromString(value)
	return FormatTag(TagSourcePathHash, sourcePathHash)
}

func FormatProjectPathHash(projectId string, source string) string {
	return fmt.Sprintf("%s_%s", strings.TrimSpace(projectId), strings.Trim(source, "/"))
}

func FormatTag(key, value string) string {
	return fmt.Sprintf("%s_%s", key, value)
}
