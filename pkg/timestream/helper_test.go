package timestream

import (
	timestreamquerytypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	"github.com/grafana/timestream-datasource/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestQueryResultToDataFrame(t *testing.T) {
	input := &timestreamquery.QueryOutput{
		ColumnInfo: []timestreamquerytypes.ColumnInfo{
			{
				Name: aws.String("time"),
				Type: &timestreamquerytypes.Type{
					ScalarType: "TIMESTAMP",
				},
			},
			{
				Name: aws.String("instance_name"),
				Type: &timestreamquerytypes.Type{
					ScalarType: "VARCHAR",
				},
			},
			{
				Name: aws.String("microservice_name"),
				Type: &timestreamquerytypes.Type{
					ScalarType: "VARCHAR",
				},
			},
			{
				Name: aws.String("value"),
				Type: &timestreamquerytypes.Type{
					ScalarType: "DOUBLE",
				},
			},
		},
		Rows: []timestreamquerytypes.Row{
			{
				Data: []timestreamquerytypes.Datum{
					{ScalarValue: aws.String("2021-03-14 09:52:44.000000000")},
					{ScalarValue: aws.String("instance-1.amazonaws.com")},
					{ScalarValue: aws.String("zeus")},
					{ScalarValue: aws.String("1.2")},
				},
			},
			{
				Data: []timestreamquerytypes.Datum{
					{ScalarValue: aws.String("2021-03-14 09:52:44.000000000")},
					{ScalarValue: aws.String("instance-1.amazonaws.com")},
					{ScalarValue: aws.String("apollo")},
					{ScalarValue: aws.String("1.3")},
				},
			},
			{
				Data: []timestreamquerytypes.Datum{
					{ScalarValue: aws.String("2021-03-14 09:57:44.000000000")},
					{ScalarValue: aws.String("instance-1.amazonaws.com")},
					{ScalarValue: aws.String("zeus")},
					{ScalarValue: aws.String("2.0")},
				},
			},
			{
				Data: []timestreamquerytypes.Datum{
					{ScalarValue: aws.String("2021-03-14 09:57:44.000000000")},
					{ScalarValue: aws.String("instance-1.amazonaws.com")},
					{ScalarValue: aws.String("apollo")},
					{ScalarValue: aws.String("1.5")},
				},
			},
		},
	}

	t.Run("table format", func(t *testing.T) {
		res := QueryResultToDataFrame(input, models.FormatOptionTable)

		// Assert that it returns one frame with four fields
		assert.Equal(t, 1, len(res.Frames))
		assert.Equal(t, 4, len(res.Frames[0].Fields))
		assert.Equal(t, "time", res.Frames[0].Fields[0].Name)
		assert.Equal(t, "instance_name", res.Frames[0].Fields[1].Name)
		assert.Equal(t, "microservice_name", res.Frames[0].Fields[2].Name)
		assert.Equal(t, "value", res.Frames[0].Fields[3].Name)
	})

	t.Run("timeseries format", func(t *testing.T) {
		res := QueryResultToDataFrame(input, models.FormatOptionTimeSeries)
		// Assert that it returns one frame with three fields
		assert.Equal(t, 1, len(res.Frames))
		assert.Equal(t, 3, len(res.Frames[0].Fields))
		assert.Equal(t, "time", res.Frames[0].Fields[0].Name)

		// And each field represents a time series
		assert.Equal(t, "value", res.Frames[0].Fields[1].Name)
		assert.Equal(t, 2, len(res.Frames[0].Fields[1].Labels))
		assert.Equal(t, "instance-1.amazonaws.com", res.Frames[0].Fields[1].Labels["instance_name"])
		assert.Equal(t, "apollo", res.Frames[0].Fields[1].Labels["microservice_name"])

		assert.Equal(t, "value", res.Frames[0].Fields[2].Name)
		assert.Equal(t, 2, len(res.Frames[0].Fields[2].Labels))
		assert.Equal(t, "instance-1.amazonaws.com", res.Frames[0].Fields[2].Labels["instance_name"])
		assert.Equal(t, "zeus", res.Frames[0].Fields[2].Labels["microservice_name"])
	})
	t.Run("timeseries format with no rows", func(t *testing.T) {
		input.Rows = []timestreamquerytypes.Row{}
		inputWithNoRows := input
		inputWithNoRows.Rows = []timestreamquerytypes.Row{}
		res := QueryResultToDataFrame(inputWithNoRows, models.FormatOptionTimeSeries)
		// Assert that it returns one frame with no fields
		assert.Equal(t, 1, len(res.Frames))
		assert.Equal(t, 4, len(res.Frames[0].Fields))
		assert.Equal(t, 0, res.Frames[0].Fields[0].Len())
	})
}
