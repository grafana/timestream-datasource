package timestream

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
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

func Test_queries_called_with_db_name_quoted(t *testing.T) {
	t.Run("db name provided without quotes runs query with quoted db name", func(t *testing.T) {
		mockRunner := &fakeRunner{output: &timestreamquery.QueryOutput{Rows: []*timestreamquery.Row{}}}
		ts := &timestreamDS{Runner: mockRunner}
		sender := &fakeSender{}

		assert.NoError(t, ts.CallResource(context.Background(), &backend.CallResourceRequest{
			Method: "POST",
			Path:   "tables",
			Body:   []byte(`{"database":"db"}`),
		}, sender))

		require.Len(t, mockRunner.calls.runQuery, 1)
		assert.Equal(t, &timestreamquery.QueryInput{QueryString: aws.String(`SHOW TABLES FROM "db"`)}, mockRunner.calls.runQuery[0])
	})

	t.Run("db name provided with quotes runs query with db name as-is", func(t *testing.T) {
		mockRunner := &fakeRunner{output: &timestreamquery.QueryOutput{Rows: []*timestreamquery.Row{}}}
		ts := &timestreamDS{Runner: mockRunner}
		sender := &fakeSender{}

		assert.NoError(t, ts.CallResource(context.Background(), &backend.CallResourceRequest{
			Method: "POST",
			Path:   "tables",
			Body:   []byte(`{"database":"\"db\""}`),
		}, sender))

		require.Len(t, mockRunner.calls.runQuery, 1)
		assert.Equal(t, &timestreamquery.QueryInput{QueryString: aws.String(`SHOW TABLES FROM "db"`)}, mockRunner.calls.runQuery[0])
	})

}
