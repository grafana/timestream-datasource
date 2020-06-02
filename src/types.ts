import { DataQuery } from '@grafana/data';
import { AwsDataSourceJsonData, AwsDataSourceSecureJsonData } from 'common/types';

export interface TimestreamQuery extends DataQuery {
  rawQuery?: string;
  noTruncation?: boolean;
}

export interface TimestreamOptions extends AwsDataSourceJsonData {
  // nothing for now
}

export interface TimestreamSecureJsonData extends AwsDataSourceSecureJsonData {
  // nothing for now
}
