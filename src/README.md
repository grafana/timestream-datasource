# Timestream Datasource

The Timestream datasource plugin provides a support for [Amazon Timestream](https://aws.amazon.com/timestream/). Add it as a data source, then you are ready to build dashboards using timestream query results

## Add the data source

1. In the side menu under the **Configuration** link, click on **Data Sources**.
1. Click the **Add data source** button.
1. Select **Timestream** in the **Time series databases** section.


## Authentication

The Timestream plugin authentication system matches the standard Cloudwatch plugin system.  See [the grafana documentation](https://grafana.com/docs/grafana/latest/datasources/cloudwatch/#authentication) for authentication options and setup.


Once authentication is configured, click "Save and Test" to verify the service is working. Once this is configured, you can specify default values for the configuration.

## Query editor

The query editor accepts timestream syntax in addition to the macros listed above and any dashboard template variables.

![query-editor](https://storage.googleapis.com/plugins-ci/plugins/timestream/timestream-query.png)

Type `ctrl+space` to open open the IntelliSense suggestions

## Macros

To simplify syntax and to allow for dynamic parts, like date range filters, the query can contain macros.

Macro example | Description
------------ | -------------
*$__database* | Will specify the selected database.  This may use the default from the datasource config, or the explicit value from the query editor.
*$__table* | Will specify the selected database.  This may use the default from the datasource config, or the explicit value from the query editor.
*$__measure* | Will specify the selected measure.  This may use the default from the datasource config, or the explicit value from the query editor.
*$__timeFilter* | Will be replaced by an expression that limits the time to the dashboard range
*$__interval_ms* | Will be replaced by a number that represents the amount of time a single pixel in the graph should cover


## Using Variables in Queries

Grafana template variables can be used in query Timestream query.  See the grafana documentaiton on [template variables](https://grafana.com/docs/grafana/latest/variables/) for more information on how to setup and use variables.

### Disabling quoting for multi-value variables

Grafana automatically creates a quoted, comma-separated string for multi-value variables. For example: if `server01` and `server02` are selected then it will be formatted as: `'server01', 'server02'`. To disable quoting, use the csv formatting option for variables:

`${servers:csv}`

Read more about variable formatting options in the [Variables](https://grafana.com/docs/grafana/latest/variables/advanced-variable-format-options/) documentation.


### Alerting

See the [Alerting](https://grafana.com/docs/grafana/latest/alerting/alerts-overview/) documentation for more on Grafana alerts.

## Configure the data source with provisioning

You can configure data sources using config files with Grafana's provisioning system. You can read more about how it works and all the settings you can set for data sources on the [provisioning docs page](https://grafana.com/docs/grafana/latest/administration/provisioning/).

Here are some provisioning examples for this data source.

### Using a credentials file

If you are using Credentials file authentication type, then you should use a credentials file with a config like this.

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
