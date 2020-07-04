import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../DataSource';
import { TimestreamQuery, TimestreamOptions, QueryType } from '../types';
import { getTemplateSrv } from '@grafana/runtime';
import { QueryField } from './Forms';

import { Segment, SegmentAsync, InlineFormLabel, CodeEditor } from '@grafana/ui';
import { sampleQueries, queryTypes } from './samples';
import { SchemaInfo } from 'SchemaInfo';

type Props = QueryEditorProps<DataSource, TimestreamQuery, TimestreamOptions>;
interface State {
  schema?: SchemaInfo;

  schemaState?: Partial<TimestreamQuery>;
}

export class QueryEditor extends PureComponent<Props, State> {
  state: State = {};

  componentDidMount = () => {
    const { datasource, query } = this.props;

    const schema = new SchemaInfo(datasource, query, getTemplateSrv());
    this.setState({ schema: schema, schemaState: schema.state });

    schema.preload().then(v => {
      console.log('Loaded schema');
    });
  };

  //-----------------------------------------------------
  //-----------------------------------------------------

  onQueryChange = (rawQuery: string) => {
    this.props.onChange({
      ...this.props.query,
      rawQuery,
      queryType: QueryType.Raw,
    });
    this.props.onRunQuery();
  };

  onDbChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      database: value.value,
    };
    if (!query.database) {
      delete query.database;
    }

    const { schema } = this.state;
    const schemaState = schema!.updateState(query);
    this.setState({ schemaState });
    this.props.onChange(query);

    if (schemaState.table) {
      this.props.onRunQuery();
    }
  };

  onTableChanged = (value: SelectableValue<string>) => {
    const query = {
      ...this.props.query,
      table: value.value,
    };
    if (!query.table) {
      delete query.table;
    }
    const { schema } = this.state;
    const schemaState = schema!.updateState(query);
    this.setState({ schemaState });

    this.props.onChange(query);
    this.props.onRunQuery();
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

    const { schema } = this.state;
    const schemaState = schema!.updateState(query);
    this.setState({ schemaState });
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

  renderDatabaseMacro = (schema: SchemaInfo, query?: string, value?: string) => {
    let placehoder = '';
    let current = '$__database = ';
    if (query) {
      current += query;
    } else {
      placehoder = current + (value ?? '?');
      current = '';
    }

    return (
      <SegmentAsync
        value={current}
        loadOptions={schema.getDatabases}
        placeholder={placehoder}
        onChange={this.onDbChanged}
        allowCustomValue
      />
    );
  };

  renderTableMacro = (schema: SchemaInfo, query?: string, value?: string) => {
    let placehoder = '';
    let current = '$__table = ';
    if (query) {
      current += query;
    } else {
      placehoder = current + (value ?? '?');
      current = '';
    }

    return (
      <SegmentAsync
        value={current}
        loadOptions={schema.getTables}
        placeholder={placehoder}
        onChange={this.onTableChanged}
        allowCustomValue
      />
    );
  };

  renderMeasureMacro = (schema: SchemaInfo, query?: string, value?: string) => {
    let placehoder = '';
    let current = '$__measure = ';
    if (query) {
      current += query;
    } else {
      placehoder = current + (value ?? '?');
      current = '';
    }

    return (
      <SegmentAsync
        value={current}
        loadOptions={schema.getMeasures}
        placeholder={placehoder}
        onChange={this.onMeasureChanged}
        allowCustomValue
      />
    );
  };

  render() {
    const { query } = this.props;
    const { schema, schemaState } = this.state;

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
          <InlineFormLabel width={8} className="query-keyword">
            Macros
          </InlineFormLabel>
          {schema && schemaState && (
            <>
              {this.renderDatabaseMacro(schema, query.database, schemaState.database)}
              {this.renderTableMacro(schema, query.table, schemaState.table)}
              {this.renderMeasureMacro(schema, query.measure, schemaState.measure)}
            </>
          )}

          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
        {queryType.value !== QueryType.Builder && schema && (
          <CodeEditor
            height={'250px'}
            language="sql"
            value={query.rawQuery || ''}
            onBlur={this.onQueryChange}
            onSave={this.onQueryChange}
            showMiniMap={false}
            showLineNumbers={true}
            getSuggestions={schema.getSuggestions}
          />
        )}
      </>
    );
  }
}
