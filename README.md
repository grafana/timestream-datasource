# AWS Timestream Datasource Development Guide

The Timestream datasource plugin provides a support for [Amazon Timestream](https://aws.amazon.com/timestream/). Add it as a data source, then you are ready to build dashboards using timestream query results

Please add feedback to the [issues](https://github.com/grafana/timestream-datasource/issues) folder, and we will follow up shortly.  Be sure to include version information for both grafana and the installed plugin.

The production plugins can be downloaded from [the Timestream plugin page](https://grafana.com/grafana/plugins/grafana-timestream-datasource/installation).

For configuration options, see: [src/README.md](src/README.md)


## Developer Guide


You will need to install [Node.js](https://nodejs.org/en/), [Yarn](https://yarnpkg.com/), [Go](https://golang.org/), and [Mage](https://magefile.org/) first.
1. `yarn install --frozen-lockfile`
1. `yarn dev` — will build the frontend changes
1. `mage build:backend` — will build the backend changes
1. (Optional) `mage -v buildAll` — this is optional if you need backend plugins for other platforms
1. The compiled plugin should be in `dist/` directory.
1. Run Grafana in [development](https://grafana.com/docs/grafana/latest/administration/configuration/#app_mode) mode, or configure Grafana to [load the unsigned plugin](https://grafana.com/docs/grafana/latest/plugins/plugin-signatures/#allow-unsigned-plugins).
1. You can install by following the [install Grafana plugins docs page](https://grafana.com/docs/grafana/latest/plugins/installation/).

For more information, please consult the [build a plugin docs page](https://grafana.com/docs/grafana/latest/developers/plugins/).
