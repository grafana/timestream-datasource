import { DataSourcePluginOptionsEditorProps, PluginType } from '@grafana/data';
import { TimestreamOptions, TimestreamQuery } from '../types';
import { DataSource } from '../DataSource';

export const mockDatasource = new DataSource({
  id: 1,
  uid: 'timestream-id',
  type: 'timestream-datasource',
  name: 'Timestream Data Source',
  jsonData: {},
  meta: {
    id: 'timestream-datasource',
    name: 'Timestream Data Source',
    type: PluginType.datasource,
    module: '',
    baseUrl: '',
    info: {
      description: '',
      screenshots: [],
      updated: '',
      version: '',
      logos: {
        small: '',
        large: '',
      },
      author: {
        name: '',
      },
      links: [],
    },
  },
});

export const mockDatasourceOptions: DataSourcePluginOptionsEditorProps<TimestreamOptions> = {
  options: {
    id: 1,
    orgId: 1,
    name: 'Redshift',
    typeLogoUrl: '',
    type: '',
    access: '',
    url: '',
    password: '',
    user: '',
    basicAuth: false,
    basicAuthPassword: '',
    basicAuthUser: '',
    database: '',
    isDefault: false,
    jsonData: {
      defaultRegion: 'us-east-2',
    },
    secureJsonFields: {},
    readOnly: false,
    withCredentials: false,
  },
  onOptionsChange: jest.fn(),
};

export const mockQuery: TimestreamQuery = { rawQuery: 'select * from foo', refId: '' };
