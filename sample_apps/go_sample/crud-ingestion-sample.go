/**
This sample application is part of the Timestream prerelease documentation. The prerelease documentation is confidential and is provided under the terms of your nondisclosure agreement with Amazon Web Services (AWS) or other agreement governing your receipt of AWS confidential information.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"os"
	"strconv"
	"time"
)

/**
  This code sample is to run the CRUD APIs and WriteRecords API in a logical order.
*/
func main() {

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	writeSvc := timestreamwrite.New(sess)

	databaseName := flag.String("database_name", "devops", "database name string")
	tableName := flag.String("table_name", "host_metrics", "table name string")

	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Creating a database, hit enter to continue")
	//reader.ReadString('\n')

	// Create database.
	createDatabaseInput := &timestreamwrite.CreateDatabaseInput{
		DatabaseName: aws.String(*databaseName),
	}

	_, err = writeSvc.CreateDatabase(createDatabaseInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Database successfully created")
	}

	fmt.Println("Describing the database, hit enter to continue")
	reader.ReadString('\n')

	// Describe database.
	describeDatabaseInput := &timestreamwrite.DescribeDatabaseInput{
		DatabaseName: aws.String(*databaseName),
	}

	describeDatabaseOutput, err := writeSvc.DescribeDatabase(describeDatabaseInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Describe database is successful, below is the output:")
		fmt.Println(describeDatabaseOutput)
	}

	fmt.Println("Listing databases, hit enter to continue")
	reader.ReadString('\n')

	// List databases.
	listDatabasesMaxResult := int64(15)

	listDatabasesInput := &timestreamwrite.ListDatabasesInput{
		MaxResults: &listDatabasesMaxResult,
	}

	listDatabasesOutput, err := writeSvc.ListDatabases(listDatabasesInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("List databases is successful, below is the output:")
		fmt.Println(listDatabasesOutput)
	}

	fmt.Println("Creating a table, hit enter to continue")
	reader.ReadString('\n')

	// Create table.
	createTableInput := &timestreamwrite.CreateTableInput{
		DatabaseName: aws.String(*databaseName),
		TableName:    aws.String(*tableName),
	}
	_, err = writeSvc.CreateTable(createTableInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Create table is successful")
	}

	fmt.Println("Describing the table, hit enter to continue")
	reader.ReadString('\n')

	// Describe table.
	describeTableInput := &timestreamwrite.DescribeTableInput{
		DatabaseName: aws.String(*databaseName),
		TableName:    aws.String(*tableName),
	}
	describeTableOutput, err := writeSvc.DescribeTable(describeTableInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Describe table is successful, below is the output:")
		fmt.Println(describeTableOutput)
	}

	fmt.Println("Listing tables, hit enter to continue")
	reader.ReadString('\n')

	// List tables.
	listTablesMaxResult := int64(15)

	listTablesInput := &timestreamwrite.ListTablesInput{
		DatabaseName: aws.String(*databaseName),
		MaxResults:   &listTablesMaxResult,
	}
	listTablesOutput, err := writeSvc.ListTables(listTablesInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("List tables is successful, below is the output:")
		fmt.Println(listTablesOutput)
	}

	fmt.Println("Updating the table, hit enter to continue")
	reader.ReadString('\n')

	// Update table.
	magneticStoreRetentionPeriodInDays := int64(7 * 365)
	memoryStoreRetentionPeriodInHours := int64(24)

	updateTableInput := &timestreamwrite.UpdateTableInput{
		DatabaseName: aws.String(*databaseName),
		TableName:    aws.String(*tableName),
		RetentionProperties: &timestreamwrite.RetentionProperties{
			MagneticStoreRetentionPeriodInDays: &magneticStoreRetentionPeriodInDays,
			MemoryStoreRetentionPeriodInHours:  &memoryStoreRetentionPeriodInHours,
		},
	}
	updateTableOutput, err := writeSvc.UpdateTable(updateTableInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Update table is successful, below is the output:")
		fmt.Println(updateTableOutput)
	}

	fmt.Println("Ingesting records, hit enter to continue")
	reader.ReadString('\n')

	// Below code will ingest cpu_utilization and memory_utilization metric for a host on
	// region=us-east-1, az=az1, and hostname=host1

	// Get current time in seconds.
	now := time.Now()
	currentTimeInSeconds := now.Unix()
	writeRecordsInput := &timestreamwrite.WriteRecordsInput{
		DatabaseName: aws.String(*databaseName),
		TableName:    aws.String(*tableName),
		Records: []*timestreamwrite.Record{
			&timestreamwrite.Record{
				Dimensions: []*timestreamwrite.Dimension{
					&timestreamwrite.Dimension{
						Name:  aws.String("region"),
						Value: aws.String("us-east-1"),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("az"),
						Value: aws.String("az1"),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("hostname"),
						Value: aws.String("host1"),
					},
				},
				MeasureName:      aws.String("cpu_utilization"),
				MeasureValue:     aws.String("13.5"),
				MeasureValueType: aws.String("DOUBLE"),
				Timestamp:        aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
				TimestampUnit:    aws.String("SECONDS"),
			},
			&timestreamwrite.Record{
				Dimensions: []*timestreamwrite.Dimension{
					&timestreamwrite.Dimension{
						Name:  aws.String("region"),
						Value: aws.String("us-east-1"),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("az"),
						Value: aws.String("az1"),
					},
					&timestreamwrite.Dimension{
						Name:  aws.String("hostname"),
						Value: aws.String("host1"),
					},
				},
				MeasureName:      aws.String("memory_utilization"),
				MeasureValue:     aws.String("40"),
				MeasureValueType: aws.String("DOUBLE"),
				Timestamp:        aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
				TimestampUnit:    aws.String("SECONDS"),
			},
		},
	}

	_, err = writeSvc.WriteRecords(writeRecordsInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Write records is successful")
	}

	fmt.Println("Ingesting records with common attributes method, hit enter to continue")
	reader.ReadString('\n')

	// Below code will ingest the same data with common attributes approach.
	now = time.Now()
	currentTimeInSeconds = now.Unix()
	writeRecordsCommonAttributesInput := &timestreamwrite.WriteRecordsInput{
		DatabaseName: aws.String(*databaseName),
		TableName:    aws.String(*tableName),
		CommonAttributes: &timestreamwrite.Record{
			Dimensions: []*timestreamwrite.Dimension{
				&timestreamwrite.Dimension{
					Name:  aws.String("region"),
					Value: aws.String("us-east-1"),
				},
				&timestreamwrite.Dimension{
					Name:  aws.String("az"),
					Value: aws.String("az1"),
				},
				&timestreamwrite.Dimension{
					Name:  aws.String("hostname"),
					Value: aws.String("host1"),
				},
			},
			MeasureValueType: aws.String("DOUBLE"),
			Timestamp:        aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimestampUnit:    aws.String("SECONDS"),
		},
		Records: []*timestreamwrite.Record{
			&timestreamwrite.Record{
				MeasureName:  aws.String("cpu_utilization"),
				MeasureValue: aws.String("13.5"),
			},
			&timestreamwrite.Record{
				MeasureName:  aws.String("memory_utilization"),
				MeasureValue: aws.String("40"),
			},
		},
	}

	_, err = writeSvc.WriteRecords(writeRecordsCommonAttributesInput)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)
	} else {
		fmt.Println("Ingest records is successful")
	}
}
