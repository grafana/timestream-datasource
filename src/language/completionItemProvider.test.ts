import { LinkedToken, SQLMonarchLanguage } from '@grafana/experimental';
import { getTimestreamCompletionProvider } from './completionItemProvider';

describe('completionItemProvider', () => {
  describe('when parsing table identifier', () => {
    const language: SQLMonarchLanguage = {
      id: 'custom-grafana-sql-language',
      tokenizer: {},
      builtinFunctions: ['SUM', 'AVG'],
      logicalOperators: ['AND', 'OR'],
      comparisonOperators: ['=', '!='],
    };
    const completionProvider = getTimestreamCompletionProvider({} as any)({} as any, language);
    const SCHEMA = '"testSchema"';
    const TABLE = '"testDatabase"';

    it('should resolve complete identifier', () => {
      const { schema, table } = completionProvider.tables!.parseName!({
        value: `${SCHEMA}.${TABLE}`,
      } as LinkedToken);

      expect(schema).toEqual(SCHEMA);
      expect(table).toEqual(TABLE);
    });

    it('should resolve schema when missing dot and missing table', () => {
      const { schema, table } = completionProvider.tables!.parseName!({
        value: `${SCHEMA}`,
      } as LinkedToken);

      expect(schema).toEqual(SCHEMA);
      expect(table).toBeUndefined();
    });

    it('should resolve schema when missing table', () => {
      const { schema, table } = completionProvider.tables!.parseName!({
        value: `${SCHEMA}.`,
      } as LinkedToken);

      expect(schema).toEqual(SCHEMA);
      expect(table).toBeUndefined();
    });
  });
});
