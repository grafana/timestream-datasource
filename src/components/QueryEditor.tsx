import React, { PureComponent, ChangeEvent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { TextArea } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions } from '../types';

type Props = QueryEditorProps<DataSource, TimestreamQuery, TimestreamOptions>;

export class QueryEditor extends PureComponent<Props> {
  onRawQueryChange = (event: ChangeEvent<any>) => {
    this.props.onChange({
      ...this.props.query,
      rawQuery: event.target.value,
    });
  };

  toggleNoTruncation = (event?: React.SyntheticEvent<HTMLInputElement>) => {
    const { query, onChange, onRunQuery } = this.props;
    onChange({
      ...query,
      noTruncation: !query.noTruncation,
    });
    onRunQuery();
  };

  render() {
    const { query } = this.props;
    return (
      <>
        <div>
          <TextArea
            className="gf-form-input"
            rows={15}
            value={query.rawQuery}
            onChange={this.onRawQueryChange}
            placeholder="timestream query"
          />
        </div>
      </>
    );
  }
}
