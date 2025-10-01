# AWS Timestream Datasource

The Timestream datasource plugin provides a support for [Amazon Timestream](https://aws.amazon.com/timestream/). Add it as a data source, then you are ready to build dashboards using timestream query results

## Getting started

1. [Install the plugin](https://grafana.com/docs/grafana/latest/administration/plugin-management/#install-grafana-plugins)
1. [Add a new data source with the UI](https://grafana.com/docs/grafana/latest/datasources/#add-a-data-source) or [provision one](https://grafana.com/docs/grafana/latest/administration/provisioning/)
1. [Configure the data source](#configuration-options)
1. Start making queries

### Configuration options

Depending on the environment in which it is run, Grafana supports different authentication providers such as keys, a credentials file, or using the "Default" provider from AWS which supports using service-based IAM roles. These providers can be manually enabled/disabled with the allowed_auth_providers field. To read more about supported authentication providers refer to [the Cloud Watch Data Source's documentation](https://grafana.com/docs/grafana/latest/datasources/aws-cloudwatch/aws-authentication/#select-an-authentication-method)

### Plugin repository

You can request new features, report issues, or contribute code directly through the [Timestream plugin Github repository](https://github.com/grafana/grafana-timestream-datasource/)