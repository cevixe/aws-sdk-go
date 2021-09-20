package model

type AwsStateRecord struct {
	Type      string                 `json:"type,omitempty"`
	ID        string                 `json:"id,omitempty"`
	Version   uint64                 `json:"version,omitempty"`
	State     map[string]interface{} `json:"state,omitempty"`
	UpdatedAt int64                  `json:"updated_at,omitempty"`
	UpdatedBy string                 `json:"updated_by,omitempty"`
	CreatedAt int64                  `json:"created_at,omitempty"`
	CreatedBy string                 `json:"created_by,omitempty"`

	ContentLocation string `json:"content_location,omitempty"`
	ContentType     string `json:"content_type,omitempty"`
	ContentEncoding string `json:"content_encoding,omitempty"`
	Content         []byte `json:"content,omitempty"`
}

type AwsStateRecordPage struct {
	Items     []*AwsStateRecord
	NextToken *string
}
