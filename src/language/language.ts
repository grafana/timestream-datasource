import { SQLMonarchLanguage, grafanaStandardSQLLanguage, grafanaStandardSQLLanguageConf } from '@grafana/experimental';

export const language: SQLMonarchLanguage = {
  ...grafanaStandardSQLLanguage,
  tokenizer: {
    ...grafanaStandardSQLLanguage.tokenizer,
    schemaTable: [
      ...grafanaStandardSQLLanguage.tokenizer.schemaTable,
      // recognize complete and incomplete database and table as identifier
      [/(\"\w+\"\.\"\w+\")/, 'identifier'], // e.g "database"."tablename"
      [/(\"\w+\"\.\"\w+)/, 'identifier'], // e.g "database"."tablename
      [/(\"\w+\"\.\w+)\"/, 'identifier'], // e.g "database".
      [/(\"\w+\"\.)/, 'identifier'], // e.g "database"
    ],
    complexIdentifiers: [], // if not resetting complexIdentifiers, database and table name would be recognized as identifier.quote and not identifier
  },
};

export const conf = grafanaStandardSQLLanguageConf;
