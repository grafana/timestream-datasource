package timestream

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
)

// MockClient ...
type MockClient struct {
	testFileNames []string
	index         int
}

func (c *MockClient) Query(context.Context, *timestreamquery.QueryInput, ...func(options *timestreamquery.Options)) (*timestreamquery.QueryOutput, error) {
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

func (c *MockClient) CancelQuery(context.Context, *timestreamquery.CancelQueryInput, ...func(options *timestreamquery.Options)) (*timestreamquery.CancelQueryOutput, error) {
	r := &timestreamquery.CancelQueryOutput{}
	return r, nil
}
