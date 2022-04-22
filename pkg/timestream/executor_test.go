package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/timestream-datasource/pkg/models"
)

func runTest(t *testing.T, names []string) *backend.DataResponse {
	mockClient := &MockClient{testFileNames: names}
	dr := ExecuteQuery(context.Background(), models.QueryModel{WaitForResult: true}, mockClient, models.DatasourceSettings{})

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

	// Set the last parameter of CheckGoldenDataResponse to true to write new golden responses
	err := experimental.CheckGoldenDataResponse("./testdata/"+names[0]+".golden.txt", &dr, false)
	if err != nil {
		t.Errorf(err.Error())
	}

	return &dr
}

func TestSavedConversions(t *testing.T) {
	runTest(t, []string{"select-consts"})
	runTest(t, []string{"describe-table"})
	runTest(t, []string{"select-star"})
	runTest(t, []string{"select-null-timestamp"})
	runTest(t, []string{"complex-timeseries"})
	runTest(t, []string{"some-timeseries"})
	runTest(t, []string{"show-measures"})
	runTest(t, []string{"show-databases"})
	runTest(t, []string{"show-tables"})
	runTest(t, []string{"pagination-off_1", "pagination-off_2"})
}

func TestGenerateTestData(t *testing.T) {
	// This will do real API calls to AWS to populate test data
	t.Skip("Integration Test") // comment line to run this
	db := "grafanaDB"
	tableName := "DevOps"
	table := db + "." + tableName

	m := make(map[string]models.QueryModel)
	m["show-databases.json"] = models.QueryModel{
		RawQuery: "SHOW DATABASES",
	}

	m["show-tables.json"] = models.QueryModel{
		RawQuery: "SHOW TABLES FROM " + db,
	}

	m["select-consts.json"] = models.QueryModel{
		RawQuery: `SELECT
		  1     as t_int32,
		  'two' as t_varchar,
		  timestamp '2020-08-08 01:00' as timestamp,
		  interval '2' day - interval '3' hour as interval_day_to_second,
		  interval '3' year - interval '5' month as interval_year_to_month,
		  time '01:00' as time,
		  date '2020-08-08' as date
		`,
	}

	m["describe-table.json"] = models.QueryModel{
		RawQuery: "DESCRIBE " + table,
	}

	m["show-measures.json"] = models.QueryModel{
		RawQuery: "SHOW MEASURES FROM " + table,
	}

	m["select-star.json"] = models.QueryModel{
		RawQuery: `SELECT * FROM ` + table + ` LIMIT 10`,
	}

	m["select-null-timetamp.json"] = models.QueryModel{
		RawQuery: `SELECT measure_name ,
		CASE WHEN measure_name = 'make_me_null' THEN (SELECT NULL) ELSE time END
		FROM ` + table + 
		` LIMIT 10`,
	}

	m["complex-timeseries.json"] = models.QueryModel{
		RawQuery: `select measure_name, availability_zone, region, cell, silo, instance_type, instance_name, create_time_series(time, measure_value::double)
		from ` + table + `
		where time > ago(30m)
			AND (measure_name = 'cpu_user' or measure_name = 'cpu_system')
			and availability_zone ='us-east-1-1'
			and microservice_name = 'hercules'
			and region = 'us-east-1'
			and cell = 'us-east-1-cell-1'
			and silo = 'us-east-1-cell-1-silo-1'
			and instance_type = 'r5.4xlarge'
			group by measure_name, availability_zone, region, cell, silo, instance_type, instance_name
		`,
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
		FROM ` + table + `
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

	m["pagination-off_1.json"] = models.QueryModel{
		RawQuery: `SELECT * FROM ` + table + " LIMIT 3",
	}
	m["pagination-off_2.json"] = models.QueryModel{
		RawQuery: `SELECT * FROM ` + table + " LIMIT 3",
	}

	nextToken := map[string]*string{}
	for filename, query := range m {
		inst, err := NewServerInstance(backend.DataSourceInstanceSettings{
			JSONData: []byte(`{"region": "us-west-2"}`),
		})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		runner := inst.(*timestreamDS).Runner
		raw, _ := Interpolate(query, models.DatasourceSettings{})
		input := &timestreamquery.QueryInput{
			QueryString: aws.String(raw),
		}
		// Custom input to test pagination
		if filename == "pagination-off_1.json" {
			// force pagination
			input.MaxRows = aws.Int64(1)
		}
		if filename == "pagination-off_2.json" {
			// continue previous page
			input.NextToken = nextToken["pagination-off_1.json"]
			if input.NextToken == nil {
				t.Fatalf("first page should be executed first. Please retry")
			}
		}

		res, err := runner.runQuery(context.Background(), input)
		if err != nil {
			fmt.Println("execute failed", err.Error())
		}
		nextToken[filename] = res.NextToken

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
}
