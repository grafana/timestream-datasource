# Change Log

All notable changes to this project will be documented in this file.

## 2.3.3

- Upgrade @grafana/ dependencies to v9.2.4

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
- Upgrade shared authenticaiton library
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
