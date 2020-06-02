import { DataSourceJsonData, SelectableValue } from '@grafana/data';

export enum AwsAuthType {
  Keys = 'keys',
  Credentials = 'credentials',
  ARN = 'arn',
}

export interface AwsDataSourceJsonData extends DataSourceJsonData {
  authType?: AwsAuthType;
  assumeRoleArn?: string;
  profile?: string; // Credentials profile name, as specified in ~/.aws/credentials
  defaultRegion?: string; // region if it is not defined by your credentials file
}

export interface AwsDataSourceSecureJsonData {
  accessKey?: string;
  secretKey?: string;
}

export const awsAuthProviderOptions = [
  { label: 'Access & secret key', value: AwsAuthType.Keys },
  { label: 'Credentials file', value: AwsAuthType.Credentials },
  { label: 'ARN', value: AwsAuthType.ARN },
] as Array<SelectableValue<AwsAuthType>>;

export const standardRegions = [
  'ap-east-1',
  'ap-northeast-1',
  'ap-northeast-2',
  'ap-northeast-3',
  'ap-south-1',
  'ap-southeast-1',
  'ap-southeast-2',
  'ca-central-1',
  'cn-north-1',
  'cn-northwest-1',
  'eu-central-1',
  'eu-north-1',
  'eu-west-1',
  'eu-west-2',
  'eu-west-3',
  'me-south-1',
  'sa-east-1',
  'us-east-1',
  'us-east-2',
  'us-gov-east-1',
  'us-gov-west-1',
  'us-iso-east-1',
  'us-isob-east-1',
  'us-west-1',
  'us-west-2',
];
