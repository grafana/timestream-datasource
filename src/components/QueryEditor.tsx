import { ResourceSelector } from '@grafana/aws-sdk';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { CodeEditor, InlineField, InlineSegmentGroup, Label, Select, Switch } from '@grafana/ui';
import React, { useEffect, useState } from 'react';

import { DataSource } from '../DataSource';
import { TimestreamOptions, TimestreamQuery } from '../types';
import { sampleQueries } from './samples';
import { selectors } from './selectors';
import { getSuggestions } from './Suggestions';

type Props = QueryEditorProps<DataSource, TimestreamQuery, TimestreamOptions>;

type QueryProperties = 'database' | 'table' | 'measure';

export function QueryEditor(props: Props) {
  const { query, datasource } = props;
  const { database, table, measure } = query;
  const { defaultDatabase, defaultTable, defaultMeasure } = datasource.options;

  // pre-populate query with default data
  useEffect(() => {
    if (!database || !table || !table) {
      props.onChange({
        ...query,
        database: database || defaultDatabase,
        table: table || defaultTable,
        measure: measure || defaultMeasure,
      });
    }
  }, []);

  const onWaitForChange = () => {
    props.onChange({ ...query, waitForResult: !query.waitForResult });
  };

  const onChangeSelector = (prop: QueryProperties) => (e: SelectableValue<string> | null) => {
    props.onChange({ ...query, [prop]: e?.value });
  };

  const onQueryChange = (rawQuery: string) => {
    props.onChange({ ...query, rawQuery });
    props.onRunQuery();
  };

  // Trigger query if all the resources are set
  useEffect(() => {
    if (database && table && measure) {
      props.onRunQuery();
    }
  }, [database, table, measure]);

  // Databases used both for the selector and editor suggestions
  const [dbs, setDBs] = useState<string[]>([]);
  useEffect(() => {
    datasource.getResource('databases').then((res) => setDBs(res));
  }, []);

  // Tables used both for the selector and editor suggestions
  const [tables, setTables] = useState<string[]>([]);
  useEffect(() => {
    if (database) {
      datasource
        .postResource('tables', {
          database: database || '',
        })
        .then((res: string[]) => {
          if (res.length > 0 && !res.some((t) => table === t)) {
            // The current list of tables do not include the current one
            // so change it to the first of the list
            props.onChange({ ...query, table: res[0] });
          }
          setTables(res);
        });
    }
  }, [database]);

  // Measures used both for the selector and editor suggestions
  const [measures, setMeasures] = useState<string[]>([]);
  // Dimensions used for editor suggestions
  const [dimensions, setDimensions] = useState<string[]>([]);
  useEffect(() => {
    if (database && table) {
      datasource
        .postResource('measures', {
          database: database,
          table: table,
        })
        .then((res: string[]) => {
          if (res.length > 0 && !res.some((t) => measure === t)) {
            // The current list of measures do not include the current one
            // so change it to the first of the list
            props.onChange({ ...query, measure: res[0] });
          }
          setMeasures(res);
        });
      datasource
        .postResource('dimensions', {
          database: database,
          table: table,
        })
        .then((res) => {
          setDimensions(res);
        });
    }
  }, [database, table]);

  return (
    <InlineSegmentGroup>
      <div className="gf-form-group">
        <h6>Macros</h6>
        <ResourceSelector
          onChange={onChangeSelector('database')}
          resources={dbs}
          value={database || null}
          tooltip="Use the selected schema with the $__database macro"
          label={selectors.components.ConfigEditor.defaultDatabase.input}
          data-testid={selectors.components.ConfigEditor.defaultDatabase.wrapper}
          labelWidth={11}
          className="width-12"
        />
        <ResourceSelector
          onChange={onChangeSelector('table')}
          resources={tables}
          value={table || null}
          tooltip="Use the selected table with the $__table macro"
          label={selectors.components.ConfigEditor.defaultTable.input}
          data-testid={selectors.components.ConfigEditor.defaultTable.wrapper}
          labelWidth={11}
          className="width-12"
        />
        <ResourceSelector
          onChange={onChangeSelector('measure')}
          resources={measures}
          value={measure || null}
          tooltip="Use the selected column with the $__measure macro"
          label={selectors.components.ConfigEditor.defaultMeasure.input}
          data-testid={selectors.components.ConfigEditor.defaultMeasure.wrapper}
          labelWidth={11}
          className="width-12"
        />
        <h6>Render</h6>
        <InlineField
          id={`${props.query.refId}-wait`}
          label={'Wait for all queries'}
          labelWidth={15}
          style={{ alignItems: 'center' }}
        >
          <Switch
            aria-labelledby={`${props.query.refId}-wait`}
            onChange={onWaitForChange}
            checked={query.waitForResult}
          />
        </InlineField>
        <h6>Sample queries</h6>
        <Label description={'Selecting a sample will modify the current query'}>
          <InlineField label="Query" labelWidth={11}>
            <Select
              aria-label={'Query'}
              options={sampleQueries}
              onChange={(e: SelectableValue<string>) => onQueryChange(e.value || '')}
              className="width-12"
            />
          </InlineField>
        </Label>
      </div>
      <div style={{ minWidth: '400px', marginLeft: '10px', flex: 1 }}>
        <CodeEditor
          language={'sql'}
          value={query.rawQuery || ''}
          onBlur={onQueryChange}
          showMiniMap={false}
          showLineNumbers={true}
          getSuggestions={() => getSuggestions(dbs, tables, measures, dimensions, database, table, measure)}
          height="240px"
        />
      </div>
    </InlineSegmentGroup>
  );
}
