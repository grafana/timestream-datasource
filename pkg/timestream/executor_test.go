package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/timestream-datasource/pkg/models"
)

func runTest(t *testing.T, name string) *backend.DataResponse {
	mockClient := &MockClient{testFileName: name}
	dr := ExecuteQuery(context.Background(), models.QueryModel{}, mockClient)

	str := ""
	if dr.Error != nil {
		str = fmt.Sprintf("%+v", dr.Error)
	}

	if dr.Frames != nil {
		for idx, frame := range dr.Frames {
			metaString := ""
			if frame.Meta != nil {
				if frame.Meta.Custom != nil {
					frame.Meta.Custom["queryId"] = "{CHANGES}"
				}
				if frame.Meta.Stats != nil {
					frame.Meta.Stats = make([]string, 0) // avoid timing changes
				}

				meta, _ := json.MarshalIndent(frame.Meta, "", "    ")
				metaString = string(meta)
			}

			str += fmt.Sprintf("Frame[%d] %s\n", idx, string(metaString))

			table, _ := frame.StringTable(100, 10)
			str += table
			str += "\n\n\n"
		}
	}

	saved, err := mockClient.readText()
	if err != nil {
		//t.Error("Error reading saved value... recreating")
		err = mockClient.saveText(str)
		if err != nil {
			t.Error(err)
		}
	}

	if diff := cmp.Diff(saved, str); diff != "" {
		t.Fatalf("mismatch %s (-want +got):\n%s", name, diff)
	}
	return &dr
}

func TestSavedConversions(t *testing.T) {
	runTest(t, "describe-table")
	runTest(t, "select-star")
	runTest(t, "single-timeseries")
	runTest(t, "some-timeseries")
}

func TestGenerateTestData(t *testing.T) {
	//t.Skip("Integration Test") // comment line to run this

	m := make(map[string]models.QueryModel)
	m["describe-table.json"] = models.QueryModel{
		RawQuery: "DESCRIBE grafanaDB.grafanaTable",
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

	raw, _ := Interpolate(query)
	input := &timestreamquery.QueryInput{
		QueryString: aws.String(raw),
	}

	res, err := runner.runQuery(context.Background(), input)
	if err != nil {
		fmt.Println("execute failed", err.Error())
	}

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
