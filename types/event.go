package types

const (
	EventTypeUnknow = "unknow"
	EventTypePurge  = "purge"
)

type Events []Event
type Event struct {
	Type string `json:"type" mapstructure:"type" validate:"required"`
	Path string `json:"path" mapstructure:"path" validate:"required"`
}
