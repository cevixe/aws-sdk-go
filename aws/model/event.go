package model

type EventObject struct {
	Reference     *ObjectRef              `json:"_ref,omitempty"`
	SourceID      string                  `json:"source_id,omitempty"`
	SourceType    string                  `json:"source_type,omitempty"`
	SourceTime    int64                   `json:"source_time,omitempty"`
	SourceOwner   string                  `json:"source_owner,omitempty"`
	SourceState   *map[string]interface{} `json:"source_state,omitempty"`
	EventID       uint64                  `json:"event_id,omitempty"`
	EventType     string                  `json:"event_type,omitempty"`
	EventTime     int64                   `json:"event_time,omitempty"`
	EventAuthor   string                  `json:"event_author,omitempty"`
	EventPayload  *map[string]interface{} `json:"event_payload,omitempty"`
	TriggerSource string                  `json:"trigger_source,omitempty"`
	TriggerID     uint64                  `json:"trigger_id,omitempty"`
	Transaction   string                  `json:"transaction,omitempty"`
}
