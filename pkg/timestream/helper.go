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
	return strconv.ParseBool(datum.ScalarValue)
}

func datumParserInt(datum *timestreamquery.Datum) (interface{}, error) {
	return strconv.ParseInt(datum.ScalarValue)
}

func datumParserFloat(datum *timestreamquery.Datum) (interface{}, error) {
	return strconv.ParseFloat(datum.ScalarValue)
}

func datumParserTime(datum *timestreamquery.Datum) (interface{}, error) {
	return time.Parse(datum.ScalarValue)
}

type columnFiller struct {
	index  int16
	field  *data.Field
	parser datumParser
}

// QueryResultToDataFrame creates a DataFrame from query results
func QueryResultToDataFrame(res *timestreamquery.QueryOutput) (*data.Frame, error) {
	fields := []*data.Field{}
	fillers := []columnFiller{}
	warnings := []string{}

	for index, columnMeta := range res.ColumnInfo {
		if columnMeta.Type.ScalarType != nil {
			var field *data.Field
			var parser datumParser

			switch columnMeta.Type.ScalarType {
			case timestreamquery.ScalarTypeBoolean:
				field = data.NewField(*columnMeta.Name, nil, make([]*time.Time, 10))
				parser = datumParserBool
			default:
				warnings = append(warnings, fmt.Sprintf("Unsupported scalar value: %s", columnMeta.Type.ScalarType))
			}

			if field != nil {
				fillers = append(fillers, columnFiller{
					index:  index,
					field:  field,
					parser: parser,
				})
				fields = append(fields, field)
			}
		} else {
			warnings = append(warnings, fmt.Sprintf("Unsupported column type", columnMeta.Type.String()))
		}
	}

	frame := data.NewFrame(refID, // TODO: shoud set the name from metadata
		fields...,
	)

	return -1
}
