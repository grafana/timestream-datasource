package timestream

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type datumParser func(datum *timestreamquery.Datum) (interface{}, error)

type fieldBuilder struct {
	name       string
	columnIdx  int
	fieldType  data.FieldType
	parser     datumParser
	timeseries bool
}

func getFieldBuilder(t *timestreamquery.Type) (*fieldBuilder, error) {
	if t.ScalarType != nil {
		switch *t.ScalarType {
		case timestreamquery.ScalarTypeTimestamp:
			return &fieldBuilder{
				fieldType: data.FieldTypeTime,
				parser:    datumParserTime,
			}, nil
		case timestreamquery.ScalarTypeBoolean:
			return &fieldBuilder{
				fieldType: data.FieldTypeBool,
				parser:    datumParserBool,
			}, nil
		case timestreamquery.ScalarTypeVarchar:
			return &fieldBuilder{
				fieldType: data.FieldTypeString,
				parser:    datumParserString,
			}, nil
		case timestreamquery.ScalarTypeDouble:
			return &fieldBuilder{
				fieldType: data.FieldTypeFloat64,
				parser:    datumParserFloat64,
			}, nil

		default:
			return nil, fmt.Errorf("Unsupported scalar value: %s", *t.ScalarType)
		}
	}

	if t.TimeSeriesMeasureValueColumnInfo != nil {
		builder, err := getFieldBuilder(t.TimeSeriesMeasureValueColumnInfo.Type)
		if err != nil {
			return nil, err
		}
		builder.timeseries = true
		return builder, nil
	}

	return nil, fmt.Errorf("Unsupported column: %s", t.GoString())
}

func datumParserBool(datum *timestreamquery.Datum) (interface{}, error) {
	return strconv.ParseBool(*datum.ScalarValue)
}

// func datumParserInt(datum *timestreamquery.Datum) (interface{}, error) {
// 	return strconv.Atoi(*datum.ScalarValue)
// }

func datumParserFloat64(datum *timestreamquery.Datum) (interface{}, error) {
	return strconv.ParseFloat(*datum.ScalarValue, 0)
}

func datumParserTime(datum *timestreamquery.Datum) (interface{}, error) {
	// 2020-03-18 17:26:30.00000000
	// 2006-01-02 15:04:05.99999999
	return time.Parse("2006-01-02 15:04:05.99999999", *datum.ScalarValue)
}

func datumParserString(datum *timestreamquery.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return "", nil // or error?
	}
	return *datum.ScalarValue, nil
}
