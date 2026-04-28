---
aliases:
  - /docs/plugins/grafana-timestream-datasource/configure/
description: Configure the Amazon Timestream data source for Grafana, including authentication, provisioning, and Terraform.
keywords:
  - grafana
  - amazon timestream
  - timestream
  - aws
  - configure
  - authentication
  - provisioning
  - terraform
  - IAM
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Configure
title: Configure the Amazon Timestream data source
weight: 200
review_date: 2026-04-28
---

# Configure the Amazon Timestream data source

This document provides instructions for configuring the Amazon Timestream data source and explains available configuration options. For general information on managing data sources, refer to [Data source management](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/data-source-management/).

## Before you begin

To configure the Amazon Timestream data source, you need:

- **Grafana permissions:** The organization administrator role.
- **AWS account:** An active AWS account with Amazon Timestream enabled.
- **AWS credentials:** An IAM identity with permissions to query Timestream. Refer to [IAM policies](#iam-policies) for the minimum required permissions.

## Key concepts

If you're new to Amazon Timestream, the following terms are used throughout this documentation.

| Term | Description |
| ---- | ----------- |
| **IAM policy** | A JSON document attached to an AWS identity that grants permissions to specific API actions. |
| **Assume role** | An AWS mechanism that lets one identity temporarily take on another role's credentials, often used for cross-account access. |
| **Database** | A top-level Timestream container that organizes tables. |
| **Table** | A collection of time-series records within a database. |
| **Measure** | A specific metric or value recorded in a Timestream table, such as CPU utilization or temperature. |

## Add the data source

To add the Amazon Timestream data source:

1. Click **Connections** in the left-side menu.
1. Click **Add new connection**.
1. Type `Amazon Timestream` in the search bar.
1. Select **Amazon Timestream**.
1. Click **Add new data source**.

## Configure settings

The following table describes the available configuration settings.

| Setting | Description |
| ------- | ----------- |
| **Name** | A display name for this data source instance. |
| **Default** | Toggle on to make this the default data source for new panels. |
| **Default Region** | The AWS region where your Timestream database is located, for example `us-east-1`. |
| **Default Endpoint** | Optional. A custom endpoint URL for Timestream queries. Leave blank to use the default AWS endpoint. Use this for VPC endpoints. The default pattern is `https://query-{cell}.timestream.{region}.amazonaws.com`. |
| **Database** | Optional. The default Timestream database used by the `$__database` macro. |
| **Table** | Optional. The default Timestream table used by the `$__table` macro. |
| **Measure** | Optional. The default Timestream measure used by the `$__measure` macro. |

## Authentication

The Amazon Timestream data source uses the shared AWS authentication provided by the Grafana AWS SDK. It supports the following authentication methods:

- **AWS SDK Default** -- Uses the default credential provider chain, which checks environment variables, the shared credentials file, and the EC2/ECS instance role in order.
- **Credentials file** -- Uses a named profile from the AWS shared credentials file (`~/.aws/credentials`).
- **Access and secret key** -- Uses an IAM access key ID and secret access key that you enter directly in the data source settings.

You can restrict which authentication methods are available by configuring the `allowed_auth_providers` option in the Grafana configuration file.

For detailed information on each method, refer to [AWS authentication](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/datasources/aws-cloudwatch/aws-authentication/).

### Assume role

To access Timestream using a different IAM role, configure the data source to assume that role:

1. Enter the **Assume Role ARN** in the data source settings, for example `arn:aws:iam::123456789012:role/TimestreamReadOnly`.
1. If the role's trust policy requires it, enter an **External ID**.

## IAM policies

The IAM identity used by Grafana needs permissions to access the Timestream API. Attach a policy to the IAM user or role configured in the authentication step.

The following example grants full Timestream access:

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

For more restrictive access, grant only the specific actions needed:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "timestream:Select",
        "timestream:DescribeEndpoints",
        "timestream:ListDatabases",
        "timestream:ListTables",
        "timestream:ListMeasures",
        "timestream:DescribeDatabase",
        "timestream:DescribeTable"
      ],
      "Resource": "*"
    }
  ]
}
```

## Verify the connection

Click **Save & test**. A **Connection success** message confirms that Grafana can connect to your Timestream instance.

If the test fails:

- Verify that your authentication credentials are correct.
- Confirm that the IAM identity has the required permissions.
- Check that the selected region matches where your Timestream database is located.

For more help, refer to [Troubleshoot Amazon Timestream data source issues](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/troubleshooting/).

## Provision the data source

You can define and configure the data source using YAML files as part of Grafana's provisioning system. For more information about provisioning, refer to [Provisioning Grafana](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/provisioning/#data-sources).

The following examples demonstrate common provisioning configurations.

### AWS SDK Default

```yaml
apiVersion: 1

datasources:
  - name: Amazon Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: default
      defaultRegion: us-east-1
```

### Credentials file profile

```yaml
apiVersion: 1

datasources:
  - name: Amazon Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: credentials
      defaultRegion: us-east-1
```

### Access and secret key

```yaml
apiVersion: 1

datasources:
  - name: Amazon Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: keys
      defaultRegion: us-east-1
    secureJsonData:
      accessKey: <YOUR_ACCESS_KEY>
      secretKey: <YOUR_SECRET_KEY>
```

### AWS SDK Default with assume role

```yaml
apiVersion: 1

datasources:
  - name: Amazon Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: default
      defaultRegion: us-east-1
      assumeRoleArn: arn:aws:iam::123456789012:role/TimestreamReadOnly
      externalId: <YOUR_EXTERNAL_ID>
```

### Include default query parameters

You can also include default database, table, and measure values in any provisioning configuration. These values populate the `$__database`, `$__table`, and `$__measure` macros.

```yaml
apiVersion: 1

datasources:
  - name: Amazon Timestream
    type: grafana-timestream-datasource
    jsonData:
      authType: default
      defaultRegion: us-east-1
      defaultDatabase: my_database
      defaultTable: my_table
      defaultMeasure: cpu_utilization
```

## Provision the data source with Terraform

You can provision the Amazon Timestream data source using the [Grafana Terraform provider](https://registry.terraform.io/providers/grafana/grafana/latest/docs).

The following example creates an Amazon Timestream data source using access and secret keys:

```hcl
resource "grafana_data_source" "timestream" {
  type = "grafana-timestream-datasource"
  name = "Amazon Timestream"

  json_data_encoded = jsonencode({
    authType      = "keys"
    defaultRegion = "us-east-1"
  })

  secure_json_data_encoded = jsonencode({
    accessKey = var.aws_access_key
    secretKey = var.aws_secret_key
  })
}
```

The following example uses the AWS SDK default authentication with an assume role:

```hcl
resource "grafana_data_source" "timestream" {
  type = "grafana-timestream-datasource"
  name = "Amazon Timestream"

  json_data_encoded = jsonencode({
    authType       = "default"
    defaultRegion  = "us-east-1"
    assumeRoleArn  = "arn:aws:iam::123456789012:role/TimestreamReadOnly"
    defaultDatabase = "my_database"
    defaultTable    = "my_table"
    defaultMeasure  = "cpu_utilization"
  })
}
```

For more information, refer to the [Grafana Terraform provider documentation](https://registry.terraform.io/providers/grafana/grafana/latest/docs/resources/data_source).
