# AWS Timestream Datasource

The Timestream datasource plugin provides a support for [Amazon Timestream](https://aws.amazon.com/timestream/). Add it as a data source, then you are ready to build dashboards using timestream query results

## Getting started

1. [Install the plugin](https://grafana.com/docs/grafana/latest/administration/plugin-management/#install-grafana-plugins)
1. [Add a new data source with the UI](https://grafana.com/docs/grafana/latest/datasources/#add-a-data-source) or [provision one](https://grafana.com/docs/grafana/latest/administration/provisioning/)
1. [Configure the data source](#configuring-the-data-source)
1. [Start making queries](#querying-the-data-source)

## Configuring the data source

### Authentication
Depending on the environment in which it is run, Grafana supports different authentication providers such as keys, a credentials file, or using the "Default" provider from AWS which supports using service-based IAM roles. These providers can be manually enabled/disabled with the `allowed_auth_providers` field in Grafana's config file. To read more about supported authentication providers refer to [the AWS authentication section](https://grafana.com/docs/grafana/latest/datasources/aws-cloudwatch/aws-authentication/#select-an-authentication-method)

### IAM policies

Grafana needs permissions granted via IAM to be able to read data from the Timestream API. You can attach these permissions to the IAM role or IAM user configured in the previous step.

Here is a policy example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["timestream:*"],
      "Resource": "*"
    }
  ]
}
```
### Configure the data source with provisioning

You can configure data sources using config files with Grafana's provisioning system. You can read more about how it works and all the settings you can set for data sources on the [provisioning docs page](https://grafana.com/docs/grafana/latest/administration/provisioning/).

Here are some provisioning examples for this data source.

### Using AWS SDK (default)

```yaml
apiVersion: 1
datasources:
  - name: Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: default
      defaultRegion: eu-west-2
```

### Using credentials' profile name (non-default)

```yaml
apiVersion: 1

datasources:
  - name: Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: credentials
      defaultRegion: us-east-1
```

### Using `accessKey` and `secretKey`

```yaml
apiVersion: 1

datasources:
  - name: Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: '<your access key>'
      secretKey: '<your secret key>'
```

### Using AWS SDK Default and ARN of IAM Role to Assume

```yaml
apiVersion: 1
datasources:
  - name: Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: default
      assumeRoleArn: arn:aws:iam::123456789012:root
      defaultRegion: eu-west-2
```

## Querying the data source

The query editor accepts timestream syntax in addition to the macros listed above and any dashboard template variables.

Type `ctrl+space` to open open the IntelliSense suggestions

## Macros

To simplify syntax and to allow for dynamic parts, like date range filters, the query can contain macros.

| Macro example          | Description                                                                                                                           |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| _$\_\_database_        | Will specify the selected database. This may use the default from the datasource config, or the explicit value from the query editor. |
| _$\_\_table_           | Will specify the selected database. This may use the default from the datasource config, or the explicit value from the query editor. |
| _$\_\_measure_         | Will specify the selected measure. This may use the default from the datasource config, or the explicit value from the query editor.  |
| _$\_\_timeFilter_      | Will be replaced by an expression that limits the time to the dashboard range.                                                        |
| _$\_\_timeFrom_        | Will be replaced by the number in milliseconds at the start of the dashboard range.                                                   |
| _$\_\_timeTo_          | Will be replaced by the number in milliseconds at the end of the dashboard range.                                                     |
| _$\_\_interval_ms_     | Will be replaced by a number in time format that represents the amount of time a single pixel in the graph should cover.              |
| _$\_\_interval_raw_ms_ | Will be replaced by the number in milliseconds that represents the amount of time a single pixel in the graph should cover.           |

## Using Variables in Queries

Instead of hard-coding server, application and sensor names in your Timestream queries, you can use variables. The variables are listed as dropdown select boxes at the top of the dashboard. These dropdowns make it easy to change the display of data in your dashboard.

For an introduction to templating and template variables, refer to the [Templating](https://grafana.com/docs/grafana/latest/variables/) documentation.

### Disabling quoting for multi-value variables

Grafana automatically creates a quoted, comma-separated string for multi-value variables. For example: if `server01` and `server02` are selected then it will be formatted as: `'server01', 'server02'`. To disable quoting, use the csv formatting option for variables:

`${servers:csv}`

Read more about variable formatting options in the [Variables](https://grafana.com/docs/grafana/latest/variables/advanced-variable-format-options/) documentation.

### Alerting

[Alerting](https://grafana.com/docs/grafana/latest/alerting/) queries should contain a time series field. Queries without this field will return an error: "input data must be a wide series but got type long". To return time series, you can use the [`CREATE_TIME_SERIES` function](https://docs.aws.amazon.com/timestream/latest/developerguide/timeseries-specific-constructs.views.html). For example:

```sql
SELECT
    silo, microservice_name, instance_name,
    CREATE_TIME_SERIES(time, measure_value::double) AS gc_pause
FROM $__database.$__table
WHERE $__timeFilter
    AND measure_name = '$__measure'
    AND region = 'ap-northeast-1'
    AND cell = 'ap-northeast-1-cell-5'
    AND silo = 'ap-northeast-1-cell-5-silo-2'
    AND availability_zone = 'ap-northeast-1-3'
    AND microservice_name = 'zeus'
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

> **Note**: Results for Timestream queries are returned in different pages (if necessary) by default. To ensure that all pages are processed before evaluating an alert, mark the "Wait for all queries" checkbox underneath the "Render" query editor section for all alert queries.

### Sample Dashboard

This plugin contains one sample dashboard. Please consult the [Sample Application section](https://docs.aws.amazon.com/timestream/latest/developerguide/Grafana.html#Grafana.sample-app) in the official Timestream doc to set it up.


### Plugin repository

You can request new features, report issues, or contribute code directly through the [Timestream plugin Github repository](https://github.com/grafana/timestream-datasource/)