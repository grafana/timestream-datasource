package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"

	"fmt"
)

// TryQueryRunner runs a query
func main() {
	// setup the query client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic("dooh")
	}

	querySvc := timestreamquery.New(sess)

	database := "devops"
	table := "host_metrics"
	hostName := "host-l9M6E"

	// Example #1
	queryStr := `SELECT region, az, hostname, BIN(time, 15s) AS binned_timestamp,
        ROUND(AVG(measure_value), 2) AS avg_cpu_utilization,
        ROUND(APPROX_PERCENTILE(measure_value, 0.9), 2) AS p90_cpu_utilization,
        ROUND(APPROX_PERCENTILE(measure_value, 0.95), 2) AS p95_cpu_utilization,
        ROUND(APPROX_PERCENTILE(measure_value, 0.99), 2) AS p99_cpu_utilization
    FROM ` + database + "." + table + `
    WHERE measure_name = 'cpu_utilization'
    AND hostname = '` + hostName + `'
        AND time > ago(5d)
    GROUP BY region, hostname, az, BIN(time, 15s)
    ORDER BY binned_timestamp ASC
    LIMIT 10
	`

	// // Example #4
	// queryStr = `WITH binned_timeseries AS (
	// 	SELECT hostname, BIN(time, 30s) AS binned_timestamp, ROUND(AVG(measure_value), 2) AS avg_cpu_utilization
	//   FROM ` + database + "." + table + `
	// 	WHERE measure_name = 'cpu_utilization'
	// 	  AND hostname = '` + hostName + `'
	// 		AND time > ago(2h)
	// 	GROUP BY hostname, BIN(time, 30s)
	// ), interpolated_timeseries AS (
	// 	SELECT hostname,
	// 		INTERPOLATE_LINEAR(
	// 			CREATE_TIME_SERIES(binned_timestamp, avg_cpu_utilization),
	// 				SEQUENCE(min(binned_timestamp), max(binned_timestamp), 15s)) AS interpolated_avg_cpu_utilization
	// 	FROM binned_timeseries
	// 	GROUP BY '` + hostName + `'
	// )
	// SELECT time, ROUND(value, 2) AS interpolated_cpu
	// FROM interpolated_timeseries
	// `
	//	CROSS JOIN UNNEST(interpolated_avg_cpu_utilization)`

	// queryStr = `select create_time_series(time, measure_value)
	//     FROM ` + database + "." + table + `
	//     where hostname = "` + hostName + `"
	//     and measure_name = 'cpu_utilization'
	//     GROUP BY hostname,measure_name`

	allowTruncation := true
	queryInput := &timestreamquery.QueryInput{
		QueryString:           &queryStr,
		AllowResultTruncation: &allowTruncation,
		// timeout setttings?
	}
	fmt.Println("QueryInput:")
	fmt.Println(queryInput)
	// execute the query
	query, err := querySvc.Query(queryInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		// process query response
		// query response metadata
		// includes column names and types
		metadata := query.ColumnInfo
		// fmt.Println("Metadata:")
		fmt.Println(metadata)
		header := ""
		for i := 0; i < len(metadata); i++ {
			header += *metadata[i].Name
			if i != len(metadata)-1 {
				header += ", "
			}
		}

		//  query.IsDataTruncated

		// query response data
		fmt.Println("Data:")
		// process rows
		rows := query.Rows
		for i := 0; i < len(rows); i++ {
			data := rows[i].Data
			// fmt.Println(data)
			value := ""
			for j := 0; j < len(data); j++ {
				if metadata[j].Type.ScalarType != nil {
					// process simple data types
					value += *data[j].ScalarValue
				} else if metadata[j].Type.TimeSeriesMeasureValueColumnInfo != nil {
					// process complex data type 'TimeSeriesValue'
					// e.g. output of create_time_series function
					datapointList := data[j].TimeSeriesValue
					value += "["
					for k := 0; k < len(datapointList); k++ {
						time := datapointList[k].Time
						value += *time + ":" + *datapointList[k].Value.ScalarValue
						if k != len(datapointList)-1 {
							value += ", "
						}
					}
					value += "]"
				}
				// comma seperated column values
				if j != len(data)-1 {
					value += ", "
				}
			}
			fmt.Println(value)
		}
		fmt.Println("Number of rows:", len(query.Rows))

		if *query.IsDataTruncated {
			fmt.Println("data is truncated")
		}
	}
}
