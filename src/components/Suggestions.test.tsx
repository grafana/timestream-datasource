import { CodeEditorSuggestionItemKind } from '@grafana/ui';

import { getSuggestions } from './Suggestions';

const templateSrv = {
  getVariables: jest.fn().mockReturnValue([{ name: 'foo' }, { name: 'bar' }]),
  replace: jest.fn(),
};

jest.mock('@grafana/runtime', () => {
  return {
    ...(jest.requireActual('@grafana/runtime') as any),
    getTemplateSrv: () => templateSrv,
  };
});

const props = {
  databases: [],
  tables: [],
  measures: [],
  dimensions: [],
};

describe('getSuggestions', () => {
  it('should return default macros', () => {
    expect(getSuggestions(props)).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Macro)',
          label: '$__timeFilter',
          kind: CodeEditorSuggestionItemKind.Method,
        }),
      ])
    );
  });

  it('should return template variables', () => {
    expect(getSuggestions(props)).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Template Variable) undefined',
          label: '$foo',
          kind: CodeEditorSuggestionItemKind.Text,
        }),
      ])
    );
  });

  it('should include the $__database', () => {
    expect(getSuggestions({ ...props, database: 'db' })).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Macro) db',
          label: '$__database',
          kind: CodeEditorSuggestionItemKind.Text,
        }),
      ])
    );
  });

  it('should include the $__table', () => {
    expect(getSuggestions({ ...props, table: 'tab' })).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Macro) tab',
          label: '$__table',
          kind: CodeEditorSuggestionItemKind.Text,
        }),
      ])
    );
  });

  it('should include the $__measure', () => {
    expect(getSuggestions({ ...props, measure: 'cpu' })).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Macro) cpu',
          label: '$__measure',
          kind: CodeEditorSuggestionItemKind.Text,
        }),
      ])
    );
  });

  it('should return the list of databases', () => {
    expect(getSuggestions({ ...props, databases: ['foo', 'bar'] })).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Database)',
          label: 'foo',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
        expect.objectContaining({
          detail: '(Database)',
          label: 'bar',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
      ])
    );
  });

  it('should return the list of tables', () => {
    expect(getSuggestions({ ...props, tables: ['foo', 'bar'] })).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Table)',
          label: 'foo',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
        expect.objectContaining({
          detail: '(Table)',
          label: 'bar',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
      ])
    );
  });

  it('should return the list of measures', () => {
    expect(getSuggestions({ ...props, measures: ['foo', 'bar'] })).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Measure)',
          label: 'foo',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
        expect.objectContaining({
          detail: '(Measure)',
          label: 'bar',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
      ])
    );
  });

  it('should return the list of dimensions', () => {
    expect(getSuggestions({ ...props, dimensions: ['foo', 'bar'] })).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Dimension)',
          label: 'foo',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
        expect.objectContaining({
          detail: '(Dimension)',
          label: 'bar',
          kind: CodeEditorSuggestionItemKind.Property,
        }),
      ])
    );
  });
});
