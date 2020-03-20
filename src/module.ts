import { DataSourcePlugin } from "@grafana/data";
import { DataSource } from "./DataSource";
import { QueryEditor, ConfigEditor } from "./components";
import { TimestreamQuery, TimestreamOptions } from "./types";

export const plugin = new DataSourcePlugin<
  DataSource,
  TimestreamQuery,
  TimestreamOptions
>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
