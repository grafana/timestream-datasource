module github.com/grafana/timestream-datasource

go 1.14

replace github.com/aws/aws-sdk-go => ./tmp/github.com/aws/aws-sdk-go

require (
	github.com/aws/aws-sdk-go v1.30.20
	github.com/grafana/grafana-plugin-sdk-go v0.60.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/magefile/mage v1.9.0
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/prometheus/client_golang v1.3.0
)
