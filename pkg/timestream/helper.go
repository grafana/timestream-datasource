package timestream

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type datumParser func(datum *timestreamquery.Datum) (interface{}, error)

func datumParserBool(datum *timestreamquery.Datum) (interface{}, error) {
	return strconv.ParseBool(*datum.ScalarValue)
}

func datumParserInt(datum *timestreamquery.Datum) (interface{}, error) {
	return strconv.Atoi(*datum.ScalarValue)
}

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

// QueryResultToDataFrame creates a DataFrame from query results
func QueryResultToDataFrame(res *timestreamquery.QueryOutput) (*data.Frame, error) {
	fields := []*data.Field{}
	warnings := []string{}

	length := len(res.Rows)
	for index, columnMeta := range res.ColumnInfo {
		if columnMeta.Type.ScalarType != nil {
			var field *data.Field
			var parser datumParser

			switch *columnMeta.Type.ScalarType {
			case timestreamquery.ScalarTypeTimestamp:
				field = data.NewField(*columnMeta.Name, nil, make([]time.Time, length))
				parser = datumParserTime
			case timestreamquery.ScalarTypeBoolean:
				field = data.NewField(*columnMeta.Name, nil, make([]bool, length))
				parser = datumParserBool
			case timestreamquery.ScalarTypeVarchar:
				field = data.NewField(*columnMeta.Name, nil, make([]string, length))
				parser = datumParserString
			case timestreamquery.ScalarTypeDouble:
				field = data.NewField(*columnMeta.Name, nil, make([]float64, length))
				parser = datumParserFloat64
			default:
				warnings = append(warnings, fmt.Sprintf("Unsupported scalar value: %s", *columnMeta.Type.ScalarType))
			}

			if field != nil {
				// Read all the values
				for i := 0; i < length; i++ {
					row := res.Rows[i]
					if row == nil {
						continue
					}
					v, err := parser(row.Data[index])
					if err != nil {
						warnings = append(warnings, fmt.Sprintf("Error parsing: row:%d, colum:%d", i, index))
					} else {
						field.Set(i, v)
					}
				}
				fields = append(fields, field)
			}
		} else {
			warnings = append(warnings, fmt.Sprintf("Unsupported column type: %s", columnMeta.Type.GoString()))
		}
	}

	frame := data.NewFrame("result", // No name
		fields...,
	)
	return frame, nil
}
