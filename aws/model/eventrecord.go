package model

type AwsEventRecord struct {
	EventSource     *string                  `json:"event_source,omitempty"`
	EventID         *string                  `json:"event_id,omitempty"`
	EventClass      *string                  `json:"event_class,omitempty"`
	EventType       *string                  `json:"event_type,omitempty"`
	EventTime       *int64                   `json:"event_time,omitempty"`
	EventAuthor     *string                  `json:"event_author,omitempty"`
	EventData       *map[string]interface{}  `json:"event_data,omitempty"`
	EntityID        *string                  `json:"entity_id,omitempty"`
	EntityType      *string                  `json:"entity_type,omitempty"`
	EntityState     *map[string]interface{}  `json:"entity_state,omitempty"`
	EntityCreatedAt *int64                   `json:"entity_created_at,omitempty"`
	EntityCreatedBy *string                  `json:"entity_created_by,omitempty"`
	TriggerSource   *string                  `json:"trigger_source,omitempty"`
	TriggerID       *string                  `json:"trigger_id,omitempty"`
	Transaction     *string                  `json:"transaction,omitempty"`
	Reference       *string 				 `json:"reference,omitempty"`
}