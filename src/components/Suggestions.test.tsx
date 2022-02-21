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

describe('getSuggestions', () => {
  it('should include macros', () => {
    expect(getSuggestions([], [], [], [], '', '', '')).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          detail: '(Macro)',
          label: '$__timeFilter',
          kind: CodeEditorSuggestionItemKind.Method,
        }),
      ])
    );
  });

  it('should include template variables', () => {
    expect(getSuggestions([], [], [], [], '', '', '')).toEqual(
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
    expect(getSuggestions([], [], [], [], 'db', '', '')).toEqual(
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
    expect(getSuggestions([], [], [], [], '', 'tab', '')).toEqual(
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
    expect(getSuggestions([], [], [], [], '', '', 'cpu')).toEqual(
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
    expect(getSuggestions(['foo', 'bar'], [], [], [], '', '', '')).toEqual(
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
    expect(getSuggestions([], ['foo', 'bar'], [], [], '', '', '')).toEqual(
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
    expect(getSuggestions([], [], ['foo', 'bar'], [], '', '', '')).toEqual(
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
    expect(getSuggestions([], [], [], ['foo', 'bar'], '', '', '')).toEqual(
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
