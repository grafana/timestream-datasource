import { SQLEditor as SQLCodeEditor } from '@grafana/plugin-ui';
import { DataSource } from '../DataSource';
import { getTimestreamCompletionProvider } from 'language/completionItemProvider';
import { DATABASE_MACRO, TABLE_MACRO } from 'language/macros';
import React, { useMemo, useCallback } from 'react';
import { TimestreamQuery } from 'types';
import timestreamLanguageDefinition from 'language/definition';

interface RawEditorProps {
  query: TimestreamQuery;
  onRunQuery: () => void;
  onChange: (query: TimestreamQuery) => void;
  datasource: DataSource;
}

export default function SQLEditor({ query, datasource, onChange }: RawEditorProps) {
  const onChangeRawQuery = (rawQuery: string) => {
    onChange({ ...query, rawQuery });
  };

  const getDatabases = useCallback(async () => {
    const databases: string[] = await datasource.postResource<string[]>('databases').catch(() => []);
    return databases.map((database) => ({ name: database, completion: database }));
  }, [datasource]);

  const getTables = useCallback(
    async (database?: string) => {
      const tables: string[] = await datasource
        .postResource<string[]>('tables', {
          database: database ?? query.database ?? '',
        })
        .catch(() => []);
      return tables.map((table) => ({ name: table, completion: table }));
    },
    [datasource, query.database]
  );

  const getColumns = useCallback(
    async (database?: string, tableName?: string) => {
      const interpolatedArgs = {
        database: database
          ? database.replace(DATABASE_MACRO, query.database ?? '')
          : query.database,
        table: tableName ? tableName.replace(TABLE_MACRO, query.table ?? '') : query.table,
      };
      const [measures, dimensions] = await Promise.all([
        datasource.postResource<string[]>('measures', interpolatedArgs).catch(() => []),
        datasource.postResource<string[]>('dimensions', interpolatedArgs).catch(() => []),
      ]);
      return [...measures, ...dimensions].map((column) => ({ name: column, completion: column }));
    },
    [datasource, query.database, query.table]
  );

  const completionProvider = useMemo(
    () =>
      getTimestreamCompletionProvider({
        getDatabases,
        getTables,
        getColumns,
      }),
    [getDatabases, getTables, getColumns]
  );

  return (
    <SQLCodeEditor
      query={query.rawQuery ?? ''}
      onChange={onChangeRawQuery}
      language={{
        ...timestreamLanguageDefinition,
        completionProvider,
      }}
    ></SQLCodeEditor>
  );
}
