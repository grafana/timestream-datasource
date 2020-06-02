import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { TimestreamOptions, TimestreamSecureJsonData } from '../types';
import CommonConfig from '../common/CommonConfig';

export type Props = DataSourcePluginOptionsEditorProps<TimestreamOptions, TimestreamSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  render() {
    return (
      <div>
        <CommonConfig {...this.props} />
      </div>
    );
  }
}
