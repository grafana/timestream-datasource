import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import Editor from '@monaco-editor/react';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions, QueryType, MeasureInfo } from '../types';
import { config } from '@grafana/runtime';
import { QueryField } from './Forms';

import { Segment } from '@grafana/ui';
import { sampleQueries, queryTypes } from './samples';

type Props = QueryEditorProps<DataSource, TimestreamQuery, TimestreamOptions>;
interface State {
  dbs?: Array<SelectableValue<string>>;
  tables?: Array<SelectableValue<string>>;
  measures?: Array<SelectableValue<MeasureInfo>>;
}

export class QueryEditor extends PureComponent<Props, State> {
  getEditorValue: any | undefined;

  state: State = {};

  onRawQueryChange = () => {
    this.props.onChange({
      ...this.props.query,
      rawQuery: this.getEditorValue(),
    });
    this.props.onRunQuery();
  };

  onEditorDidMount = (getEditorValue: any) => {
    this.getEditorValue = getEditorValue;
  };

  onDbChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      database: value.value,
    };
    this.props.onChange(query);
    console.log('DB Changed!', query);
  };

  onTableChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      table: value.value,
    };

    this.props.onChange(query);
    console.log('Table Changed!', query);
  };

  onQueryTypeChange = (value: SelectableValue<QueryType>) => {
    this.props.onChange({
      ...this.props.query,
      queryType: value.value || QueryType.Samples,
    });
  };

  onSampleChange = (value: SelectableValue<string>) => {
    this.props.onChange({
      ...this.props.query,
      rawQuery: value.value,
    });
  };

  render() {
    const { query } = this.props;
    const { dbs, tables } = this.state;
    const queryType = queryTypes.find(v => v.value === query.queryType) || queryTypes[1]; // Samples

    return (
      <>
        <div className={'gf-form-inline'}>
          <QueryField label="Query">
            <Segment value={queryType} options={queryTypes} onChange={this.onQueryTypeChange} />
          </QueryField>
          {queryType.value === QueryType.Samples && (
            <Segment value={''} placeholder="Select Example" options={sampleQueries} onChange={this.onSampleChange} />
          )}
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
        <div className={'gf-form-inline'}>
          <QueryField label="Database">
            <Segment
              value={query.database || '${database}'}
              options={dbs || []}
              onChange={this.onDbChanged}
              allowCustomValue
            />
          </QueryField>
          {query.database && (
            <QueryField label="Table">
              <Segment
                value={query.table || '${table}'}
                options={tables || []}
                onChange={this.onTableChanged}
                allowCustomValue
              />
            </QueryField>
          )}
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
        {queryType.value !== QueryType.Builder && (
          <div onBlur={this.onRawQueryChange}>
            <Editor
              height={'250px'}
              language="sql"
              value={query.rawQuery}
              editorDidMount={this.onEditorDidMount}
              theme={config.theme.isDark ? 'dark' : 'light'}
              options={{
                wordWrap: 'off',
                codeLens: false, // too small to bother
                minimap: {
                  enabled: false,
                  renderCharacters: false,
                },
                lineNumbersMinChars: 4,
                lineDecorationsWidth: 0,
                overviewRulerBorder: false,
                automaticLayout: true,
              }}
            />
          </div>
        )}
      </>
    );
  }
}
