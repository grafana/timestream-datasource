/**
This sample application is part of the Timestream prerelease documentation. The prerelease documentation is confidential and is provided under the terms of your nondisclosure agreement with Amazon Web Services (AWS) or other agreement governing your receipt of AWS confidential information.
*/

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"io"
	"os"
	"strconv"
	"time"
)

/**
  This code sample is to read data from a CSV file and ingest data into a Timestream table. Each line of the CSV file is a record to ingest.
  The record schema is fixed, the format is [dimension_name_1, dimension_value_1, dimension_name_2, dimension_value_2, dimension_name_2, dimension_value_2, measure_name, measure_value, measure_data_type, timestamp, timestamp_unit].
  The code will replace the timestamp in the record with a timestamp in the range [current_epoch_in_seconds - number_of_records * 10, current_epoch_in_seconds].
*/
func main() {

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	writeSvc := timestreamwrite.New(sess)
	querySvc := timestreamquery.New(sess)

	databaseName := flag.String("database_name", "devops", "database name string")
	tableName := flag.String("table_name", "host_metrics", "table name string")
	testFileName := flag.String("test_file", "../data/test.csv", "CSV file containing the data to ingest")

	flag.Parse()

	// Describe database.
	describeDatabaseInput := &timestreamwrite.DescribeDatabaseInput{
		DatabaseName: aws.String(*databaseName),
	}

	describeDatabaseOutput, err := writeSvc.DescribeDatabase(describeDatabaseInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
		// Create database if database doesn't exist.
		serr, ok := err.(*timestreamwrite.ResourceNotFoundException)
		fmt.Println(serr)
		if ok {
			fmt.Println("Creating database")
			createDatabaseInput := &timestreamwrite.CreateDatabaseInput{
				DatabaseName: aws.String(*databaseName),
			}

			_, err = writeSvc.CreateDatabase(createDatabaseInput)

			if err != nil {
				fmt.Println("Error:")
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("Database exists")
		fmt.Println(describeDatabaseOutput)
	}

	// Describe table.
	describeTableInput := &timestreamwrite.DescribeTableInput{
		DatabaseName: aws.String(*databaseName),
		TableName:    aws.String(*tableName),
	}
	describeTableOutput, err := writeSvc.DescribeTable(describeTableInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
		serr, ok := err.(*timestreamwrite.ResourceNotFoundException)
		fmt.Println(serr)
		if ok {
			// Create table if table doesn't exist.
			fmt.Println("Creating the table")
			createTableInput := &timestreamwrite.CreateTableInput{
				DatabaseName: aws.String(*databaseName),
				TableName:    aws.String(*tableName),
			}
			_, err = writeSvc.CreateTable(createTableInput)

			if err != nil {
				fmt.Println("Error:")
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("Table exists")
		fmt.Println(describeTableOutput)
	}

	csvfile, err := os.Open(*testFileName)
	records := make([]*timestreamwrite.Record, 0)
	if err != nil {
		fmt.Println("Couldn't open the csv file", err)
	}

	// Get current time in nano seconds.
	currentTimeInMilliSeconds := time.Now().UnixNano() / int64(time.Millisecond)
	// Counter for number of records.
	counter := int64(0)
	reader := csv.NewReader(csvfile)
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		records = append(records, &timestreamwrite.Record{
			Dimensions: []*timestreamwrite.Dimension{
				&timestreamwrite.Dimension{
					Name:  aws.String(record[0]),
					Value: aws.String(record[1]),
				},
				&timestreamwrite.Dimension{
					Name:  aws.String(record[2]),
					Value: aws.String(record[3]),
				},
				&timestreamwrite.Dimension{
					Name:  aws.String(record[4]),
					Value: aws.String(record[5]),
				},
			},
			MeasureName:      aws.String(record[6]),
			MeasureValue:     aws.String(record[7]),
			MeasureValueType: aws.String(record[8]),
			Timestamp:        aws.String(strconv.FormatInt(currentTimeInMilliSeconds-counter*int64(50), 10)),
			TimestampUnit:    aws.String("MILLISECONDS"),
		})

		counter++
		// WriteRecordsRequest has 100 records limit per request.
		if counter%100 == 0 {
			writeRecordsInput := &timestreamwrite.WriteRecordsInput{
				DatabaseName: aws.String(*databaseName),
				TableName:    aws.String(*tableName),
				Records:      records,
			}
			_, err = writeSvc.WriteRecords(writeRecordsInput)

			if err != nil {
				fmt.Println("Error:")
				fmt.Println(err)
			} else {
				fmt.Print("Ingested ", counter)
				fmt.Println(" records to the table.")
			}
			records = make([]*timestreamwrite.Record, 0)
		}
	}

	queryInput := &timestreamquery.QueryInput{
		QueryString: aws.String("select count(*) from " + *databaseName + "." + *tableName),
	}
	// execute the query
	query, err := querySvc.Query(queryInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println(query)
	}
}
