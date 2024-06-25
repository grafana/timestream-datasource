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
  getDatabases: () => Promise<TableDefinition[]>;
  getTables: (d?: string) => Promise<TableDefinition[]>;
  getColumns: (database?: string, table?: string) => Promise<ColumnDefinition[]>;
}

export const getTimestreamCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider =
  ({ getDatabases, getTables, getColumns }) =>
  (monaco, language) => {
    return {
      // get standard SQL completion provider which will resolve functions and macros
      ...(language && getStandardSQLCompletionProvider(monaco, language)),
      triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
      schemas: {
        resolve: getDatabases,
      },
      tables: {
        resolve: (t?: TableIdentifier | null) => {
          return getTables(t?.schema);
        },
        parseName: (token?: LinkedToken | null) => {
          const tablePath = token?.value ?? '';
          const parts = tablePath.split('.');

          return {
            schema: parts.length >= 1 && parts[0] ? parts[0] : undefined,
            table: parts.length >= 2 && parts[1] ? parts[1] : undefined,
          };
        },
      },
      columns: {
        resolve: (t?: TableIdentifier) => getColumns(t?.schema, t?.table),
      },
      supportedMacros: () => MACROS,
    };
  };
