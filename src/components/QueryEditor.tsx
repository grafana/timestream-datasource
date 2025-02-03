import { ResourceSelector } from '@grafana/aws-sdk';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { InlineField, InlineSegmentGroup, Label, Select, Switch } from '@grafana/ui';
import React, { useEffect, useState } from 'react';

import { DataSource } from '../DataSource';
import { FormatOptions, SelectableFormatOptions, TimestreamOptions, TimestreamQuery } from '../types';
import { sampleQueries } from './samples';
import { selectors } from './selectors';
import SQLEditor from './SQLEditor';

type Props = QueryEditorProps<DataSource, TimestreamQuery, TimestreamOptions>;

type QueryProperties = 'database' | 'table' | 'measure';

export function QueryEditor(props: Props) {
  const { query, datasource, onChange, onRunQuery } = props;
  const { database, table, measure, format } = query;
  const { defaultDatabase, defaultTable, defaultMeasure } = datasource.options;

  // pre-populate query with default data
  useEffect(() => {
    if (!database || !table || !measure) {
      onChange({
        ...query,
        database: database || defaultDatabase,
        table: table || defaultTable,
        measure: measure || defaultMeasure,
        format: format || FormatOptions.Table,
      });
    }
    // Run only once
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onWaitForChange = () => {
    onChange({ ...query, waitForResult: !query.waitForResult });
  };

  const onChangeSelector = (prop: QueryProperties) => (e: SelectableValue | null) => {
    onChange({ ...query, [prop]: e?.value });
  };

  const onChangeFormat = (e: SelectableValue) => {
    onChange({ ...query, format: e.value || 0 });
    onRunQuery();
  };

  const onQueryChange = (rawQuery: string) => {
    onChange({ ...query, rawQuery });
    onRunQuery();
  };

  // Databases used both for the selector and editor suggestions
  const [databases, setDatabases] = useState<string[]>([]);
  useEffect(() => {
    datasource.getResource('databases').then((res) => setDatabases(res));
  }, [datasource]);

  // Tables used both for the selector and editor suggestions
  const [tables, setTables] = useState<string[]>([]);
  useEffect(() => {
    if (database) {
      datasource
        .postResource<string[]>('tables', {
          database: database || '',
        })
        .then((res: string[]) => {
          if (res.length > 0 && !res.some((t) => table === t)) {
            // The current list of tables do not include the current one
            // so change it to the first of the list
            onChange({ ...query, table: res[0] });
          }
          setTables(res);
        });
    }
    // Run only on database change
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [database]);

  // Measures used both for the selector and editor suggestions
  const [measures, setMeasures] = useState<string[]>([]);
  // Dimensions used for editor suggestions
  useEffect(() => {
    if (database && table) {
      datasource
        .postResource<string[]>('measures', {
          database: database,
          table: table,
        })
        .then((res: string[]) => {
          if (res.length > 0 && !res.some((t) => measure === t)) {
            // The current list of measures do not include the current one
            // so change it to the first of the list
            onChange({ ...query, measure: res[0] });
          }
          setMeasures(res);
        });
    }
    // Run only on database or table change
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [database, table]);

  return (
    <InlineSegmentGroup>
      <div className="gf-form-group">
        <h6>Macros</h6>
        <ResourceSelector
          id="database"
          onChange={onChangeSelector('database')}
          resources={databases}
          value={database || null}
          tooltip="Use the selected schema with the $__database macro"
          label={selectors.components.ConfigEditor.defaultDatabase.input}
          data-testid={selectors.components.ConfigEditor.defaultDatabase.wrapper}
          labelWidth={11}
          className="width-12"
        />
        <ResourceSelector
          id="table"
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
          id="measure"
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
          labelWidth={16}
          style={{ alignItems: 'center' }}
          htmlFor={`${props.query.refId}-wait-for-all-queries`}
        >
          <Switch
            id={`${props.query.refId}-wait-for-all-queries`}
            aria-labelledby={`${props.query.refId}-wait`}
            onChange={onWaitForChange}
            value={query.waitForResult}
          />
        </InlineField>
        <h6>Frames</h6>
        <InlineField
          label="Format as"
          htmlFor="format-as"
          labelWidth={13}
          tooltip={
            <>
              {
                'Timeseries queries must have times in ascending order, which can be done by adding "ORDER BY <time field> ASC" to the query. '
              }
              <a
                href="https://docs.aws.amazon.com/timestream/latest/developerguide/supported-sql-constructs.SELECT.html"
                target="_blank"
                rel="noopener noreferrer"
              >
                See the AWS Docs for more details.
              </a>
            </>
          }
          interactive
        >
          <Select
            inputId="format-as"
            options={SelectableFormatOptions}
            value={props.query.format || FormatOptions.Table}
            onChange={onChangeFormat}
            className="width-11"
            menuShouldPortal={true}
          />
        </InlineField>
        <h6>Sample queries</h6>
        <Label description={'Selecting a sample will modify the current query'}>
          <InlineField label="Query" labelWidth={11} htmlFor={`${props.query.refId}-sample-query`}>
            <Select
              aria-label={'Query'}
              inputId={`${props.query.refId}-sample-query`}
              options={sampleQueries}
              onChange={(e: SelectableValue) => onQueryChange(e.value || '')}
              className="width-12"
            />
          </InlineField>
        </Label>
      </div>
      <div style={{ minWidth: '400px', marginLeft: '10px', flex: 1 }}>
        <SQLEditor
          query={query}
          onRunQuery={props.onRunQuery}
          onChange={props.onChange}
          datasource={props.datasource}
        />
      </div>
    </InlineSegmentGroup>
  );
}
