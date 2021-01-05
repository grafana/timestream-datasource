import React, { PureComponent } from 'react';
import { InlineFormLabel, AsyncSelect } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { TimestreamOptions, TimestreamSecureJsonData, TimestreamQuery } from '../types';
import ConnectionConfig from '../common/ConnectionConfig';
import { SchemaInfo } from 'SchemaInfo';
import { getDataSourceSrv } from '@grafana/runtime';
import { DataSource } from '../DataSource';

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
    console.log('componentDidMountX', this.props);
    const d = await getDataSourceSrv().get(this.props.options.name);
    const { options } = this.props;
    const ds = d as DataSource;
    const query = {
      refId: 'X',
      database: options.jsonData.defaultDatabase,
      table: options.jsonData.defaultTable,
      measure: options.jsonData.defaultMeasure,
    };
    this.setState({ schema: new SchemaInfo(ds, query) });
  };

  // Try harder to get the state
  componentDidUpdate = async () => {
    if (!this.state.schema) {
      console.log('no state... try again');
      this.componentDidMount();
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

    return (
      <div>
        <br />
        <h3>Default Query Macros</h3>
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel className={widthKey}>$__database</InlineFormLabel>
            <AsyncSelect
              className={widthVal}
              loadOptions={schema.getDatabases}
              value={currentDatabase}
              onChange={this.onDatabaseChange}
              defaultOptions
              loadingMessage="..."
              allowCustomValue={true}
            />
          </div>
        </div>
        {defaultDatabase && (
          <div className="gf-form-inline">
            <div className="gf-form">
              <InlineFormLabel className={widthKey}>$__table</InlineFormLabel>
              <AsyncSelect
                className={widthVal}
                loadOptions={schema.getTables}
                value={currentTable}
                onChange={this.onTableChange}
                defaultOptions
                loadingMessage="..."
                allowCustomValue={true}
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
                loadOptions={schema.getMeasures}
                value={currentMeasure}
                onChange={this.onMeasureChange}
                defaultOptions
                loadingMessage="..."
                allowCustomValue={true}
                formatCreateLabel={v => `Use unknown measure: ${v}`}
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
          <ConnectionConfig {...this.props} defaultEndpoint="https://query-{cell}.timestream.{region}.amazonaws.com" />
        </div>

        {schema && this.renderDefaultChoices(schema)}
      </>
    );
  }
}
