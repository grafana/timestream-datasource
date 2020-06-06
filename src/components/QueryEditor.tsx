import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import Editor from '@monaco-editor/react';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions, QueryType, MeasureInfo } from '../types';
import { config, getTemplateSrv } from '@grafana/runtime';
import { QueryField } from './Forms';

import { Segment, SegmentAsync, InlineFormLabel } from '@grafana/ui';
import { sampleQueries, queryTypes } from './samples';

type Props = QueryEditorProps<DataSource, TimestreamQuery, TimestreamOptions>;
interface State {
  measures?: Array<SelectableValue<MeasureInfo>>;
}

export class QueryEditor extends PureComponent<Props, State> {
  getEditorValue: any | undefined;

  state: State = {};

  //-----------------------------------------------------
  //-----------------------------------------------------

  onRawQueryChange = () => {
    const rawQuery = this.getEditorValue();
    if (this.props.query.rawQuery === rawQuery) {
      return; // no change
    }

    this.props.onChange({
      ...this.props.query,
      rawQuery,
      queryType: QueryType.Raw,
    });
    this.props.onRunQuery();
  };

  onEditorDidMount = (getEditorValue: any) => {
    this.getEditorValue = getEditorValue;
  };

  getCurrentVars() {
    let { database, table } = this.props.query;
    const templateSrv = getTemplateSrv();

    if (isVar(database)) {
      database = templateSrv.replace(database!);
    } else if (!database) {
      database = templateSrv.replace('${database}');
    }

    if (isVar(table)) {
      table = templateSrv.replace(table!);
    } else if (!table) {
      table = templateSrv.replace('${table}');
    }

    return {
      database,
      table,
    };
  }

  getDatabaseOptions = async (q?: string) => {
    const vals = await this.props.datasource.getDatabases();
    const opts = vals.map(v => {
      return { value: v, label: v };
    });
    if (this.props.query.database) {
      opts.push({ value: '', label: '-- remove --' });
    }
    return opts;
  };

  getTablesOptions = async (q?: string) => {
    const vars = this.getCurrentVars();
    const vals = await this.props.datasource.getTables(vars.database);
    const opts = vals.map(v => {
      return { value: v, label: v };
    });
    if (this.props.query.table) {
      opts.push({ value: '', label: '-- remove --' });
    }
    return opts;
  };

  getMeasureOptions = async (q?: string) => {
    const vars = this.getCurrentVars();
    const vals = await this.props.datasource.getMeasures(vars.database, vars.table);
    const opts = vals.map(v => {
      return { value: v, label: v };
    });
    if (this.props.query.table) {
      opts.push({ value: '', label: '-- remove --' });
    }
    return opts;
  };

  onDbChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      database: value.value,
    };
    if (!query.database) {
      delete query.database;
    }

    // Remove the table when Db changes
    if (!isVar(query.table)) {
      delete query.table;
    }

    this.props.onChange(query);
    this.props.onRunQuery();
  };

  onTableChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      table: value.value,
    };
    if (!query.table) {
      delete query.table;
    }

    this.props.onChange(query);
    console.log('Table Changed!', query);
  };

  onMeasureChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      measure: value.value,
    };
    if (!query.measure) {
      delete query.measure;
    }

    this.props.onChange(query);
    console.log('Measure Changed!', query);
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
    this.props.onRunQuery();
  };

  render() {
    const { query } = this.props;
    const queryType = queryTypes.find(v => v.value === query.queryType) || queryTypes[1]; // Samples

    const vars = {
      database: '${database}',
      table: '${table}',
      measure: '${measure}',
    };

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
          <InlineFormLabel width={8} className="query-keyword">
            Variables
          </InlineFormLabel>
          <InlineFormLabel className="keyword" width={6}>
            {vars.database}
          </InlineFormLabel>
          <SegmentAsync
            value={query.database || '+'}
            loadOptions={this.getDatabaseOptions}
            onChange={this.onDbChanged}
            allowCustomValue
          />

          {query.database && (
            <>
              <InlineFormLabel className="keyword" width={5}>
                {vars.table}
              </InlineFormLabel>
              <SegmentAsync
                value={query.table || '+'}
                loadOptions={this.getTablesOptions}
                onChange={this.onTableChanged}
                allowCustomValue
              />
            </>
          )}

          {query.table && (
            <>
              <InlineFormLabel className="keyword" width={5}>
                {vars.measure}
              </InlineFormLabel>
              <SegmentAsync
                value={query.measure || '+'}
                loadOptions={this.getMeasureOptions}
                onChange={this.onMeasureChanged}
                allowCustomValue
              />
            </>
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

function isVar(txt?: string): boolean {
  return !!txt && txt.startsWith('${');
}
