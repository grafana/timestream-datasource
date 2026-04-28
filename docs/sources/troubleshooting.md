---
aliases:
  - /docs/plugins/grafana-timestream-datasource/troubleshooting/
description: Troubleshoot common issues with the Amazon Timestream data source in Grafana, including authentication, connection, and query errors.
keywords:
  - grafana
  - amazon timestream
  - timestream
  - aws
  - troubleshooting
  - errors
  - authentication
  - query
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Troubleshooting
title: Troubleshoot Amazon Timestream data source issues
weight: 700
review_date: 2026-04-28
---

# Troubleshoot Amazon Timestream data source issues

This document provides solutions to common issues you may encounter when configuring or using the Amazon Timestream data source. For configuration instructions, refer to [Configure the Amazon Timestream data source](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/configure/).

## Authentication errors

These errors occur when credentials are invalid, missing, or don't have the required permissions.

### "Access denied" or "Authorization failed"

**Symptoms:**

- **Save & test** fails with authorization errors.
- Queries return access denied messages.
- Database, table, or measure drop-downs don't load.

**Possible causes and solutions:**

| Cause | Solution |
| ----- | -------- |
| Missing IAM permissions | Attach a policy granting `timestream:*` or the specific actions listed in the [IAM policies](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/configure/#iam-policies) section. |
| Invalid credentials | Verify credentials in the AWS IAM console. Regenerate the access key if necessary. |
| Expired credentials | Create a new access key and update the data source configuration. |
| Wrong region | Verify that the **Default Region** matches the region where your Timestream database is located. |
| Assume role misconfigured | Verify the role ARN, check the trust policy, and confirm the external ID matches if one is required. |

### "ExpiredTokenException" or "Token has expired"

**Symptoms:**

- Queries fail intermittently with token expiration errors.
- The data source works initially but stops after a period of time.

**Solutions:**

1. If you use temporary credentials, verify that the session token is current.
1. For assume role configurations, verify that the source identity has permission to call `sts:AssumeRole`.
1. Check that the Grafana server clock is synchronized. Expired token errors can occur when the system clock drifts.

## Connection errors

These errors occur when Grafana can't reach the Timestream API endpoints.

### "Connection refused" or timeout errors

**Symptoms:**

- **Save & test** times out.
- Queries fail with network errors.
- Intermittent connection failures.

**Solutions:**

1. Verify network connectivity from the Grafana server to Timestream endpoints in the configured region.
1. Check that firewall rules allow outbound HTTPS traffic on port 443.
1. If you use a VPC endpoint, verify the endpoint is correctly configured and set the custom endpoint URL in the data source settings.
1. For Grafana Cloud, configure [Private data source connect](https://grafana.com/docs/grafana-cloud/connect-externally-hosted/private-data-source-connect/) if accessing private resources.

### Custom endpoint issues

**Symptoms:**

- Queries fail after configuring a custom endpoint.
- 404 errors related to `DescribeEndpoints`.
- Errors related to endpoint discovery.

**Solutions:**

1. Verify the custom endpoint URL is correct and reachable from the Grafana server.
1. When a custom endpoint is configured (for example, a VPC endpoint), the plugin automatically disables AWS endpoint discovery. This prevents 404 errors from VPC endpoints that don't implement the `DescribeEndpoints` API. Verify that the endpoint URL points directly to the Timestream query endpoint.
1. If the VPC endpoint blocks egress to public AWS endpoints, ensure endpoint discovery is not being forced by another configuration.
1. To revert to the default endpoint, clear the **Default Endpoint** field in the data source configuration.

## Query errors

These errors occur when executing queries against Timestream.

### "No data" or empty results

**Symptoms:**

- Queries run without error but return no data.
- Panels display a "No data" message.

**Possible causes and solutions:**

| Cause | Solution |
| ----- | -------- |
| Time range doesn't contain data | Expand the dashboard time range or verify data exists in the AWS Timestream console. |
| Wrong database, table, or measure selected | Verify you've selected the correct database, table, and measure in the query editor. |
| Missing read permissions | Verify the IAM identity has `timestream:Select` permission on the target resources. |
| Filter excludes all results | Review `WHERE` clause conditions. Try removing filters to confirm data exists, then add them back incrementally. |

### "input data must be a wide series but got type long"

**Symptoms:**

- Alert rules fail with this error.
- The query works in dashboard panels but fails in alerting.

**Solutions:**

1. Use the [`CREATE_TIME_SERIES`](https://docs.aws.amazon.com/timestream/latest/developerguide/timeseries-specific-constructs.views.html) function to return data in wide time-series format. Refer to [Alerting](https://grafana.com/docs/plugins/grafana-timestream-datasource/latest/alerting/) for examples.
1. Enable **Wait for all queries** in the query editor to ensure all result pages are processed before the alert evaluates.

### Query timeout

**Symptoms:**

- Queries run for a long time and then fail.
- Errors mention timeout or query limits.

**Solutions:**

1. Narrow the dashboard time range to reduce the volume of data scanned.
1. Add filters to the `WHERE` clause to reduce the result set.
1. Use `bin(time, <interval>)` with a larger interval to reduce the number of returned rows.
1. Break complex queries into smaller parts using multiple panels.

### Pagination and incomplete results

**Symptoms:**

- Only a subset of expected results appears.
- Results differ between dashboard panels and alerting.

**Solutions:**

1. Enable **Wait for all queries** in the query editor to fetch all result pages before returning data.
1. If you don't enable this option, the plugin streams pages incrementally. This works for dashboard panels but can cause incomplete results in alerting evaluations.

## Template variable errors

These errors occur when using template variables with the data source.

### Variables return no values

**Symptoms:**

- Variable drop-downs are empty.
- The variable preview shows no results when you click **Run query**.

**Solutions:**

1. Verify the data source connection is working by running **Save & test** in the data source settings.
1. Check that the variable query is valid SQL that returns at least one column.
1. For cascading variables (where one depends on another), verify that the parent variable has a valid selection.
1. Verify the IAM identity has permission to list databases, tables, or run the specific query used by the variable.

### Variables are slow to load

**Solutions:**

1. Set the variable refresh to **On dashboard load** instead of **On time range change** to avoid reloading on every time range adjustment.
1. Simplify variable queries. For example, use `SHOW DATABASES` instead of a full `SELECT DISTINCT` query when listing databases.
1. Reduce the scope of variable queries by adding filters.

## Performance issues

These issues relate to slow queries or AWS API limits.

### API throttling or rate limit errors

**Symptoms:**

- "Rate exceeded" or throttling errors.
- Dashboard panels intermittently fail to load.
- Multiple panels fail simultaneously.

**Solutions:**

1. Reduce the frequency of dashboard auto-refresh.
1. Use larger time intervals in `bin()` to reduce the number of data points per query.
1. Enable [query caching](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/administration/data-source-management/#query-caching) in Grafana (available in Grafana Enterprise and Grafana Cloud).
1. Request a quota increase from AWS through the [Service Quotas console](https://console.aws.amazon.com/servicequotas/).

## Enable debug logging

To capture detailed error information for troubleshooting:

1. Set the Grafana log level to `debug` in the configuration file:

   ```ini
   [log]
   level = debug
   ```

1. Review logs in `/var/log/grafana/grafana.log` or your configured log location.
1. Look for entries containing `timestream` for request and response details.
1. Reset the log level to `info` after troubleshooting to avoid excessive log volume.

## Minimum Grafana version

The Amazon Timestream data source requires Grafana 10.4 or later. If you encounter unexpected errors, verify that your Grafana version meets this requirement.

## Get additional help

If you've tried these solutions and still encounter issues:

1. Check the [Grafana community forums](https://community.grafana.com/) for similar issues.
1. Review the [Timestream plugin GitHub issues](https://github.com/grafana/timestream-datasource/issues) for known bugs.
1. Refer to the [Amazon Timestream documentation](https://docs.aws.amazon.com/timestream/) for service-specific guidance.
1. Contact Grafana Support if you're a Cloud Pro, Cloud Advanced, or Enterprise user.
1. When reporting issues, include:
   - Grafana version and plugin version
   - Error messages (redact sensitive information)
   - Steps to reproduce
   - Relevant configuration (redact credentials)
