package model

type AwsEventHeaderRecord struct {
	EventSource *string `json:"event_source,omitempty"`
	EventID     *string `json:"event_id,omitempty"`
	EventClass  *string `json:"event_class,omitempty"`
	EventType   *string `json:"event_type,omitempty"`
	EventTime   *int64  `json:"event_time,omitempty"`
	EventDay    *string `json:"event_day,omitempty"`
	EventAuthor *string `json:"event_author,omitempty"`
}

type AwsEventHeaderRecordPage struct {
	Items     []*AwsEventHeaderRecord
	NextToken *string
}
