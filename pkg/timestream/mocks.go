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

func (c *MockClient) readText() (string, error) {
	bs, err := ioutil.ReadFile("./testdata/" + c.testFileName + ".txt")
	if err != nil {
		return "", err
	}
	txt := string(bs)
	return txt, nil
}

func (c *MockClient) saveText(text string) error {
	return ioutil.WriteFile("./testdata/"+c.testFileName+".txt", []byte(text), 0644)
}
