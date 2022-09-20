package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/timestreamquery"
)

// MockClient ...
type MockClient struct {
	testFileNames []string
	index         int
}

func (c *MockClient) runQuery(ctx context.Context, input *timestreamquery.QueryInput) (*timestreamquery.QueryOutput, error) {
	bs, err := os.ReadFile("./testdata/" + c.testFileNames[c.index] + ".json")
	if err != nil {
		return nil, err
	}
	r := &timestreamquery.QueryOutput{}
	err = json.Unmarshal(bs, r)
	if err != nil {
		fmt.Println(err)
	}
	c.index++
	return r, nil
}

func (c *MockClient) cancelQuery(ctx context.Context, input *timestreamquery.CancelQueryInput) (*timestreamquery.CancelQueryOutput, error) {
	r := &timestreamquery.CancelQueryOutput{}
	return r, nil
}
