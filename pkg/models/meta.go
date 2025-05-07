package models

import (
	timestreamquerytypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
)

// TimestreamCustomMeta is the standard metadata
type TimestreamCustomMeta struct {
	StartTime  int64 `json:"executionStartTime,omitempty"`
	FinishTime int64 `json:"executionFinishTime,omitempty"`

	NextToken string `json:"nextToken,omitempty"`
	QueryID   string `json:"queryId,omitempty"`
	RequestID string `json:"requestId,omitempty"`
	HasSeries bool   `json:"hasSeries,omitempty"`

	Status *timestreamquerytypes.QueryStatus `json:"status,omitempty"`
}
