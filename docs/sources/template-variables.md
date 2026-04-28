---
aliases:
  - /docs/plugins/grafana-timestream-datasource/template-variables/
description: Use template variables with the Amazon Timestream data source to create dynamic, reusable dashboards.
keywords:
  - grafana
  - amazon timestream
  - timestream
  - aws
  - template variables
  - variables
  - dashboard
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Template variables
title: Amazon Timestream template variables
weight: 400
review_date: 2026-04-28
---

# Amazon Timestream template variables

Use template variables to create dynamic, reusable dashboards. Instead of hard-coding database names, table names, or filter values in your queries, you can use variables that appear as drop-down selectors at the top of the dashboard.

For an introduction to template variables, refer to [Variables](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/).

## Before you begin

- [Configure the Amazon Timestream data source](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/configure/).
- Understand [Grafana template variables](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/).

## Supported variable types

| Variable type | Supported |
| ------------- | --------- |
| Query | Yes |
| Custom | Yes |
| Data source | Yes |

## Create a query variable

Query variables let you dynamically populate a drop-down with values returned from a Timestream SQL query.

To create a query variable:

1. Navigate to **Dashboard settings** > **Variables**.
1. Click **Add variable**.
1. Select **Query** as the variable type.
1. Select your Amazon Timestream data source.
1. Enter a Timestream SQL query in the **Query** field. The first column of the result set populates the variable options.
1. Click **Run query** to preview the values.
1. Click **Apply**.

### Variable query examples

The following queries demonstrate common patterns for populating variables.

**List all databases:**

```sql
SHOW DATABASES
```

**List tables in a database (using a `database` variable):**

```sql
SHOW TABLES FROM ${database}
```

**List distinct values for a dimension:**

```sql
SELECT DISTINCT region FROM ${database}.${table}
```

**Cascading variables (filter by a parent variable):**

Create a `region` variable first, then create an `instance` variable that depends on it:

```sql
SELECT DISTINCT instance_name
FROM ${database}.${table}
WHERE region = '${region}'
```

When the `region` selection changes, the `instance` variable automatically refreshes to show only instances in the selected region.

## Use variables in queries

Reference variables in your Timestream queries with the `$variable_name` or `${variable_name}` syntax:

```sql
SELECT
  bin(time, $__interval_ms) AS binned_time,
  avg(measure_value::double) AS avg_value
FROM $__database.$__table
WHERE $__timeFilter
  AND region = '$region'
  AND instance_name = '$instance'
ORDER BY binned_time ASC
```

Variables also work in the **Database**, **Table**, and **Measure** selector fields in the query editor. For example, you can set a variable as the database and all queries using `$__database` will reflect the selected value.

## Multi-value variables

When a multi-value variable has multiple selections, the plugin automatically wraps each value in single quotes and joins them with commas. For example, if `server01` and `server02` are selected, `$servers` renders as `'server01','server02'`. This format works directly in SQL `IN` clauses:

```sql
SELECT
  bin(time, $__interval_ms) AS binned_time,
  hostname,
  avg(measure_value::double) AS avg_cpu
FROM $__database.$__table
WHERE $__timeFilter
  AND measure_name = 'cpu_utilization'
  AND hostname IN ($servers)
GROUP BY bin(time, $__interval_ms), hostname
ORDER BY binned_time ASC
```

### Disable quoting for multi-value variables

To disable the automatic quoting, use the `csv` formatting option:

```sql
${servers:csv}
```

This renders the values as `server01,server02` without quotes, which is useful in contexts where SQL string quoting isn't needed.

For more information about variable formatting options, refer to [Advanced variable format options](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/variables/variable-syntax/).
