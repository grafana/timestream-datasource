---
aliases:
  - /docs/plugins/grafana-timestream-datasource/alerting/
description: Set up Grafana alerting with the Amazon Timestream data source, including CREATE_TIME_SERIES and pagination requirements.
keywords:
  - grafana
  - amazon timestream
  - timestream
  - aws
  - alerting
  - alerts
  - CREATE_TIME_SERIES
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Alerting
title: Amazon Timestream alerting
weight: 600
review_date: 2026-04-28
---

# Amazon Timestream alerting

You can use the Amazon Timestream data source with Grafana Alerting to create alert rules that trigger when your time-series data meets specific conditions.

For general information about Grafana Alerting, refer to [Alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/).

## Before you begin

- [Configure the Amazon Timestream data source](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/configure/).
- Understand the [Amazon Timestream query editor](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/query-editor/).
- Review [Grafana Alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/).

## Alert query requirements

Timestream alert queries must return data in **wide time-series format**. If your query returns data in the default long format, the alert evaluation fails with the error `input data must be a wide series but got type long`.

To return time-series data, use the Timestream [`CREATE_TIME_SERIES`](https://docs.aws.amazon.com/timestream/latest/developerguide/timeseries-specific-constructs.views.html) function, which converts individual rows into a time-series structure.

## Enable Wait for all queries

Timestream returns large result sets across multiple pages. By default, the plugin streams pages incrementally. For alerting, you must enable **Wait for all queries** to ensure all pages are processed before the alert condition is evaluated.

To enable this setting:

1. Open the alert rule query editor.
1. Toggle **Wait for all queries** to on.

{{< admonition type="caution" >}}
If **Wait for all queries** is not enabled, alert evaluations may use incomplete data from only the first page of results, leading to missed or false alerts.
{{< /admonition >}}

## Alert query examples

Alert queries support the same macros and template variables as regular queries. Test your alert query in a dashboard panel first to verify it returns data in the expected format before creating an alert rule.

### Simple threshold alert

The following query monitors CPU utilization per instance. Each `instance_name` becomes a separate time series that you can set a threshold condition on:

```sql
SELECT
  instance_name,
  CREATE_TIME_SERIES(time, measure_value::double) AS cpu_utilization
FROM $__database.$__table
WHERE $__timeFilter
  AND measure_name = 'cpu_utilization'
GROUP BY instance_name
```

### Alert on string values

Grafana alerting requires numeric values. If your Timestream table stores status as strings (for example, `"HEALTHY"` or `"UNHEALTHY"`), use a `CASE` expression to convert them to numeric values that you can set a threshold on:

```sql
SELECT
  instance_name,
  CREATE_TIME_SERIES(time,
    CASE measure_value::varchar
      WHEN 'UNHEALTHY' THEN 1.0
      WHEN 'DEGRADED' THEN 0.5
      ELSE 0.0
    END
  ) AS health_status
FROM $__database.$__table
WHERE $__timeFilter
  AND measure_name = 'status'
GROUP BY instance_name
```

In this example, set the alert condition to trigger when `health_status` is greater than `0` to detect any non-healthy state.

### Multi-dimensional alert

The following query creates time series grouped by multiple dimensions. Use `LIMIT` and `ORDER BY` to alert only on the top offenders:

```sql
SELECT
  silo,
  microservice_name,
  instance_name,
  CREATE_TIME_SERIES(time, measure_value::double) AS gc_pause
FROM $__database.$__table
WHERE $__timeFilter
  AND measure_name = '$__measure'
  AND region = 'us-east-1'
  AND cell = 'us-east-1-cell-1'
  AND silo = 'us-east-1-cell-1-silo-1'
  AND availability_zone = 'us-east-1-1'
  AND microservice_name = 'apollo'
GROUP BY region,
  cell,
  silo,
  availability_zone,
  microservice_name,
  instance_name,
  process_name,
  jdk_version
ORDER BY AVG(measure_value::double) DESC
LIMIT 3
```

## Duplicate alerts in high availability setups

In Grafana high availability (HA) deployments, each Grafana instance evaluates alert rules independently. This can result in duplicate notifications for the same alert condition.

To prevent duplicate alerts:

1. Configure a single Grafana instance as the alert evaluator, or use [Grafana Alerting HA](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/set-up/configure-high-availability/) to coordinate alert evaluation across instances.
1. Ensure that all instances share the same database for alert state, so only one instance fires each alert.
