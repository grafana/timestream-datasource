import { DataSourcePlugin } from '@grafana/data';
import { MetaInspector } from 'components/MetaInspector';

import { ConfigEditor, QueryEditor } from './components';
import { DataSource } from './DataSource';
import { TimestreamOptions, TimestreamQuery } from './types';

export const plugin = new DataSourcePlugin<DataSource, TimestreamQuery, TimestreamOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setMetadataInspector(MetaInspector)
  .setQueryEditor(QueryEditor);
