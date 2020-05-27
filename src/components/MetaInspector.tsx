import React, { PureComponent } from 'react';
import { MetadataInspectorProps, DataFrame } from '@grafana/data';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions } from '../types';

export type Props = MetadataInspectorProps<DataSource, TimestreamQuery, TimestreamOptions>;

export class MetaInspector extends PureComponent<Props> {
  state = { index: 0 };

  renderInfo = (frame: DataFrame) => {
    const query = frame.meta?.custom?.executedQuery as string;
    if (!query) {
      return null;
    }

    return (
      <div>
        <h3>Executed Query</h3>
        <pre>{query}</pre>
      </div>
    );
  };

  render() {
    const { data } = this.props;
    if (!data || !data.length) {
      return <div>No Data</div>;
    }
    return (
      <div>
        {data.map(frame => {
          return this.renderInfo(frame);
        })}
      </div>
    );
  }
}
