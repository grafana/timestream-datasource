import { DataSourceInstanceSettings } from "@grafana/data";
import { DataSourceWithBackend } from "@grafana/runtime";

import { TimestreamQuery, TimestreamOptions } from "./types";

export class DataSource extends DataSourceWithBackend<
  TimestreamQuery,
  TimestreamOptions
> {
  constructor(
    instanceSettings: DataSourceInstanceSettings<TimestreamOptions>
  ) {
    super(instanceSettings);
  }
}
