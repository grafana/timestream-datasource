import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface TimestreamQuery extends DataQuery {
  rawQuery?: string;
  noTruncation?: boolean;
}

export interface TimestreamOptions extends DataSourceJsonData {
  // nothing for now
}

export interface TimestreamSecureJsonData {
  // nothing for now
}
