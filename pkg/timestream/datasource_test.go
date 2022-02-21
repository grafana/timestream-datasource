package timestream

import (
	"context"
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
}

func (f *fakeRunner) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
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
