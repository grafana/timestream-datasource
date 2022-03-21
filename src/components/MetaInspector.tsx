import { DataFrame, MetadataInspectorProps } from '@grafana/data';
import React from 'react';

import { DataSource } from '../DataSource';
import { TimestreamOptions, TimestreamQuery } from '../types';

export type Props = MetadataInspectorProps<DataSource, TimestreamQuery, TimestreamOptions>;

export function MetaInspector(props: Props) {
  const { data } = props;

  const renderInfo = (frame: DataFrame, idx: number) => {
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

  if (!data || !data.length) {
    return <div>No Data</div>;
  }
  return (
    <div>
      {data.map((frame, idx) => {
        return renderInfo(frame, idx);
      })}
    </div>
  );
}
