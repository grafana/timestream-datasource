module github.com/grafana/timestream-datasource

go 1.14

replace github.com/aws/aws-sdk-go => ./tmp/github.com/aws/aws-sdk-go

require (
	github.com/aws/aws-sdk-go v1.29.27
	github.com/golang/protobuf v1.3.3 // indirect; indirect	github.com/grafana/grafana-plugin-sdk-go v0.28.0
	github.com/grafana/grafana-plugin-sdk-go v0.31.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/magefile/mage v1.9.0
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.3.0
)
