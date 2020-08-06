import React, { PureComponent } from 'react';
import { MetadataInspectorProps, DataFrame } from '@grafana/data';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions } from '../types';

export type Props = MetadataInspectorProps<DataSource, TimestreamQuery, TimestreamOptions>;

export class MetaInspector extends PureComponent<Props> {
  state = { index: 0 };

  renderInfo = (frame: DataFrame, idx: number) => {
    const custom = frame.meta?.custom;
    if (!custom) {
      return null;
    }

    return (
      <div key={idx}>
        <h3>Query ID</h3>
        <pre>{custom.queryId}</pre>
        {custom.nextToken && (
          <>
            <h3>Next Token</h3>
            <pre>{custom.nextToken}</pre>
          </>
        )}

        <h3>Details</h3>
        <pre>{JSON.stringify(custom, null, 2)}</pre>
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
        {data.map((frame, idx) => {
          return this.renderInfo(frame, idx);
        })}
      </div>
    );
  }
}
