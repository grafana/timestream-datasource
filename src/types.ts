import { DataQuery } from '@grafana/data';
import { AwsDataSourceJsonData, AwsDataSourceSecureJsonData } from 'common/types';

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

export enum QueryType {
  Builder = 'builder',
  Samples = 'samples',
  Raw = 'raw',
}

export enum DataType {
  varchar = 'varchar',
  double = 'double',
  timestamp = 'timestamp',
}

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

  // Not a real parameter...
  // nextToken?: string;
}

export interface TimestreamOptions extends AwsDataSourceJsonData {
  defaultDatabase?: string;
  defaultTable?: string;
  defaultMeasure?: string;
}

export interface TimestreamSecureJsonData extends AwsDataSourceSecureJsonData {
  // nothing for now
}
