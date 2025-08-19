package types

const (
	HeaderCachePurge = "X-Cache-Tag"
)

type PurgeCache interface {
	//Type() string
	Purge(events Events)
}
