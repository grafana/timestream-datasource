package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	gaws "github.com/grafana/timestream-datasource/pkg/common/aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/timestream-datasource/pkg/models"
)

func runTest(t *testing.T, name string) *backend.DataResponse {
	mockClient := &MockClient{testFileName: name}
	dr := ExecuteQuery(context.Background(), models.QueryModel{}, mockClient, gaws.DatasourceSettings{})

	// Remove changable fields
	for _, frame := range dr.Frames {
		if frame.Meta != nil {
			meta := frame.Meta.Custom.(*models.TimestreamCustomMeta)
			meta.StartTime = 1111
			meta.FinishTime = 2222
			if meta.QueryID != "" {
				meta.QueryID = "$queryId$"
			}
			if meta.NextToken != "" {
				meta.NextToken = "$nextToken$"
			}

			// TODO: fix: https://github.com/grafana/grafana-plugin-sdk-go/issues/213
			frame.Meta.Custom = nil
		}
	}

	err := experimental.CheckGoldenDataResponse("./testdata/"+name+".golden.txt", &dr, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	return &dr
}

func TestSavedConversions(t *testing.T) {
	runTest(t, "select-consts")
	runTest(t, "describe-table")
	runTest(t, "select-star")
	runTest(t, "single-timeseries")
	runTest(t, "some-timeseries")
	runTest(t, "show-measures")
}

func TestGenerateTestData(t *testing.T) {
	t.Skip("Integration Test") // comment line to run this

	m := make(map[string]models.QueryModel)
	m["select-consts.json"] = models.QueryModel{
		RawQuery: `SELECT 
		  1     as t_int32, 
		  'two' as t_varchar
		`,
	}

	m["describe-table.json"] = models.QueryModel{
		RawQuery: "DESCRIBE grafanaDB.grafanaTable",
	}

	m["show-measures.json"] = models.QueryModel{
		RawQuery: "SHOW MEASURES FROM grafanaDB.grafanaTable",
	}

	m["select-star.json"] = models.QueryModel{
		RawQuery: `SELECT * FROM grafanaDB.grafanaTable LIMIT 10`,
	}

	m["single-timeseries.json"] = models.QueryModel{
		RawQuery: `SELECT region, cell, silo, availability_zone, microservice_name,
		instance_name, process_name, jdk_version,
		CREATE_TIME_SERIES(time, measure_value::double) AS gc_reclaimed
	FROM grafanaDB.grafanaTable
	WHERE time > ago(2h)
		AND measure_name = 'gc_reclaimed'
		AND region = 'ap-northeast-1' AND cell = 'ap-northeast-1-cell-5' AND silo = 'ap-northeast-1-cell-5-silo-2'
		AND availability_zone = 'ap-northeast-1-3' AND microservice_name = 'zeus'
		AND instance_name = 'i-zaZswmJk-zeus-0002.amazonaws.com' AND process_name = 'server' AND jdk_version = 'JDK_11'
	GROUP BY region, cell, silo, availability_zone, microservice_name,
		instance_name, process_name, jdk_version`,
	}

	m["some-timeseries.json"] = models.QueryModel{
		RawQuery: `SELECT 
			region, 
			cell, 
			silo, 
			availability_zone, 
			microservice_name,
			instance_name, 
			process_name, 
			jdk_version,
			CREATE_TIME_SERIES(time, measure_value::double) AS gc_reclaimed
		FROM grafanaDB.grafanaTable
		WHERE time > ago(1h)
			AND measure_name = 'gc_reclaimed'
			AND region = 'ap-northeast-1' 
			AND cell = 'ap-northeast-1-cell-5' 
			AND silo = 'ap-northeast-1-cell-5-silo-2'
			AND availability_zone = 'ap-northeast-1-3' 
			AND microservice_name = 'zeus'
		GROUP BY region, 
			cell, 
			silo, 
			availability_zone, 
			microservice_name,
			instance_name, 
			process_name, 
			jdk_version		
		ORDER BY AVG(measure_value::double) DESC
		LIMIT 10`,
	}

	for key, value := range m {
		writeTestData(key, value, t)
	}
}

// This will write the results to local json file
func writeTestData(filename string, query models.QueryModel, t *testing.T) {

	inst, err := newDataSourceInstance(backend.DataSourceInstanceSettings{
		JSONData: []byte(`{"region": "us-east-1"}`),
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	runner := inst.(*instanceSettings).Runner

	raw, _ := Interpolate(query, gaws.DatasourceSettings{})
	input := &timestreamquery.QueryInput{
		QueryString: aws.String(raw),
	}

	res, err := runner.runQuery(context.Background(), input)
	if err != nil {
		fmt.Println("execute failed", err.Error())
	}

	// This changes with every request, so make it the same
	res.QueryId = aws.String("#QueryId#")

	json, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		fmt.Println("marshalling results failed", err.Error())
	}

	f, err := os.Create("./testdata/" + filename)
	if err != nil {
		fmt.Println("create file failed: ", filename)
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = f.Write(json)
	if err != nil {
		fmt.Println("write file failed: ", filename)
	}
}
