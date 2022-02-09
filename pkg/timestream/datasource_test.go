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
	for _, r := range f.resources {
		res.Rows = append(res.Rows, &timestreamquery.Row{Data: []*timestreamquery.Datum{
			{ScalarValue: aws.String(r)},
		}})
	}
	return res, nil
}

func TestCallResource(t *testing.T) {
	tests := []struct {
		description string
		req         *backend.CallResourceRequest
	}{
		{
			"databases request",
			&backend.CallResourceRequest{
				Path: "databases",
			},
		},
		{
			"tables request",
			&backend.CallResourceRequest{
				Method: "POST",
				Path:   "tables",
				Body:   []byte(`{"database":"db"}`),
			},
		},
		{
			"databases request",
			&backend.CallResourceRequest{
				Method: "POST",
				Path:   "tables",
				Body:   []byte(`{"database":"db","table":"t"}`),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ts := &timestreamDS{Runner: &fakeRunner{resources: []string{"foo", "bar"}}}
			sender := &fakeSender{}
			err := ts.CallResource(context.Background(), test.req, sender)
			if err != nil {
				t.Error(err)
			}
			if string(sender.res.Body) != `["foo","bar"]` {
				t.Errorf("unexpected result %s", string(sender.res.Body))
			}
		})
	}
}
