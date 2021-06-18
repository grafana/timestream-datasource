# Change Log

All notable changes to this project will be documented in this file.

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
- Support $__timefilter on armhf (#52, @mg-arne)
- Add $__now_ms macro (#49, @squalou)
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
