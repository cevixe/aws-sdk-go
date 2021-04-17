package model

type ObjectRef struct {
	Bucket  string  `json:"bucket,omitempty"`
	Region  string  `json:"region,omitempty"`
	Key     string  `json:"key,omitempty"`
	Version *string `json:"version,omitempty"`
}
