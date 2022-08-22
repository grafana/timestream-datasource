import {
  ColumnDefinition,
  getStandardSQLCompletionProvider,
  LanguageCompletionProvider,
  LinkedToken,
  TableDefinition,
  TableIdentifier,
} from '@grafana/experimental';
import { MACROS } from './macros';

interface CompletionProviderGetterArgs {
  getDatabases: React.MutableRefObject<() => Promise<TableDefinition[]>>;
  getTables: React.MutableRefObject<(d?: string) => Promise<TableDefinition[]>>;
  getColumns: React.MutableRefObject<(database?: string, table?: string) => Promise<ColumnDefinition[]>>;
}

export const getTimestreamCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getDatabases, getTables, getColumns }) =>
  (monaco, language) => {
    return {
      // get standard SQL completion provider which will resolve functions and macros
      ...(language && getStandardSQLCompletionProvider(monaco, language)),
      triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
      schemas: {
        resolve: () => getDatabases.current(),
      },
      tables: {
        resolve: (t: TableIdentifier) => {
          return getTables.current(t?.schema);
        },
        parseName: (token: LinkedToken) => {
          const tablePath = token?.value ?? '';
          const parts = tablePath.split('.');

          return {
            schema: parts.length >= 1 && parts[0] ? parts[0] : undefined,
            table: parts.length >= 2 && parts[1] ? parts[1] : undefined,
          };
        },
      },
      columns: {
        resolve: (t: TableIdentifier) => getColumns.current(t.schema, t.table),
      },
      supportedMacros: () => MACROS,
    };
  };
