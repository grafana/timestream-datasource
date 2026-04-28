---
aliases:
  - /docs/plugins/grafana-timestream-datasource/query-editor/
description: Use the Amazon Timestream query editor to write SQL queries, use macros, and format query results.
keywords:
  - grafana
  - amazon timestream
  - timestream
  - aws
  - query editor
  - sql
  - macros
  - time series
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Query editor
title: Amazon Timestream query editor
weight: 300
review_date: 2026-04-28
---

# Amazon Timestream query editor

The Amazon Timestream query editor lets you write SQL queries against your Timestream databases. It includes a code editor with IntelliSense, macros for dynamic values, and configurable output formats.

## Before you begin

- [Configure the Amazon Timestream data source](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/configure/).
- Verify that your IAM identity has permissions to query the target database and table.

## Key concepts

The following terms are used throughout the query editor documentation.

| Term | Description |
| ---- | ----------- |
| **Database** | A top-level Timestream container that organizes tables. |
| **Table** | A collection of time-series records within a database. |
| **Measure** | A named metric or value in a Timestream table, such as `cpu_utilization` or `temperature`. |
| **Macro** | A placeholder in a query that Grafana replaces with a dynamic value at execution time. |

## Query editor fields

The query editor provides the following fields and controls.

| Field | Description |
| ----- | ----------- |
| **Database** | The Timestream database to query. Populates the `$__database` macro. Falls back to the default database set in the data source configuration. |
| **Table** | The table within the selected database. Populates the `$__table` macro. The table list updates when you change the database. |
| **Measure** | The measure within the selected table. Populates the `$__measure` macro. The measure list updates when you change the database or table. |
| **Wait for all queries** | When enabled, the plugin fetches all paginated result pages before returning data. Enable this for [alerting queries](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/alerting/). |
| **Format as** | Controls the output format: **Table** (default) or **Time Series**. Time-series queries must return times in ascending order using `ORDER BY time ASC`. |
| **Sample queries** | A drop-down of pre-built queries to help you get started. Selecting a sample replaces the current query. |

## Write a query

Use the SQL editor to write [Timestream SQL](https://docs.aws.amazon.com/timestream/latest/developerguide/reference.html) queries. The editor supports IntelliSense for column names, table names, and macros. Press `Ctrl+Space` to trigger auto-complete suggestions.

To run a query, press `Ctrl+Enter` or click **Run query**.

### Basic query example

The following query selects all records from the configured database and table within the dashboard time range:

```sql
SELECT *
FROM $__database.$__table
WHERE $__timeFilter
ORDER BY time ASC
LIMIT 100
```

## Macros

Use macros in your queries to insert dynamic values like time ranges, intervals, and configured defaults. The query engine replaces macros with their computed values before sending the query to Timestream.

| Macro | Description |
| ----- | ----------- |
| `$__database` | The database selected in the query editor, or the default database from the data source configuration. |
| `$__table` | The table selected in the query editor, or the default table from the data source configuration. |
| `$__measure` | The measure selected in the query editor, or the default measure from the data source configuration. |
| `$__timeFilter` | An expression that limits results to the dashboard time range, for example `time BETWEEN from_milliseconds(1234) AND from_milliseconds(5678)`. |
| `$__timeFrom` | The start of the dashboard time range in milliseconds. |
| `$__timeTo` | The end of the dashboard time range in milliseconds. |
| `$__interval` | A Timestream duration literal representing the calculated interval for the panel width, for example `60000ms`. |
| `$__interval_ms` | Same as `$__interval`. A Timestream duration literal representing the calculated interval in milliseconds. |
| `$__interval_raw_ms` | The calculated interval as a plain integer in milliseconds, for example `60000`. |
| `$__now_ms` | The current time in milliseconds. |

### Macro example

The following query combines several macros to aggregate data into dashboard-appropriate intervals:

```sql
SELECT
  bin(time, $__interval_ms) AS binned_time,
  measure_name,
  avg(measure_value::double) AS avg_value
FROM $__database.$__table
WHERE $__timeFilter
  AND measure_name = '$__measure'
GROUP BY bin(time, $__interval_ms), measure_name
ORDER BY binned_time ASC
```

This query uses `$__database`, `$__table`, and `$__measure` to reference the selections in the query editor, `$__timeFilter` to scope results to the dashboard time range, and `$__interval_ms` to group data into intervals that match the panel width.

## Sample queries

The query editor includes built-in sample queries you can select from the **Sample queries** drop-down. These queries use macros and adapt to your selected database and table.

| Sample query | Description |
| ------------ | ----------- |
| **Show databases** | Lists all databases in the Timestream instance. |
| **Show tables** | Lists tables in the selected database. |
| **Describe table** | Describes the schema of the selected table. |
| **Show measurements** | Lists all measures in the selected table. |
| **First 10 rows** | Returns the first 10 rows from the selected table. |

## Use cases

The following examples demonstrate common query patterns for Timestream data.

### Aggregate metrics over time

Use `bin()` with `$__interval_ms` to group data into intervals that match the panel resolution. This is the most common pattern for time-series visualizations.

```sql
SELECT
  bin(time, $__interval_ms) AS binned_time,
  instance_name,
  avg(measure_value::double) AS avg_cpu
FROM $__database.$__table
WHERE $__timeFilter
  AND measure_name = 'cpu_utilization'
GROUP BY bin(time, $__interval_ms), instance_name
ORDER BY binned_time ASC
```

### Compare metrics across dimensions

Group by a dimension column to compare values across hosts, regions, or other attributes in a single panel.

```sql
SELECT
  bin(time, $__interval_ms) AS binned_time,
  hostname,
  max(measure_value::double) AS peak_memory,
  avg(measure_value::double) AS avg_memory
FROM $__database.$__table
WHERE $__timeFilter
  AND measure_name = 'memory_utilization'
GROUP BY bin(time, $__interval_ms), hostname
ORDER BY binned_time ASC
```

### Explore table schema

Use Timestream's metadata queries to discover available data before building visualizations.

```sql
SHOW MEASURES FROM $__database.$__table
```

```sql
DESCRIBE $__database.$__table
```

### Query the latest values

Use `$__now_ms` and `$__timeFrom` to build custom time-range logic, for example, to find the most recent reading per sensor.

```sql
SELECT
  sensor_id,
  measure_name,
  time,
  measure_value::double AS value
FROM $__database.$__table
WHERE time BETWEEN from_milliseconds($__timeFrom) AND from_milliseconds($__now_ms)
  AND measure_name = 'temperature'
ORDER BY time DESC
LIMIT 10
```

### IoT multi-measure queries

Query multiple measure values from multi-measure records in a single query.

```sql
SELECT
  bin(time, $__interval_ms) AS binned_time,
  device_id,
  avg(temperature) AS avg_temp,
  avg(humidity) AS avg_humidity
FROM $__database.$__table
WHERE $__timeFilter
GROUP BY bin(time, $__interval_ms), device_id
ORDER BY binned_time ASC
```
