import { SQLEditor as SQLCodeEditor } from '@grafana/experimental';
import { DataSource } from '../DataSource';
import { getTimestreamCompletionProvider } from 'language/completionItemProvider';
import { DATABASE_MACRO, TABLE_MACRO } from 'language/macros';
import React, { useRef, useMemo, useCallback, useEffect } from 'react';
import { TimestreamQuery } from 'types';
import timestreamLanguageDefinition from 'language/definition';

interface RawEditorProps {
  query: TimestreamQuery;
  onRunQuery: () => void;
  onChange: (q: string) => void;
  datasource: DataSource;
}

export default function SQLEditor({ query, datasource, onRunQuery, onChange }: RawEditorProps) {
  const queryRef = useRef<TimestreamQuery>(query);
  useEffect(() => {
    queryRef.current = query;
  }, [query]);

  const getDatabases = useCallback(async () => {
    const databases: string[] = await datasource.postResource('databases').catch(() => []);
    return databases.map((database) => ({ name: database, completion: database }));
  }, [datasource]);

  const getTables = useCallback(
    async (database?: string) => {
      const tables: string[] = await datasource
        .postResource('tables', {
          database: database ?? queryRef.current.database ?? '',
        })
        .catch(() => []);
      return tables.map((table) => ({ name: table, completion: table }));
    },
    [datasource]
  );

  const getColumns = useCallback(
    async (database?: string, tableName?: string) => {
      const interpolatedArgs = {
        database: database ? database.replace(TABLE_MACRO, queryRef.current.table ?? '') : queryRef.current.table,
        table: tableName
          ? tableName.replace(DATABASE_MACRO, queryRef.current.database ?? '')
          : queryRef.current.database,
      };
      const [measures, dimensions] = await Promise.all([
        datasource.postResource('measures', interpolatedArgs).catch(() => []),
        datasource.postResource('dimensions', interpolatedArgs).catch(() => []),
      ]);
      return [...measures, ...dimensions].map((column) => ({ name: column, completion: column }));
    },
    [datasource]
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
      onBlur={() => onRunQuery()}
      onChange={(rawQuery) => onChange(rawQuery)}
      language={{
        ...timestreamLanguageDefinition,
        completionProvider,
      }}
    ></SQLCodeEditor>
  );
}
