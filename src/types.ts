import { DataQuery, KeyValue } from '@grafana/data';
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
  bigint = 'bigint',
  timestamp = 'timestamp',
}

export interface MeasureInfo {
  name: string;
  type: DataType;
  dimensions: KeyValue<DataType>;
}

export interface TimestreamQuery extends DataQuery {
  // Standard templates
  database?: string;
  table?: string;
  measure?: string; // single measure

  // The rendered query
  rawQuery?: string;
}

export interface TimestreamOptions extends AwsDataSourceJsonData {
  // nothing for now
}

export interface TimestreamSecureJsonData extends AwsDataSourceSecureJsonData {
  // nothing for now
}
