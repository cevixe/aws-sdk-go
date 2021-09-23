package model

type AwsEventHeaderRecord struct {
	EventSource   *string `json:"event_source,omitempty"`
	EventID       *string `json:"event_id,omitempty"`
	EventClass    *string `json:"event_class,omitempty"`
	EventType     *string `json:"event_type,omitempty"`
	EventTime     *int64  `json:"event_time,omitempty"`
	EventDay      *string `json:"event_day,omitempty"`
	EventAuthor   *string `json:"event_author,omitempty"`
	EntityID      *string `json:"entity_id,omitempty"`
	EntityType    *string `json:"entity_type,omitempty"`
	EntityDeleted bool    `json:"entity_deleted,omitempty"`
	Transaction   *string `json:"transaction,omitempty"`
	TriggerSource *string `json:"trigger_source,omitempty"`
	TriggerID     *string `json:"trigger_id,omitempty"`
}

type AwsEventHeaderRecordPage struct {
	Items     []*AwsEventHeaderRecord
	NextToken *string
}
