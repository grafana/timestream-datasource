package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	timestreamquerytypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/timestream-datasource/pkg/models"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeSender struct {
	res *backend.CallResourceResponse
}

func (f *fakeSender) Send(res *backend.CallResourceResponse) error {
	f.res = res
	return nil
}

type fakeClient struct {
	output *timestreamquery.QueryOutput

	calls runnerCalls
}

type runnerCalls struct {
	runQuery []*timestreamquery.QueryInput
}

func (f *fakeClient) Query(_ context.Context, input *timestreamquery.QueryInput, _ ...func(*timestreamquery.Options)) (*timestreamquery.QueryOutput, error) {
	f.calls.runQuery = append(f.calls.runQuery, input)
	return f.output, nil
}

func (f *fakeClient) CancelQuery(context.Context, *timestreamquery.CancelQueryInput, ...func(*timestreamquery.Options)) (*timestreamquery.CancelQueryOutput, error) {
	return nil, nil
}

func TestCallResource(t *testing.T) {
	tests := []struct {
		description string
		output      *timestreamquery.QueryOutput
		req         *backend.CallResourceRequest
		result      string
	}{
		{
			"databases request",
			&timestreamquery.QueryOutput{
				Rows: []timestreamquerytypes.Row{
					{Data: []timestreamquerytypes.Datum{{ScalarValue: aws.String("foo")}}},
					{Data: []timestreamquerytypes.Datum{{ScalarValue: aws.String("bar")}}},
				},
			},
			&backend.CallResourceRequest{
				Path: "databases",
			},
			`["\"foo\"","\"bar\""]`,
		},
		{
			"tables request",
			&timestreamquery.QueryOutput{
				Rows: []timestreamquerytypes.Row{
					{Data: []timestreamquerytypes.Datum{{ScalarValue: aws.String("foo")}}},
					{Data: []timestreamquerytypes.Datum{{ScalarValue: aws.String("bar")}}},
				},
			},
			&backend.CallResourceRequest{
				Method: "POST",
				Path:   "tables",
				Body:   []byte(`{"database":"db"}`),
			},
			`["\"foo\"","\"bar\""]`,
		},
		{
			"measures request",
			&timestreamquery.QueryOutput{
				Rows: []timestreamquerytypes.Row{
					{Data: []timestreamquerytypes.Datum{{ScalarValue: aws.String("foo")}}},
					{Data: []timestreamquerytypes.Datum{{ScalarValue: aws.String("bar")}}},
				},
			},
			&backend.CallResourceRequest{
				Method: "POST",
				Path:   "measures",
				Body:   []byte(`{"database":"db","table":"t"}`),
			},
			`["foo","bar"]`,
		},
		{
			"dimensions request",
			&timestreamquery.QueryOutput{
				Rows: []timestreamquerytypes.Row{
					{Data: []timestreamquerytypes.Datum{
						{}, // measure
						{}, // measure type
						{ArrayValue: []timestreamquerytypes.Datum{ // dimensions
							{RowValue: &timestreamquerytypes.Row{Data: []timestreamquerytypes.Datum{
								{ScalarValue: aws.String("foo")},
							}}},
							{RowValue: &timestreamquerytypes.Row{Data: []timestreamquerytypes.Datum{
								{ScalarValue: aws.String("bar")},
							}}},
						}},
					}},
				},
			},
			&backend.CallResourceRequest{
				Method: "POST",
				Path:   "dimensions",
				Body:   []byte(`{"database":"db","table":"t"}`),
			},
			`["foo","bar"]`,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ts := &timestreamDS{Client: &fakeClient{output: test.output}}
			sender := &fakeSender{}
			err := ts.CallResource(context.Background(), test.req, sender)
			if err != nil {
				t.Error(err)
			}
			if string(sender.res.Body) != test.result {
				t.Errorf("unexpected result %s", string(sender.res.Body))
			}
		})
	}
}

func Test_runQuery_always_wraps_db_and_table_name_in_quotes(t *testing.T) {
	testCases := []struct {
		name, resource, requestBody, expectedQuery string
	}{
		{
			resource:      "tables",
			requestBody:   `{"database":"db"}`,
			expectedQuery: `SHOW TABLES FROM "db"`,
		},
		{
			resource:      "tables",
			requestBody:   `{"database":"\"db\""}`,
			expectedQuery: `SHOW TABLES FROM "db"`,
		},
		{
			resource:      "measures",
			requestBody:   `{"database":"db","table":"some_table_name"}`,
			expectedQuery: `SHOW MEASURES FROM "db"."some_table_name"`,
		},
		{
			resource:      "measures",
			requestBody:   `{"database":"\"db\"","table":"\"some_table_name\""}`,
			expectedQuery: `SHOW MEASURES FROM "db"."some_table_name"`,
		},
		{
			resource:      "dimensions",
			requestBody:   `{"database":"db","table":"some_table_name"}`,
			expectedQuery: `SHOW MEASURES FROM "db"."some_table_name"`,
		},
		{
			resource:      "dimensions",
			requestBody:   `{"database":"\"db\"","table":"\"some_table_name\""}`,
			expectedQuery: `SHOW MEASURES FROM "db"."some_table_name"`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			client := &fakeClient{output: &timestreamquery.QueryOutput{Rows: []timestreamquerytypes.Row{}}}
			ds := &timestreamDS{Client: client}

			assert.NoError(t, ds.CallResource(context.Background(),
				&backend.CallResourceRequest{
					Method: "POST",
					Path:   test.resource,
					Body:   []byte(test.requestBody),
				}, &fakeSender{}))

			require.Len(t, client.calls.runQuery, 1)
			assert.Equal(t, &timestreamquery.QueryInput{QueryString: &test.expectedQuery}, client.calls.runQuery[0])
		})
	}
}

// The following were formerly in executor_test.go

func runTest(t *testing.T, names []string) *backend.DataResponse {
	mockClient := &MockClient{testFileNames: names}
	ds := timestreamDS{Client: mockClient}
	dr := ds.ExecuteQuery(context.Background(), models.QueryModel{WaitForResult: true})

	// Remove changeable fields
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
	experimental.CheckGoldenJSONResponse(t, "./testdata", names[0], &dr, false)

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
	runTest(t, []string{"time-series-with-null-data-points"})
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

	m["select-null-timestamp.json"] = models.QueryModel{
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

	m["time-series-with-null-data-points.json"] = models.QueryModel{
		RawQuery: `WITH
			binnedTimeseries AS (
				SELECT BIN(time, 10000ms) AS binnedTime, SUM(measure_value::double) as value
				FROM ` + table + ` WHERE time BETWEEN from_milliseconds(1654008860795) AND from_milliseconds(1654030460795) 
					AND measure_name = 'any_metric_name'
				GROUP BY BIN(time, 10000ms)
				ORDER BY binnedTime
			),
			interpolatedTimeseries AS (
				SELECT 
					INTERPOLATE_FILL(
						CREATE_TIME_SERIES(binnedTime, value),
						SEQUENCE(min(binnedTime), max(binnedTime), 10000ms),
						sqrt(-1)
					) AS interpolatedValue
				FROM binnedTimeseries
			)
			SELECT * FROM interpolatedTimeseries`,
	}

	nextToken := map[string]*string{}
	for filename, query := range m {
		inst, err := NewDatasource(context.Background(), backend.DataSourceInstanceSettings{
			JSONData: []byte(`{"region": "us-west-2"}`),
		})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		ds := inst.(*timestreamDS)
		raw, _ := Interpolate(query, models.DatasourceSettings{})
		input := &timestreamquery.QueryInput{
			QueryString: aws.String(raw),
		}
		// Custom input to test pagination
		if filename == "pagination-off_1.json" {
			// force pagination
			input.MaxRows = aws.Int32(1)
		}
		if filename == "pagination-off_2.json" {
			// continue previous page
			input.NextToken = nextToken["pagination-off_1.json"]
			if input.NextToken == nil {
				t.Fatalf("first page should be executed first. Please retry")
			}
		}

		res, err := ds.Client.Query(context.Background(), input)
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
