---
aliases:
  - /docs/plugins/grafana-timestream-datasource/annotations/
description: Use annotations with the Amazon Timestream data source to mark events on Grafana dashboard panels.
keywords:
  - grafana
  - amazon timestream
  - timestream
  - aws
  - annotations
  - events
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Annotations
title: Amazon Timestream annotations
weight: 500
review_date: 2026-04-28
---

# Amazon Timestream annotations

Annotations allow you to overlay event information on graphs, providing context for metric changes. The Amazon Timestream data source supports annotation queries that pull event data directly from your Timestream tables.

For general information about annotations in Grafana, refer to [Annotate visualizations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/build-dashboards/annotate-visualizations/).

## Before you begin

- [Configure the Amazon Timestream data source](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/configure/).
- Verify that your Timestream table contains event data with a timestamp column.

## Create an annotation query

To add a Timestream annotation to a dashboard:

1. Click **Dashboard settings** (gear icon).
1. Click **Annotations**.
1. Click **Add annotation query**.
1. Select your Amazon Timestream data source.
1. Enter a SQL query that returns the required columns.
1. Click **Apply**.

## Required columns

Your annotation query must return at least a time column. Grafana automatically maps the following column names to annotation properties.

| Column | Required | Description |
| ------ | -------- | ----------- |
| `time` | Yes | The timestamp for the annotation. |
| `timeEnd` | No | The end timestamp for range annotations. When present, the annotation spans from `time` to `timeEnd`. |
| `text` | No | The annotation body text displayed on hover. |
| `title` | No | A title for the annotation. |
| `tags` | No | Comma-separated tags used to filter annotations. |

## Annotation query examples

The following examples demonstrate common annotation query patterns. Annotation queries support the same macros and template variables as regular queries.

### Mark point-in-time events

The following query retrieves deployment events and displays them as point annotations:

```sql
SELECT
  time,
  measure_value::varchar AS text
FROM $__database.deployment_events
WHERE $__timeFilter
  AND measure_name = 'deployment'
ORDER BY time ASC
```

### Categorize annotations with tags

Add a `tags` column to categorize annotations and filter them in the dashboard:

```sql
SELECT
  time,
  measure_value::varchar AS text,
  environment AS tags
FROM $__database.deployment_events
WHERE $__timeFilter
  AND measure_name = 'deployment'
ORDER BY time ASC
```

### Mark time ranges

Use `timeEnd` to create range annotations that highlight a span of time, such as a maintenance window:

```sql
SELECT
  time,
  timeEnd,
  measure_value::varchar AS text
FROM $__database.maintenance_windows
WHERE $__timeFilter
  AND measure_name = 'maintenance'
ORDER BY time ASC
```
