# Change Log

All notable changes to this project will be documented in this file.

## v2.0.0
_Not released yet_

- Breaking Change: Timestream data source now requires Grafana 8.0+ to run.

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
