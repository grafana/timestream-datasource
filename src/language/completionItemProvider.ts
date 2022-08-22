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

export const getTimestreamCompletionProvider: (args: CompletionProviderGetterArgs) => LanguageCompletionProvider = ({
  getDatabases,
  getTables,
  getColumns,
}) => (monaco, language) => {
  return {
    // get standard SQL completion provider which will resolve functions and macros
    ...(language && getStandardSQLCompletionProvider(monaco, language)),
    triggerCharacters: ['.', ' ', '$', ',', '(', "'"],
    schemas: {
      resolve: async () => getDatabases.current(),
    },
    tables: {
      resolve: async (t: TableIdentifier) => {
        return await getTables.current(t?.schema);
      },
      parseName: (token: LinkedToken) => {
        let tablePath = token?.value ?? '';

        const parts = tablePath.split('.');
        if (parts.length === 1) {
          return { schema: parts[0] || undefined };
        } else if (parts.length === 2) {
          return { schema: parts[0] || undefined, table: parts[1] || undefined };
        }

        return null;
      },
    },
    columns: {
      resolve: async (t: TableIdentifier) => getColumns.current(t.schema, t.table),
    },
    supportedMacros: () => MACROS,
  };
};
