# Building and releasing

## How to build the Timestream data source plugin locally

## Dependencies

Make sure you have the following dependencies installed first:

- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/) (see [go.mod](../go.mod#L3) for minimum required version)
- [Mage](https://magefile.org/)
- [Node.js (Long Term Support)](https://nodejs.org)
- [Yarn](https://yarnpkg.com)

## Frontend

1. Install dependencies

   ```bash
   yarn install --pure-lockfile
   ```

2. Build plugin in development mode or run in watch mode

   ```bash
   yarn dev
   ```

   or

   ```bash
   yarn watch
   ```

3. Build plugin in production mode

   ```bash
   yarn build
   ```

## Backend

1. Build the backend binaries

   ```bash
   mage -v
   ```

## E2E Tests

1. `yarn playwright install --with-deps`
1. `yarn server`
1. `yarn e2e`

### Golden files

Golden files check that data frames are being generated correctly based on the Timestream API response. They have two parts, the json files represent the raw API response and the golden files represent the expected data frame. Both are generated in executor_test.go.

#### Re-generating json API response

> **Note:** Only members of the Grafana team can re-generate these files. If you need help with this, ping the `@grafana/aws-datasources` team on GitHub and they will help out.

1. Make sure to comment out the [t.Skip("Integration Test")](https://github.com/grafana/timestream-datasource/blob/5b3f07edb13cb3e3bbeeca284f5b9228a30de451/pkg/timestream/executor_test.go#L64) line in the executor_test.go file.
2. Run the `TestGenerateTestData`. This should regenerate the json files.
3. Uncomment the `t.Skip("Integration Test")` again.

#### Re-generating golden files

1. Change the last argument in the [CheckGoldenDataResponse](https://github.com/grafana/timestream-datasource/blob/5b3f07edb13cb3e3bbeeca284f5b9228a30de451/pkg/timestream/executor_test.go#L40) call to true. This will re-generate the golden files.
2. Run the test, and then undo the change from step 1.
3. Re-run the test and they should now pass.

## Releasing

1. Update the version number in the `package.json` file.
2. Update the `CHANGELOG.md` with the changes contained in the release.
3. Commit the changes to master and push to GitHub.
4. Follow the release process that you can find [here](https://enghub.grafana-ops.net/docs/default/component/grafana-plugins-platform/plugins-ci-github-actions/010-plugins-ci-github-actions/#cd_1)
