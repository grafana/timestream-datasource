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

Annotations let you mark points in time on dashboard panels to highlight events such as deployments, incidents, or configuration changes. The Amazon Timestream data source supports annotation queries that pull event data directly from your Timestream tables.

For general information about annotations, refer to [Annotate visualizations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/dashboards/build-dashboards/annotate-visualizations/).

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

Your annotation query must return at least a time column. The following columns are recognized.

| Column | Required | Description |
| ------ | -------- | ----------- |
| `time` | Yes | The timestamp for the annotation. |
| `text` | No | The annotation body text displayed on hover. |
| `title` | No | A title for the annotation. |
| `tags` | No | Comma-separated tags used to filter annotations. |

## Annotation query example

The following query retrieves deployment events from a Timestream table and displays them as annotations:

```sql
SELECT
  time,
  measure_value::varchar AS text
FROM my_database.deployment_events
WHERE $__timeFilter
  AND measure_name = 'deployment'
ORDER BY time ASC
```

### Annotation with tags

The following query includes tags to categorize annotations:

```sql
SELECT
  time,
  measure_value::varchar AS text,
  environment AS tags
FROM my_database.deployment_events
WHERE $__timeFilter
  AND measure_name = 'deployment'
ORDER BY time ASC
```
