# Change Log

All notable changes to this project will be documented in this file.

## 2.12.5

- Publish to every environment on Cloud and Add On-Prem promotion in [#525](https://github.com/grafana/timestream-datasource/pull/525)
- Chore(deps): Bump the aws-sdk-go-v2 group with 2 updates in [#500](https://github.com/grafana/timestream-datasource/pull/500)
- Chore(deps): Bump actions/stale from 9 to 10 in [#486](https://github.com/grafana/timestream-datasource/pull/486)

## 2.12.4

- Bump transitive dep brace-expansion to 1.1.12 and 2.0.2 in [#522](https://github.com/grafana/timestream-datasource/pull/522)
- Update gitignore for yarn files in [#512](https://github.com/grafana/timestream-datasource/pull/512)
- Upgrade yarn to v4.10.3 in [#507](https://github.com/grafana/timestream-datasource/pull/507)
- Chore: Move docs to docs/sources to enable publishing to website in [#508](https://github.com/grafana/timestream-datasource/pull/508)
- CD: Publish Github release directly, skipping draft step in [#506](https://github.com/grafana/timestream-datasource/pull/506)
- Update workflows and templates in [#494](https://github.com/grafana/timestream-datasource/pull/494)
- Chore(deps-dev): Bump the npm-patch-dev-dependencies group with 5 updates in [#502](https://github.com/grafana/timestream-datasource/pull/502)
- Chore(deps): Bump github.com/grafana/grafana-aws-sdk from 1.1.1 to 1.2.0 in [#503](https://github.com/grafana/timestream-datasource/pull/503)
- Chore(deps): Bump github.com/grafana/grafana-plugin-sdk-go from 0.278.0 to 0.279.0 in [#501](https://github.com/grafana/timestream-datasource/pull/501)
- Update dependabot groups in [#493](https://github.com/grafana/timestream-datasource/pull/493)
- Chore(deps-dev): Bump @swc/core from 1.12.7 to 1.13.5 in [#495](https://github.com/grafana/timestream-datasource/pull/495)
- Chore(deps-dev): Bump sass from 1.89.2 to 1.92.0 in [#492](https://github.com/grafana/timestream-datasource/pull/492)
- Chore(deps-dev): Bump jest-environment-jsdom from 30.0.2 to 30.1.2 in [#489](https://github.com/grafana/timestream-datasource/pull/489)

## 2.12.3

- Remove pr_commands by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/436
- Chore: Use ci github action to run e2e tests by @idastambuk in https://github.com/grafana/timestream-datasource/pull/438
- Add support for auto-merging dependabot updates by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/427
- Tweak dependabot schedule by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/450
- Bump github.com/aws/aws-sdk-go-v2 from 1.36.4 to 1.36.6 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/442
- Bump @swc/jest from 0.2.38 to 0.2.39 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/444
- Remove @types/glob by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/445
- Bump dotenv from 17.2.0 to 17.2.1 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/446
- Bump github.com/grafana/grafana-plugin-sdk-go from 0.277.1 to 0.278.0 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/447
- Bump @eslint/js from 9.30.0 to 9.31.0 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/449
- Bump github.com/aws/aws-sdk-go-v2/service/timestreamquery from 1.31.1 to 1.31.3 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/443
- Bump github.com/grafana/grafana-aws-sdk from 0.38.7 to 1.1.0 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/451
- Stop running queries onblur by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/453
- Bump the grafana-dependencies group and run grafana/create-plugin update by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/448
- Bump eslint-plugin-prettier from 5.5.1 to 5.5.4 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/458
- Bump actions/checkout from 4 to 5 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/461
- Bump github.com/grafana/grafana-aws-sdk from 1.1.0 to 1.1.1 by @dependabot[bot] in https://github.com/grafana/timestream-datasource/pull/464
- Fix: disable endpoint discovery when a custom endpoint is set by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/475

## 2.12.2

- Releasing for internal deployment testing, no actual changes

## 2.12.1

- Mark legacy query format as downstream error by @njvrzm in [#431](https://github.com/grafana/timestream-datasource/pull/431)
- Chore: migrate to Github actions by @idastambuk in [#425](https://github.com/grafana/timestream-datasource/pull/425)
- Chore: Update CONTRIBUTING.md with release instructions by @idastambuk in [#426](https://github.com/grafana/timestream-datasource/pull/426)
- Dependency updates:
  - Bump the all-go-dependencies group across 1 directory with 3 updates by @dependabot in [#423](https://github.com/grafana/timestream-datasource/pull/423)

## 2.12.0

- Bump the all-node-dependencies group across 1 directory with 31 updates in [#412](https://github.com/grafana/timestream-datasource/pull/412)
- Bump github.com/grafana/grafana-aws-sdk from 0.38.3 to 0.38.4 in the all-go-dependencies group in [#417](https://github.com/grafana/timestream-datasource/pull/417)
- Bump the all-go-dependencies group across 1 directory with 4 updates in [#410](https://github.com/grafana/timestream-datasource/pull/410)
- Update for new v2-style auth in grafana-aws-sdk in [#392](https://github.com/grafana/timestream-datasource/pull/392)
- Bump golang.org/x/net from 0.36.0 to 0.38.0 in the go_modules group in [#398](https://github.com/grafana/timestream-datasource/pull/398)
- Update for aws-sdk-go-v2 in [#356](https://github.com/grafana/timestream-datasource/pull/356)

## 2.11.0

- Bump github.com/grafana/grafana-plugin-sdk-go from 0.268.1 to 0.269.1 in the all-go-dependencies group across 1 directory in [#384](https://github.com/grafana/timestream-datasource/pull/384)
- Bump golang.org/x/net from 0.35.0 to 0.36.0 in the go_modules group in [#385](https://github.com/grafana/timestream-datasource/pull/385)
- Bump the all-node-dependencies group across 1 directory with 18 updates in [#387](https://github.com/grafana/timestream-datasource/pull/387)
- Bump @babel/runtime from 7.26.9 to 7.26.10 in the npm_and_yarn group in [#386](https://github.com/grafana/timestream-datasource/pull/386)
- Add PDC support in [#374](https://github.com/grafana/timestream-datasource/pull/374)

## 2.10.1

- Update minimum Grafana supported version in plugin.json

## 2.10.0

- Migrate form to new styling in [#375](https://github.com/grafana/timestream-datasource/pull/375)
- Add external PRs to project board in [#366](https://github.com/grafana/timestream-datasource/pull/366)
- Chore: add label to external contributions in [#362](https://github.com/grafana/timestream-datasource/pull/362)
- Migrate E2E tests to Playwright in [#358](https://github.com/grafana/timestream-datasource/pull/358)
- Dependabot:
  - Bump the all-node-dependencies group across 1 directory with 7 updates in [#377](https://github.com/grafana/timestream-datasource/pull/377)
  - Bump the all-go-dependencies group across 1 directory with 3 updates in [#376](https://github.com/grafana/timestream-datasource/pull/376)
  - Bump github.com/grafana/grafana-plugin-sdk-go from 0.265.0 to 0.266.0 in the all-go-dependencies group in [#372](https://github.com/grafana/timestream-datasource/pull/372)
  - Bump the all-node-dependencies group across 1 directory with 22 updates in [#370](https://github.com/grafana/timestream-datasource/pull/370)
  - Bump github.com/grafana/grafana-plugin-sdk-go from 0.262.0 to 0.265.0 in the all-go-dependencies group across 1 directory in [#368](https://github.com/grafana/timestream-datasource/pull/368)

## 2.9.13

- Bump the all-node-dependencies group across 1 directory with 21 updates in [#352](https://github.com/grafana/timestream-datasource/pull/352)
- Bump the all-go-dependencies group across 1 directory with 3 updates in [#351](https://github.com/grafana/timestream-datasource/pull/351)

## 2.9.12

- Bump cross-spawn from 7.0.3 to 7.0.6 in the npm_and_yarn group in [#330](https://github.com/grafana/timestream-datasource/pull/330)
- Bump the all-go-dependencies group across 1 directory with 2 updates in [#334](https://github.com/grafana/timestream-datasource/pull/334)
- Add ap-south-1 to regions list in [#331](https://github.com/grafana/timestream-datasource/pull/331)
- Bump the all-node-dependencies group across 1 directory with 30 updates in [#337](https://github.com/grafana/timestream-datasource/pull/337)

## 2.9.11

- Bugfix: interpolate interval on the backend [#327](https://github.com/grafana/timestream-datasource/pull/327)

## 2.9.10

- Bugfix: Account for template variable being a number
- Chore: update dependabot config (#317)
- Dependency updates:
  - github.com/grafana/grafana-plugin-sdk-go from 0.251.0 to 0.258.0 in [#314](https://github.com/grafana/timestream-datasource/pull/314),[#315](https://github.com/grafana/timestream-datasource/pull/315), [#319](https://github.com/grafana/timestream-datasource/pull/319)
  - github.com/aws/aws-sdk-go from 1.51.31 to 1.55.5 in [#319](https://github.com/grafana/timestream-datasource/pull/319)
  - github.com/grafana/grafana-aws-sdk from 0.31.2 to 0.31.4 in [#319](https://github.com/grafana/timestream-datasource/pull/319)
  - actions/checkout from 2 to 4 in [#318](https://github.com/grafana/timestream-datasource/pull/318)
  - tibdex/github-app-token from 1.8.0 to 2.1.0 in [#318](https://github.com/grafana/timestream-datasource/pull/318)
  - github.com/grafana/sqlds/v4 from v4.1.0 to v4.1.2 in [#322](https://github.com/grafana/timestream-datasource/pull/322)

## 2.9.9

- Fix "Wait for All Queries" toggle in [#313](https://github.com/grafana/timestream-datasource/pull/313)
- Fix errors in LongToWide transformation in [#311](https://github.com/grafana/timestream-datasource/pull/311)
- Chore: Update plugin.json keywords in [#310](https://github.com/grafana/timestream-datasource/pull/310)
- Update grafana-plugin-sdk-go and grafana-aws-sdk in [#309](https://github.com/grafana/timestream-datasource/pull/309)
- fix: linter complaints in [#308](https://github.com/grafana/timestream-datasource/pull/308)
- Bump path-to-regexp from 1.8.0 to 1.9.0 in [#303](https://github.com/grafana/timestream-datasource/pull/303)
- Docs: Updates and improvements in [#302](https://github.com/grafana/timestream-datasource/pull/302)
- Add dependabot for grafana/plugin-sdk-go in [#307](https://github.com/grafana/timestream-datasource/pull/307)

## 2.9.8

- Bump webpack from 5.92.1 to 5.94.0 in [#301](https://github.com/grafana/timestream-datasource/pull/301)
- Bump micromatch from 4.0.7 to 4.0.8 in [#299](https://github.com/grafana/timestream-datasource/pull/299)
- Bump fast-loops from 1.1.3 to 1.1.4 in [#298](https://github.com/grafana/timestream-datasource/pull/298)

## 2.9.7

- feat: add errorsource [#296](https://github.com/grafana/timestream-datasource/pull/296)
- chore: refactor macros to avoid macro-length bug in [#295](https://github.com/grafana/timestream-datasource/pull/295)

## 2.9.6

- Bugfix: Fix $interval variable interpolation in [#291](https://github.com/grafana/timestream-datasource/pull/291)

## 2.9.5

- Chore: update dependencies in [#290](https://github.com/grafana/timestream-datasource/pull/290)

## 2.9.4

- Fix: use ReadAuthSettings to get authSettings in [#289](https://github.com/grafana/timestream-datasource/pull/289)

## 2.9.3

- Upgrade grafana-aws-sdk and other packages [#285](https://github.com/grafana/timestream-datasource/pull/285)

## 2.9.2

- Add keywords by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/278
- Bring in [security fixes in go 1.21.8](https://groups.google.com/g/golang-announce/c/5pwGVUPoMbg)

## 2.9.1

- Update grafana/aws-sdk to 0.20.0 to add a new supported region in [#274](https://github.com/grafana/timestream-datasource/pull/274)
- Query Editor: Fix table and database mapping in [#272](https://github.com/grafana/timestream-datasource/pull/272)

## 2.9.0

- Bump jest-dom in [#270](https://github.com/grafana/timestream-datasource/pull/270)
- Query Editor: Stop running query automatically when all macros are selected in [#269](https://github.com/grafana/timestream-datasource/pull/269)
- Bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.45.0 to 0.46.0 in [#267](https://github.com/grafana/timestream-datasource/pull/267)

## 2.8.0

- Support Node 18 by @kevinwcyu in https://github.com/grafana/timestream-datasource/pull/245
- Bump word-wrap from 1.2.3 to 1.2.5 by @dependabot in https://github.com/grafana/timestream-datasource/pull/246
- Bump semver from 5.7.1 to 5.7.2 by @dependabot in https://github.com/grafana/timestream-datasource/pull/235
- Bump golang.org/x/net from 0.9.0 to 0.17.0 by @dependabot in https://github.com/grafana/timestream-datasource/pull/250
- Bump postcss from 8.4.18 to 8.4.31 by @dependabot in https://github.com/grafana/timestream-datasource/pull/249
- Bump @babel/traverse from 7.18.13 to 7.23.2 by @dependabot in https://github.com/grafana/timestream-datasource/pull/252
- Bump google.golang.org/grpc from 1.54.0 to 1.56.3 by @dependabot in https://github.com/grafana/timestream-datasource/pull/253
- Bump go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace from 0.37.0 to 0.44.0 by @dependabot in https://github.com/grafana/timestream-datasource/pull/251
- Bump @babel/traverse from 7.18.13 to 7.23.2 by @dependabot in https://github.com/grafana/timestream-datasource/pull/255
- Bump google.golang.org/grpc from 1.58.2 to 1.58.3 by @dependabot in https://github.com/grafana/timestream-datasource/pull/254
- Upgrade underscore, d3-color, debug, cosmiconfig, yaml dependencies by @fridgepoet in https://github.com/grafana/timestream-datasource/pull/263
- Bump yaml from 2.1.3 to 2.3.4 by @dependabot in https://github.com/grafana/timestream-datasource/pull/264

**Full Changelog**: https://github.com/grafana/timestream-datasource/compare/v2.7.1...v2.8.0

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
