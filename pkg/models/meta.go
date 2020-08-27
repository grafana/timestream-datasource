package models

// TimestreamCustomMeta is the standard metadata
type TimestreamCustomMeta struct {
	StartTime  int64 `json:"executionStartTime,omitempty"`
	FinishTime int64 `json:"executionFinishTime,omitempty"`

	NextToken string `json:"nextToken,omitempty"`
	QueryID   string `json:"queryId,omitempty"`
	RequestID string `json:"requestId,omitempty"`
	HasSeries bool   `json:"hasSeries,omitempty"`
}
