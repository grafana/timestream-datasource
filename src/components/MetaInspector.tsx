import React, { PureComponent } from 'react';
import { MetadataInspectorProps, DataFrame } from '@grafana/data';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions } from '../types';

export type Props = MetadataInspectorProps<DataSource, TimestreamQuery, TimestreamOptions>;

export class MetaInspector extends PureComponent<Props> {
  state = { index: 0 };

  renderInfo = (frame: DataFrame) => {
    const custom = frame.meta?.custom;
    if (!custom) {
      return null;
    }

    // meta["queryId"] = output.QueryId
    // meta["nextToken"] = output.NextToken

    return (
      <div>
        <h3>Query ID</h3>
        <pre>{custom.queryId}</pre>
        {custom.nextToken && (
          <>
            <h3>Next Token</h3>
            <pre>{custom.nextToken}</pre>
          </>
        )}

        <h3>Executed Query</h3>
        <pre>{(frame.meta as any).executedQueryString}</pre>
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
