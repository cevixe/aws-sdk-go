package model

type AwsControlRecord struct {
	ControlGroup   string         `json:"control_group,omitempty"`
	ControlID      string         `json:"control_id,omitempty"`
	ControlIntent  uint64         `json:"control_intent,omitempty"`
	ControlType    AwsControlType `json:"control_type,omitempty"`
	ControlTime    int64          `json:"control_time,omitempty"`
	HandlerID      string         `json:"handler_id,omitempty"`
	HandlerVersion string         `json:"handler_version,omitempty"`
	HandlerTimeout uint64         `json:"handler_timeout,omitempty"`
	EventSource    string         `json:"event_source,omitempty"`
	EventID        string         `json:"event_id,omitempty"`
	Transaction    string         `json:"transaction,omitempty"`
}

type AwsControlType string

const (
	BlockControl   AwsControlType = "B"
	ConfirmControl AwsControlType = "C"
)
