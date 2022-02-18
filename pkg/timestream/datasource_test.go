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
	resources []string
}

func (f *fakeRunner) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	res := &timestreamquery.QueryOutput{}
	dimensions := []*timestreamquery.Datum{}
	for _, r := range f.resources {
		res.Rows = append(res.Rows, &timestreamquery.Row{Data: []*timestreamquery.Datum{
			{ScalarValue: aws.String(r)},
			{},
			{}, // Dimension data
		}})
		// Populate dimensions using same resources
		dimensions = append(dimensions, &timestreamquery.Datum{
			RowValue: &timestreamquery.Row{
				Data: []*timestreamquery.Datum{{ScalarValue: aws.String(r)}},
			},
		})
	}
	// Only populate dimensions data in the first row because it's the only one we are reading
	res.Rows[0].Data[2].ArrayValue = dimensions
	return res, nil
}

func TestCallResource(t *testing.T) {
	tests := []struct {
		description string
		resources   []string
		req         *backend.CallResourceRequest
		result      string
	}{
		{
			"databases request",
			[]string{"foo", "bar"},
			&backend.CallResourceRequest{
				Path: "databases",
			},
			`["\"foo\"","\"bar\""]`,
		},
		{
			"tables request",
			[]string{"foo", "bar"},
			&backend.CallResourceRequest{
				Method: "POST",
				Path:   "tables",
				Body:   []byte(`{"database":"db"}`),
			},
			`["\"foo\"","\"bar\""]`,
		},
		{
			"measures request",
			[]string{"foo", "bar"},
			&backend.CallResourceRequest{
				Method: "POST",
				Path:   "measures",
				Body:   []byte(`{"database":"db","table":"t"}`),
			},
			`["foo","bar"]`,
		},
		{
			"dimensions request",
			[]string{"foo", "bar"},
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
			ts := &timestreamDS{Runner: &fakeRunner{resources: test.resources}}
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
