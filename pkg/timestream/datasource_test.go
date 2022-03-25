package timestream

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
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

type fakeRunner struct {
	queryRunner
	output *timestreamquery.QueryOutput

	calls runnerCalls
}

type runnerCalls struct {
	runQuery []*timestreamquery.QueryInput
}

func (f *fakeRunner) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	f.calls.runQuery = append(f.calls.runQuery, input)

	return f.output, nil
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
				Rows: []*timestreamquery.Row{
					{Data: []*timestreamquery.Datum{{ScalarValue: aws.String("foo")}}},
					{Data: []*timestreamquery.Datum{{ScalarValue: aws.String("bar")}}},
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
				Rows: []*timestreamquery.Row{
					{Data: []*timestreamquery.Datum{{ScalarValue: aws.String("foo")}}},
					{Data: []*timestreamquery.Datum{{ScalarValue: aws.String("bar")}}},
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
				Rows: []*timestreamquery.Row{
					{Data: []*timestreamquery.Datum{{ScalarValue: aws.String("foo")}}},
					{Data: []*timestreamquery.Datum{{ScalarValue: aws.String("bar")}}},
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
				Rows: []*timestreamquery.Row{
					{Data: []*timestreamquery.Datum{
						{}, // measure
						{}, // measure type
						{ArrayValue: []*timestreamquery.Datum{ // dimensions
							{RowValue: &timestreamquery.Row{Data: []*timestreamquery.Datum{
								{ScalarValue: aws.String("foo")},
							}}},
							{RowValue: &timestreamquery.Row{Data: []*timestreamquery.Datum{
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
			ts := &timestreamDS{Runner: &fakeRunner{output: test.output}}
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

func Test_runQuery_always_sends_db_name_with_quotes(t *testing.T) {
	testCases := []struct {
		name, resource, requestBody, expectedQuery string
	}{
		{
			name:          "tables: db name without quotes runs query with quoted db name",
			resource:      "tables",
			requestBody:   `{"database":"db"}`,
			expectedQuery: `SHOW TABLES FROM "db"`,
		},
		{
			name:          "tables: db name with quotes runs query with db name as-is",
			resource:      "tables",
			requestBody:   `{"database":"\"db\""}`,
			expectedQuery: `SHOW TABLES FROM "db"`,
		},
		{
			name:          "measures: db name without quotes runs query with quoted db name",
			resource:      "measures",
			requestBody:   `{"database":"db","table":"some_table_name"}`,
			expectedQuery: `SHOW MEASURES FROM "db".some_table_name`,
		},
		{
			name:          "measures: db name with quotes runs query with db name as-is",
			resource:      "measures",
			requestBody:   `{"database":"\"db\"","table":"some_table_name"}`,
			expectedQuery: `SHOW MEASURES FROM "db".some_table_name`,
		},
		{
			name:          "dimensions: db name without quotes runs query with quoted db name",
			resource:      "dimensions",
			requestBody:   `{"database":"db","table":"some_table_name"}`,
			expectedQuery: `SHOW MEASURES FROM "db".some_table_name`,
		},
		{
			name:          "dimensions: db name with quotes runs query with db name as-is",
			resource:      "dimensions",
			requestBody:   `{"database":"\"db\"","table":"some_table_name"}`,
			expectedQuery: `SHOW MEASURES FROM "db".some_table_name`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockRunner := &fakeRunner{output: &timestreamquery.QueryOutput{Rows: []*timestreamquery.Row{}}}

			assert.NoError(t, (&timestreamDS{Runner: mockRunner}).CallResource(context.Background(),
				&backend.CallResourceRequest{
					Method: "POST",
					Path:   test.resource,
					Body:   []byte(test.requestBody),
				}, &fakeSender{}))

			require.Len(t, mockRunner.calls.runQuery, 1)
			assert.Equal(t, &timestreamquery.QueryInput{QueryString: &test.expectedQuery}, mockRunner.calls.runQuery[0])
		})
	}
}
