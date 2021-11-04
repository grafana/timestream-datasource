import React, { PureComponent } from 'react';
import { InlineFormLabel, AsyncSelect } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { TimestreamOptions, TimestreamSecureJsonData, TimestreamQuery } from '../types';
import { SchemaInfo } from 'SchemaInfo';
import { getDataSourceSrv } from '@grafana/runtime';
import { DataSource } from '../DataSource';
import { ConnectionConfig } from '@grafana/aws-sdk';
import { standardRegions } from 'regions';

export type Props = DataSourcePluginOptionsEditorProps<TimestreamOptions, TimestreamSecureJsonData>;

interface State {
  schema?: SchemaInfo;
  schemaState?: Partial<TimestreamQuery>;
}

export class ConfigEditor extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  componentDidMount = async () => {
    const d = await getDataSourceSrv().get(this.props.options.name);
    const { options } = this.props;
    const ds = (d as unknown) as DataSource;
    const query = {
      refId: 'X',
      database: options.jsonData.defaultDatabase,
      table: options.jsonData.defaultTable,
      measure: options.jsonData.defaultMeasure,
    };
    this.setState({ schema: new SchemaInfo(ds, query) });
  };

  // Try harder to get the state
  componentDidUpdate = async (oldProps: Props) => {
    if (!this.state.schema) {
      console.log('no state... try again');
      this.componentDidMount();
    } else {
      const oldData = oldProps.options.jsonData;
      const { jsonData } = this.props.options;
      if (jsonData.authType !== oldData.authType || jsonData.defaultRegion !== jsonData.defaultRegion) {
        this.componentDidMount(); // update
      }
    }
  };

  onDatabaseChange = (value: SelectableValue<string>) => {
    const { options } = this.props;
    const jsonData = {
      ...options.jsonData,
      defaultDatabase: value.value,
    };
    delete jsonData.defaultTable;
    delete jsonData.defaultMeasure;

    this.props.onOptionsChange({
      ...options,
      jsonData,
    });

    const { schema } = this.state;
    const schemaState = schema!.updateState({ database: jsonData.defaultDatabase ?? '' });
    this.setState({ schemaState });
  };

  onTableChange = (value: SelectableValue<string>) => {
    const { options } = this.props;
    const jsonData = {
      ...options.jsonData,
      defaultTable: value.value,
    };
    delete jsonData.defaultMeasure;

    this.props.onOptionsChange({
      ...options,
      jsonData,
    });

    const { schema } = this.state;
    const schemaState = schema!.updateState({ table: jsonData.defaultTable ?? '' });
    this.setState({ schemaState });
  };

  onMeasureChange = (value: SelectableValue<string>) => {
    const { options } = this.props;
    const jsonData = {
      ...options.jsonData,
      defaultMeasure: value.value,
    };

    this.props.onOptionsChange({
      ...options,
      jsonData,
    });

    const { schema } = this.state;
    const schemaState = schema!.updateState({ measure: jsonData.defaultMeasure ?? '' });
    this.setState({ schemaState });
  };

  renderDefaultChoices(schema: SchemaInfo) {
    const widthKey = 'width-14';
    const widthVal = 'width-30';

    const { options } = this.props;
    const { defaultDatabase, defaultTable, defaultMeasure } = options.jsonData;

    const currentDatabase = defaultDatabase
      ? { label: defaultDatabase, value: defaultDatabase }
      : { label: 'Select database', value: '' };

    const currentTable = defaultTable
      ? { label: defaultTable, value: defaultTable }
      : { label: 'Select table', value: '' };

    const currentMeasure = defaultMeasure
      ? { label: defaultMeasure, value: defaultMeasure }
      : { label: 'Select measure', value: '' };

    // Reload the dropdowns when config changes
    return (
      <div key={hashCode(JSON.stringify(this.props.options.jsonData))}>
        <br />
        <h3>Default Query Macros</h3>
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel className={widthKey}>$__database</InlineFormLabel>
            <AsyncSelect
              className={widthVal}
              cacheOptions={false}
              loadOptions={schema.getDatabases}
              value={currentDatabase}
              onChange={this.onDatabaseChange}
              defaultOptions
              loadingMessage="..."
              allowCustomValue={true}
              formatCreateLabel={(t) => `DB: ${t}`}
            />
          </div>
        </div>
        {defaultDatabase && (
          <div className="gf-form-inline">
            <div className="gf-form">
              <InlineFormLabel className={widthKey}>$__table</InlineFormLabel>
              <AsyncSelect
                className={widthVal}
                cacheOptions={false}
                loadOptions={schema.getTables}
                value={currentTable}
                onChange={this.onTableChange}
                defaultOptions
                loadingMessage="..."
                allowCustomValue={true}
                formatCreateLabel={(t) => `Table: ${t}`}
              />
            </div>
          </div>
        )}
        {defaultDatabase && defaultTable && (
          <div className="gf-form-inline">
            <div className="gf-form">
              <InlineFormLabel className={widthKey}>$__measure</InlineFormLabel>
              <AsyncSelect
                className={widthVal}
                cacheOptions={false}
                loadOptions={schema.getMeasures}
                value={currentMeasure}
                onChange={this.onMeasureChange}
                defaultOptions
                loadingMessage="..."
                allowCustomValue={true}
                formatCreateLabel={(v) => `Use unknown measure: ${v}`}
              />
            </div>
          </div>
        )}
      </div>
    );
  }

  render() {
    const { schema } = this.state;

    return (
      <>
        <div>
          <ConnectionConfig
            {...this.props}
            standardRegions={standardRegions}
            defaultEndpoint="https://query-{cell}.timestream.{region}.amazonaws.com"
          />
        </div>

        {schema && this.renderDefaultChoices(schema)}
      </>
    );
  }
}

function hashCode(s: string) {
  return s.split('').reduce((a, b) => {
    a = (a << 5) - a + b.charCodeAt(0);
    return a & a;
  }, 0);
}
