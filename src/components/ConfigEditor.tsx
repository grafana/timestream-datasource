import { ConfigSelect, ConnectionConfig } from '@grafana/aws-sdk';
import { DataSourcePluginOptionsEditorProps, DataSourceSettings, SelectableValue } from '@grafana/data';
import { getBackendSrv } from '@grafana/runtime';
import React, { useState } from 'react';
import { standardRegions } from 'regions';

import { TimestreamDataSourceSettings, TimestreamOptions, TimestreamSecureJsonData } from '../types';
import { selectors } from './selectors';

export type Props = DataSourcePluginOptionsEditorProps<TimestreamOptions, TimestreamSecureJsonData>;

export type ResourceType = 'defaultDatabase' | 'defaultTable' | 'defaultMeasure';

export function ConfigEditor(props: Props) {
  const baseURL = `/api/datasources/${props.options.id}`;
  const resourcesURL = `${baseURL}/resources`;
  const [saved, setSaved] = useState(!!props.options.jsonData.defaultRegion);
  const saveOptions = async () => {
    if (saved) {
      return;
    }
    await getBackendSrv()
      .put(baseURL, props.options)
      .then((result: { datasource: TimestreamDataSourceSettings }) => {
        props.onOptionsChange({
          ...props.options,
          version: result.datasource.version,
        });
      });
    setSaved(true);
  };

  // Databases
  const fetchDatabases = async () => {
    const loaded: string[] = await getBackendSrv().get(resourcesURL + '/databases');
    return loaded;
  };
  // Tables
  const fetchTables = async () => {
    const loaded: string[] = await getBackendSrv().post(resourcesURL + '/tables', {
      database: props.options.jsonData.defaultDatabase,
    });
    return loaded;
  };
  // Measures
  const fetchMeasures = async () => {
    const loadedWorkgroups: string[] = await getBackendSrv().post(resourcesURL + '/measures', {
      database: props.options.jsonData.defaultDatabase,
      table: props.options.jsonData.defaultTable,
    });
    return loadedWorkgroups;
  };

  const onOptionsChange = (options: DataSourceSettings<TimestreamOptions, TimestreamSecureJsonData>) => {
    setSaved(false);
    props.onOptionsChange(options);
  };

  const onChange = (resource: ResourceType) => (e: SelectableValue<string> | null) => {
    let value = e?.value ?? '';
    props.onOptionsChange({
      ...props.options,
      jsonData: {
        ...props.options.jsonData,
        [resource]: value,
      },
    });
  };

  return (
    <div className="gf-form-group">
      <ConnectionConfig
        {...props}
        standardRegions={standardRegions}
        defaultEndpoint="https://query-{cell}.timestream.{region}.amazonaws.com"
        onOptionsChange={onOptionsChange}
      />
      <h3>Timestream Details</h3>
      <p>Default values to be used as macros</p>
      <ConfigSelect
        {...props}
        value={props.options.jsonData.defaultDatabase ?? ''}
        onChange={onChange('defaultDatabase')}
        fetch={fetchDatabases}
        label={selectors.components.ConfigEditor.defaultDatabase.input}
        data-testid={selectors.components.ConfigEditor.defaultDatabase.wrapper}
        saveOptions={saveOptions}
      />
      <ConfigSelect
        {...props}
        value={props.options.jsonData.defaultTable ?? ''}
        onChange={onChange('defaultTable')}
        fetch={fetchTables}
        label={selectors.components.ConfigEditor.defaultTable.input}
        data-testid={selectors.components.ConfigEditor.defaultTable.wrapper}
        dependencies={[props.options.jsonData.defaultDatabase || '']}
        saveOptions={saveOptions}
      />
      <ConfigSelect
        {...props}
        value={props.options.jsonData.defaultMeasure ?? ''}
        onChange={onChange('defaultMeasure')}
        fetch={fetchMeasures}
        label={selectors.components.ConfigEditor.defaultMeasure.input}
        data-testid={selectors.components.ConfigEditor.defaultMeasure.wrapper}
        dependencies={[props.options.jsonData.defaultDatabase || '', props.options.jsonData.defaultTable || '']}
        saveOptions={saveOptions}
      />
    </div>
  );
}
