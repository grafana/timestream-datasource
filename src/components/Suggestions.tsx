import { appendTemplateVariablesAsSuggestions } from '@grafana/aws-sdk';
import { getTemplateSrv } from '@grafana/runtime';
import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';

type Props = {
  databases: string[];
  tables: string[];
  measures: string[];
  dimensions: string[];
  database?: string;
  table?: string;
  measure?: string;
};

export const getSuggestions = ({ databases, tables, measures, dimensions, database, table, measure }: Props) => {
  const sugs: CodeEditorSuggestionItem[] = [
    {
      label: '$__timeFilter',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__interval_ms',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__now_ms',
      kind: CodeEditorSuggestionItemKind.Method,
      detail: '(Macro)',
    },
    {
      label: '$__database',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${database}`,
    },
    {
      label: '$__table',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${table}`,
    },
    {
      label: '$__measure',
      kind: CodeEditorSuggestionItemKind.Text,
      detail: `(Macro) ${measure}`,
    },
  ];

  databases.forEach((r) => sugs.push({ label: r, kind: CodeEditorSuggestionItemKind.Property, detail: '(Database)' }));
  tables.forEach((r) => sugs.push({ label: r, kind: CodeEditorSuggestionItemKind.Property, detail: '(Table)' }));
  measures.forEach((r) => sugs.push({ label: r, kind: CodeEditorSuggestionItemKind.Property, detail: '(Measure)' }));
  dimensions.forEach((r) =>
    sugs.push({ label: r, kind: CodeEditorSuggestionItemKind.Property, detail: '(Dimension)' })
  );

  return appendTemplateVariablesAsSuggestions(getTemplateSrv, sugs);
};
