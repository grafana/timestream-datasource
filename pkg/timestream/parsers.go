package timestream

import (
	"fmt"
	"strconv"
	"time"

	timestreamquerytypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type datumParser func(datum timestreamquerytypes.Datum) (interface{}, error)

type fieldBuilder struct {
	name       string
	columnIdx  int
	fieldType  data.FieldType
	config     *data.FieldConfig
	parser     datumParser
	asJSON     bool // if true, the results will be marshaled to json first
	timeseries bool
}

func getFieldBuilder(t *timestreamquerytypes.Type) (*fieldBuilder, error) {
	if t.ScalarType != "" {
		switch t.ScalarType {
		case timestreamquerytypes.ScalarTypeTimestamp:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableTime,
				parser:    datumParserTimestamp,
			}, nil
		case timestreamquerytypes.ScalarTypeBoolean:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableBool,
				parser:    datumParserBool,
			}, nil
		case timestreamquerytypes.ScalarTypeVarchar:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableString,
				parser:    datumParserString,
			}, nil
		case timestreamquerytypes.ScalarTypeDouble:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableFloat64,
				parser:    datumParserFloat64,
			}, nil
		case timestreamquerytypes.ScalarTypeBigint:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableInt64,
				parser:    datumParserInt64,
			}, nil

		case timestreamquerytypes.ScalarTypeInteger:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableInt32,
				parser:    datumParserInt32,
			}, nil

		case timestreamquerytypes.ScalarTypeIntervalDayToSecond:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableString,
				parser:    datumParserInterval,
			}, nil

		case timestreamquerytypes.ScalarTypeIntervalYearToMonth:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableString,
				parser:    datumParserInterval,
			}, nil

		case timestreamquerytypes.ScalarTypeDate:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableTime,
				parser:    datumParserDate,
			}, nil

		case timestreamquerytypes.ScalarTypeTime:
			return &fieldBuilder{
				fieldType: data.FieldTypeNullableTime,
				parser:    datumParserTime,
			}, nil

		default:
			return nil, fmt.Errorf("unsupported scalar value: %s", t.ScalarType)
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

	return nil, fmt.Errorf("unsupported column: %+v", t)
}

func getArrayBuilder(column *timestreamquerytypes.ColumnInfo) (*fieldBuilder, error) {
	elem, err := getFieldBuilder(column.Type)
	if err != nil {
		return nil, err
	}

	parser := func(datum timestreamquerytypes.Datum) (interface{}, error) {
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

func getRowBuilder(columns []timestreamquerytypes.ColumnInfo) (*fieldBuilder, error) {
	count := len(columns)
	cols := make([]*fieldBuilder, count)
	for i := 0; i < len(columns); i++ {
		elem, err := getFieldBuilder(columns[i].Type)
		if err != nil {
			return nil, err
		}
		cols[i] = elem
	}

	parser := func(datum timestreamquerytypes.Datum) (interface{}, error) {
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

func datumParserBool(datum timestreamquerytypes.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := strconv.ParseBool(*datum.ScalarValue)
	return &v, err
}

func datumParserInt32(datum timestreamquerytypes.Datum) (interface{}, error) {
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

func datumParserInt64(datum timestreamquerytypes.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := strconv.ParseInt(*datum.ScalarValue, 10, 64)
	return &v, err
}

func datumParserFloat64(datum timestreamquerytypes.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := strconv.ParseFloat(*datum.ScalarValue, 64)
	return &v, err
}

func datumParserTimestamp(datum timestreamquerytypes.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := time.Parse("2006-01-02 15:04:05.99999999", *datum.ScalarValue)
	return &v, err
}

func datumParserDate(datum timestreamquerytypes.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := time.Parse("2006-01-02", *datum.ScalarValue)
	return &v, err
}

func datumParserTime(datum timestreamquerytypes.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	v, err := time.Parse("15:04:05.99999999", *datum.ScalarValue)
	if err != nil {
		return nil, err
	}
	// the default is that parse will use year 0 which will not display properly
	properTime := v.AddDate(1970, 0, 0)
	return &properTime, nil
}

func datumParserString(datum timestreamquerytypes.Datum) (interface{}, error) {
	return datum.ScalarValue, nil
}

func datumParserInterval(datum timestreamquerytypes.Datum) (interface{}, error) {
	if datum.ScalarValue == nil {
		return nil, nil
	}
	// TODO: parse into a better datatype, maybe an integer for millisecond?
	// Right now this string is consistent with Timestream console
	return datum.ScalarValue, nil
}
