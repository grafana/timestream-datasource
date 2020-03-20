package models

// QueryModel represents a spreadsheet query.
type QueryModel struct {
	RawQuery     string `json:"rawQuery"`
	NoTruncation bool   `json:"noTruncation"`
}

// TimestreamConfig contains config properties (share with other AWS services?)
type TimestreamConfig struct {
	// Nothing for now...
}
