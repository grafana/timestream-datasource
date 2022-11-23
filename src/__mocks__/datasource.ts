import { DataSourcePluginOptionsEditorProps, PluginType } from '@grafana/data';

import { DataSource } from '../DataSource';
import { TimestreamOptions, TimestreamQuery } from '../types';

export const mockDatasource = new DataSource({
  id: 1,
  uid: 'timestream-id',
  type: 'timestream-datasource',
  name: 'Timestream Data Source',
  jsonData: {},
  access: 'direct',
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
  readOnly: false
});

export const mockDatasourceOptions: DataSourcePluginOptionsEditorProps<TimestreamOptions> = {
  options: {
    id: 1,
    uid: '1',
    orgId: 1,
    name: 'Timestream',
    typeLogoUrl: '',
    type: '',
    access: '',
    url: '',
    user: '',
    basicAuth: false,
    basicAuthUser: '',
    database: '',
    isDefault: false,
    jsonData: {
      defaultRegion: 'us-east-2',
    },
    secureJsonFields: {},
    readOnly: false,
    withCredentials: false,
    typeName: '',
  },
  onOptionsChange: jest.fn(),
};

export const mockQuery: TimestreamQuery = { rawQuery: 'select * from foo', refId: '' };
