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
	config     *data.FieldConfig
	parser     datumParser
	asJSON     bool // if true, the results will be marshaled to json first
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
				fieldType: data.FieldTypeNullableBool,
				parser:    datumParserBool,
			}, nil
		case timestreamquery.ScalarTypeVarchar:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableString,
				parser:    datumParserString,
			}, nil
		case timestreamquery.ScalarTypeDouble:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableFloat64,
				parser:    datumParserFloat64,
			}, nil
		case timestreamquery.ScalarTypeBigint:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableInt64,
				parser:    datumParserInt64,
			}, nil

		case timestreamquery.ScalarTypeInteger:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableInt32,
				parser:    datumParserInt32,
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

	if t.RowColumnInfo != nil {
		return getRowBuilder(t.RowColumnInfo)
	}

	if t.ArrayColumnInfo != nil {
		return getArrayBuilder(t.ArrayColumnInfo)
	}

	return nil, fmt.Errorf("Unsupported column: %s", t.GoString())
}

func getArrayBuilder(column *timestreamquery.ColumnInfo) (*fieldBuilder, error) {
	elem, err := getFieldBuilder(column.Type)
	if err != nil {
		return nil, err
	}

	parser := func(datum *timestreamquery.Datum) (interface{}, error) {
		count := len(datum.ArrayValue)
		vals := make([]interface{}, count)
		for i, d := range datum.ArrayValue {
			v, err := elem.parser(d)
			if err != nil {
				return nil, err
			}
			vals[i] = v
		}
		return vals, nil
	}

	tableProps := make(map[string]interface{})
	tableProps["displayMode"] = "json-view"
	return &fieldBuilder{
		fieldType: data.FieldTypeString,
		parser:    parser,
		asJSON:    true,
		config: &data.FieldConfig{
			Custom: tableProps,
		},
	}, nil
}

func getRowBuilder(columns []*timestreamquery.ColumnInfo) (*fieldBuilder, error) {
	count := len(columns)
	cols := make([]*fieldBuilder, count)
	for i := 0; i < len(columns); i++ {
		elem, err := getFieldBuilder(columns[i].Type)
		if err != nil {
			return nil, err
		}
		cols[i] = elem
	}

	parser := func(datum *timestreamquery.Datum) (interface{}, error) {
		vals := make(map[string]interface{})
		for i, d := range datum.RowValue.Data {
			v, err := cols[i].parser(d)
			if err != nil {
				return nil, err
			}
			vals[*columns[i].Name] = v
		}
		return vals, nil
	}

	tableProps := make(map[string]interface{})
	return &fieldBuilder{
		fieldType: data.FieldTypeString,
		parser:    parser,
		asJSON:    true,
		config: &data.FieldConfig{
			Custom: tableProps,
		},
	}, nil
}

//---------------------------------------------------

func datumParserBool(datum *timestreamquery.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := strconv.ParseBool(*datum.ScalarValue)
	return &v, err
}

func datumParserInt32(datum *timestreamquery.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	i64, err := strconv.ParseInt(*datum.ScalarValue, 10, 32)
	if err != nil {
		return nil, err
	}
	i32 := int32(i64)
	return &i32, nil
}

func datumParserInt64(datum *timestreamquery.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := strconv.ParseInt(*datum.ScalarValue, 10, 64)
	return &v, err
}

func datumParserFloat64(datum *timestreamquery.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := strconv.ParseFloat(*datum.ScalarValue, 0)
	return &v, err
}

func datumParserTime(datum *timestreamquery.Datum) (interface{}, error) {
	// 2020-03-18 17:26:30.00000000
	// 2006-01-02 15:04:05.99999999
	return time.Parse("2006-01-02 15:04:05.99999999", *datum.ScalarValue)
}

func datumParserString(datum *timestreamquery.Datum) (interface{}, error) {
	return datum.ScalarValue, nil
}
