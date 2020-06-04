package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

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
			meta, _ := json.MarshalIndent(frame.Meta, "", "    ")
			str += fmt.Sprintf("Frame[%d] %s\n", idx, string(meta))

			table, err := dr.Frames[0].StringTable(100, 10)
			if err != nil {
				t.Fatalf("error writing string table: %s", err)
			}
			str += table
			str += "\n\n\n"
		}
	}

	saved, err := mockClient.readText()
	if err != nil {
		t.Error("Error reading saved value... recreating")
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
	runTest(t, "simple-query")
}

func TestGenerateTestData(t *testing.T) {
	//t.Skip("Integration Test") // comment line to run this

	m := make(map[string]models.QueryModel)
	m["describe-table.json"] = models.QueryModel{
		RawQuery: "DESCRIBE grafanaDB.grafanaTable",
	}

	m["simple-query.json"] = models.QueryModel{
		Interval: time.Second * 5,
		TimeRange: backend.TimeRange{
			From: time.Unix(0, int64(time.Millisecond)*1588698110284),
			To:   time.Unix(0, int64(time.Millisecond)*1588700180087),
		},
		Database: "grafanaDB",
		Table:    "grafanaTable",
		RawQuery: `SELECT region, cell, silo, microservice_name,
		BIN(time, $__intervalStr) AS time_bin,
		ROUND(AVG(measure_value::double), 2) AS avg_value
	FROM ${database}.${table}
	WHERE $__timeFilter
		AND measure_name = 'cpu_user'
		AND region = 'us-east-1' AND cell = 'us-east-1-cell-1' AND microservice_name = 'apollo'
	GROUP BY region, cell, silo, microservice_name, BIN(time, $__intervalStr)
	ORDER BY time_bin DESC`,
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
