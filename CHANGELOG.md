# Change Log

All notable changes to this project will be documented in this file.

## 2.7.1

- Update @grafana/aws-sdk to fix a bug in temporary credentials

## 2.7.0

- Update grafana-aws-sdk to v0.19.1 to add `il-central-1` to opt-in region list

## 2.6.2

- Update grafana-aws-sdk-react version to use grafana/runtime instead of grafanaBootData [#237](https://github.com/grafana/grafana-aws-sdk/pull/237)
- Remove code coverage workflow [#234](https://github.com/grafana/grafana-aws-sdk/pull/234)

## 2.6.1

- Update grafana-aws-sdk version to include new region in opt-in region list https://github.com/grafana/grafana-aws-sdk/pull/80
- Security: Upgrade Go in build process to 1.20.4
- Update grafana-plugin-sdk-go version to 0.161.0 to avoid a potential http header problem. https://github.com/grafana/athena-datasource/issues/233

## 2.6.0

- Update backend dependencies

## 2.5.0

- Update @grafana/aws-sdk by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/216
- Increase label width to fix overflow. by @chinu-anand in https://github.com/grafana/timestream-datasource/pull/217
- migrate to create-plugin by @iwysiu in https://github.com/grafana/timestream-datasource/pull/199
- Upgrade grafana-aws-sdk by @fridgepoet in https://github.com/grafana/timestream-datasource/pull/223

**Full Changelog**: https://github.com/grafana/timestream-datasource/compare/v2.4.0...v2.5.0

## 2.4.0

- Fix: SQLEditor: Use queryRef to call onChange [#209](https://github.com/grafana/timestream-datasource/pull/209)
- Chore: Update version of code-coverage [#211](https://github.com/grafana/timestream-datasource/pull/211)
- Feature: Timestream is now available in us-gov-west-1 [#207](https://github.com/grafana/timestream-datasource/pull/207)

## 2.3.2

- Security: Upgrade Go in build process to 1.19.3

## 2.3.1

- Security: Upgrade Go in build process to 1.19.2

## v2.3.0

- Change timestamp fieldType to be nullable by @nekketsuuu in https://github.com/grafana/timestream-datasource/pull/184
- Upgrade to grafana-aws-sdk v0.11.0 by @fridgepoet in https://github.com/grafana/timestream-datasource/pull/195

## v2.2.0

- Add support for context aware autocompletion by @sunker in https://github.com/grafana/timestream-datasource/pull/188

## v2.1.0

- Add 'ap-southeast-2' and 'ap-northeast-1' regions [#178](https://github.com/grafana/timestream-datasource/pull/178)

## v2.0.1

- Bug fix for issue logging in with incorrect keys: https://github.com/grafana/timestream-datasource/pull/176
- Code Coverage Check updates

## v2.0.0

- Breaking Change: Timestream data source now requires Grafana 8.0+ to run.
- Fix: Allow null data points for time series [#170](https://github.com/grafana/timestream-datasource/pull/170)

## v1.5.2

- Fix Panic while parsing null timestamps (#165)

## v1.5.1

- Always apply double quotes to database and table name (#155)

## v1.5.0

- Revamp query editor.
- Add toggle to avoid streaming responses.
- Add `$__interval` variable.
- Modify the User-Agent for requests. Now it will follow this form: `"aws-sdk-go/$aws-sdk-version ($go-version; $OS;) Timestream/$timestream-version-$git-hash Grafana/$grafana-version"`.
- Fixes bugs for Endpoint and Assume Role settings.

## v1.4.0

- Add macros for raw values of interval, from, to [#98](https://github.com/grafana/timestream-datasource/pull/98)
- Quote and join multiple variables [#118](https://github.com/grafana/timestream-datasource/pull/118)
- Add stats for bytes metered and scanned [#110](https://github.com/grafana/timestream-datasource/pull/110)

## v1.3.3

- Support for multiple timeseries columns
- Improved support for custom endpoint

## v1.3.2

- Adding eu-central-1 region
- renamed "master" branch to "main"
- build with Golang 1.6

## v1.3.1

- Execute each query in its own request, this will support multiple queries that
  require multiple pages to complete
- Upgrade shared authentication library
- Bump minimum grafana runtime to 7.5

## v1.3.0

- fix bug with supporting multi-page timeseries results
- Use a shared authentication library and UI component
- Bump minimum grafana runtime to 7.4

## v1.2.0

- Support $\_\_timefilter on armhf (#52, @mg-arne)
- Add $\_\_now_ms macro (#49, @squalou)
- Fixed region picker default values

## v1.1.2

- Fix template variable queries
- Only show valid regions

## v1.1.1

- Avoid double escaping
- support template variables in query

## v1.1.0

- Updated authentication to match builtin cloudwatch authentication
- Include query status in metadata
- Examples and query suggestions now quote all names

## v1.0.0

- Initial Release
