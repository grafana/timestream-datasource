---
aliases:
  - /docs/plugins/grafana-timestream-datasource/
description: Use the Amazon Timestream data source to query and visualize time-series data from Amazon Timestream in Grafana.
keywords:
  - grafana
  - amazon timestream
  - timestream
  - aws
  - data source
  - time series
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Amazon Timestream
title: Amazon Timestream data source
weight: 100
review_date: 2026-04-28
---

# Amazon Timestream data source

The Amazon Timestream data source lets you query and visualize time-series data stored in [Amazon Timestream](https://aws.amazon.com/timestream/) directly within Grafana dashboards. Amazon Timestream is a fully managed, serverless time-series database designed for IoT and operational workloads that automatically scales to handle trillions of events per day.

## Supported features

The following table lists the features available with the Amazon Timestream data source.

| Feature | Supported |
| ----------- | --------- |
| Metrics | Yes |
| Logs | No |
| Traces | No |
| Annotations | Yes |
| Alerting | Yes |

## Get started

The following guides help you set up and use the Amazon Timestream data source:

- [Configure the Amazon Timestream data source](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/configure/)
- [Amazon Timestream query editor](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/query-editor/)
- [Template variables](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/template-variables/)
- [Annotations](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/annotations/)
- [Alerting](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/alerting/)
- [Troubleshooting](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/troubleshooting/)

## Additional features

After you configure the data source, you can:

- Use [Explore](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/explore/) to run ad-hoc queries without building a dashboard.
- Add [transformations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/panels-visualizations/query-transform-data/transform-data/) to manipulate query results.
- Set up [alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/) rules to get notified when data meets specific conditions.

## Pre-built dashboards

The Amazon Timestream data source includes a **Sample (DevOps)** dashboard. To import it:

1. Navigate to the Amazon Timestream data source configuration page.
1. Click the **Dashboards** tab.
1. Click **Import** next to **Sample (DevOps)**.

Refer to the [Sample Application section](https://docs.aws.amazon.com/timestream/latest/developerguide/Grafana.html#Grafana.sample-app) in the official Timestream documentation to set up the sample data this dashboard uses.

## Plugin updates

Always ensure that your plugin version is up to date so you have access to all current features and improvements. Navigate to **Plugins and data** > **Plugins** to check for updates. Grafana recommends upgrading to the latest Grafana version, and this applies to plugins as well.

{{< admonition type="note" >}}
Plugins are automatically updated in Grafana Cloud.
{{< /admonition >}}

## Related resources

- [Amazon Timestream documentation](https://docs.aws.amazon.com/timestream/)
- [Amazon Timestream query language reference](https://docs.aws.amazon.com/timestream/latest/developerguide/reference.html)
- [Timestream plugin GitHub repository](https://github.com/grafana/timestream-datasource/)
- [Grafana community forum](https://community.grafana.com/)
