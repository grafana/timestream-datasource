import { AwsAuthDataSourceJsonData, AwsAuthDataSourceSecureJsonData } from '@grafana/aws-sdk';
import { DataSourceSettings, SelectableValue } from '@grafana/data';
import { type DataQuery } from '@grafana/schema';

export interface ColumnInfo {
  column: string;
  type: string;
  category: string;
}

export enum StandardFunction {
  Value = 'value',
  Avg = 'avg',
  Min = 'min',
  Max = 'mean',
  P90 = 'p90',
  P95 = 'p95',
  P99 = 'p99',
}

export enum DataType {
  varchar = 'varchar',
  double = 'double',
  timestamp = 'timestamp',
}

export enum FormatOptions {
  Table,
  TimeSeries,
}

export const SelectableFormatOptions: Array<SelectableValue<FormatOptions>> = [
  {
    label: 'Table',
    value: FormatOptions.Table,
  },
  {
    label: 'Time Series',
    value: FormatOptions.TimeSeries,
  },
];

export interface MeasureInfo {
  name: string;
  type: DataType;
  dimensions: string[]; // only the strings for now
}

export interface SchemaInfo {
  databases?: string[];
  tables?: string[];
  measures?: MeasureInfo[];
}

export interface TimestreamCustomMeta {
  queryId: string;
  nextToken?: string;
  hasSeries?: boolean;

  executionStartTime?: number; // The backend clock
  executionFinishTime?: number; // The backend clock

  fetchStartTime?: number; // The frontend clock
  fetchEndTime?: number; // The frontend clock
  fetchTime?: number; // The frontend clock

  status: {
    CumulativeBytesMetered?: number;
    CumulativeBytesScanned?: number;
  };

  // when multiple queries exist we keep track of each request
  subs?: TimestreamCustomMeta[];
}

export interface TimestreamQuery extends DataQuery {
  // When specified, use this rather than the default for macros
  database?: string;
  table?: string;
  measure?: string;

  // The rendered query
  rawQuery?: string;

  // Avoid pagination
  waitForResult?: boolean;

  format?: FormatOptions;

  // Not a real parameter...
  // nextToken?: string;
}

export interface TimestreamOptions extends AwsAuthDataSourceJsonData {
  defaultDatabase?: string;
  defaultTable?: string;
  defaultMeasure?: string;
}

export interface TimestreamSecureJsonData extends AwsAuthDataSourceSecureJsonData {
  // nothing for now
}

export type TimestreamDataSourceSettings = DataSourceSettings<TimestreamOptions, TimestreamSecureJsonData>;
