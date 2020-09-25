# Timestream Datasource

The Timestream datasource plugin provides a support for [Amazon Timestream](https://aws.amazon.com/timestream/). Add it as a data source, then you are ready to build dashboards using timestream query results

## Add the data source

1. In the side menu under the **Configuration** link, click on **Data Sources**.
1. Click the **Add data source** button.
1. Select **Timestream** in the **Time series databases** section.

| Name                     | Description                                                                                                             |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| Name                     | The data source name. This is how you refer to the data source in panels and queries.                                   |
| Auth Provider            | Specify the provider to get credentials.                                                                                |
| Default Region           | Used in query editor to set region. (can be changed on per query basis)                                                 |
| Credentials profile name | Specify the name of the profile to use (if you use `~/.aws/credentials` file), leave blank for default.                 |
| Assume Role Arn          | Specify the ARN of the role to assume.                                                                                  |
| Endpoint (optional)      | If you need to specify an alternate service endpoint                                                                    |

## Authentication

In this section we will go through the different type of authentication you can use for X-Ray data source.

### Example AWS credentials

If the Auth Provider is `Credentials file`, then Grafana tries to get credentials in the following order:

- Hard-code credentials
- Environment variables (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`)
- Existing default config files
- ~/.aws/credentials
- IAM role for Amazon EC2

Refer to [Configuring the AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html) in the AWS documentation for more information.

### AWS credentials file

Create a file at `~/.aws/credentials`. That is the `HOME` path for user running grafana-server.

> **Note:** If the credentials file in the right place, but it is not working, then try moving your .aws file to '/usr/share/grafana/'. Make sure your credentials file has at most 0644 permissions.

Example credential file:

```bash
[default]
aws_access_key_id = <your access key>
aws_secret_access_key = <your access key>
region = us-west-2
```

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
    type: datasource
    jsonData:
      authType: credentials
      defaultRegion: us-east-1
```

### Using `accessKey` and `secretKey`

```yaml
apiVersion: 1

datasources:
  - name: Timestream
    type: datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: '<your access key>'
      secretKey: '<your secret key>'
```
