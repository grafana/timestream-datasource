package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

// MockClient ...
type MockClient struct {
	testFileName string
}

func (c *MockClient) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	bs, err := ioutil.ReadFile("./testdata/" + c.testFileName + ".json")
	if err != nil {
		return nil, err
	}
	r := &timestreamquery.QueryOutput{}
	err = json.Unmarshal(bs, r)
	if err != nil {
		fmt.Println(err)
	}
	return r, nil
}

func (c *MockClient) cancelQuery(ctx context.Context, input *timestreamquery.CancelQueryInput) (*timestreamquery.CancelQueryOutput, error) {
	r := &timestreamquery.CancelQueryOutput{}
	return r, nil
}
