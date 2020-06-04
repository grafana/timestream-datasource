module github.com/grafana/timestream-datasource

go 1.14

replace github.com/aws/aws-sdk-go => ./tmp/github.com/aws/aws-sdk-go

require (
	github.com/aws/aws-sdk-go v1.31.10
	github.com/google/go-cmp v0.4.1
	github.com/grafana/grafana-plugin-sdk-go v0.66.0
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
)
