import { TimestreamQuery } from './types';
import { DataSource } from './DataSource';
import { SelectableValue, KeyValue } from '@grafana/data';
import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';
import { TemplateSrv } from '@grafana/runtime';

export class SchemaInfo {
  state: Partial<TimestreamQuery>;

  databases?: Array<SelectableValue<string>>;
  tables?: Array<SelectableValue<string>>;
  measures?: Array<SelectableValue<string>>;
  dimensions?: KeyValue<string[]>;

  constructor(private ds: DataSource, q: Partial<TimestreamQuery>, private templateSrv?: TemplateSrv) {
    this.state = { ...q };

    if (this.templateSrv) {
      this.updateState(q);
    }
  }

  updateState(state: Partial<TimestreamQuery>): Partial<TimestreamQuery> {
    if (state.database) {
      this.databases = undefined;
      this.tables = undefined;
      this.measures = undefined;
      this.dimensions = undefined;
    } else if (state.table) {
      this.measures = undefined;
      this.dimensions = undefined;
    } else if (state.measure) {
      this.dimensions = undefined;
    }

    const merged = { ...this.state, ...state };
    if (this.templateSrv) {
      const { defaultDatabase, defaultTable, defaultMeasure } = this.ds.options;
      if (!merged.database && defaultDatabase) {
        merged.database = defaultDatabase;
      }
      if (!merged.table && defaultTable) {
        merged.table = defaultTable;
      }
      if (!merged.measure && defaultMeasure) {
        merged.measure = defaultMeasure;
      }

      if (merged.database) {
        merged.database = this.templateSrv.replace(merged.database);
      }
      if (merged.table) {
        merged.table = this.templateSrv.replace(merged.table);
      }
      if (merged.measure) {
        merged.measure = this.templateSrv.replace(merged.measure);
      }
    }
    return (this.state = merged);
  }

  async preload() {
    await this.getDatabases();
    await this.getTables();
    return this.getMeasures();
  }

  getSuggestions = (): CodeEditorSuggestionItem[] => {
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
        label: '$__database',
        kind: CodeEditorSuggestionItemKind.Method,
        detail: `(Macro) ${this.state.database}`,
      },
      {
        label: '$__table',
        kind: CodeEditorSuggestionItemKind.Method,
        detail: `(Macro) ${this.state.table}`,
      },
      {
        label: '$__measure',
        kind: CodeEditorSuggestionItemKind.Method,
        detail: `(Macro) ${this.state.measure}`,
      },
    ];

    if (this.templateSrv) {
      this.templateSrv.getVariables().forEach(variable => {
        const label = '${' + variable.name + '}';
        let val = this.templateSrv!.replace(label);
        if (val === label) {
          val = '';
        }
        sugs.push({
          label,
          kind: CodeEditorSuggestionItemKind.Text,
          //origin: VariableOrigin.Template,
          detail: `(Template Variable) ${val}`,
        });
      });
    }

    if (this.databases) {
      for (const v of this.databases) {
        const label = getValidSuggestion(v);
        if (label) {
          sugs.push({
            label,
            kind: CodeEditorSuggestionItemKind.Property,
            detail: `(Database)`,
          });
        }
      }
    }

    if (this.tables) {
      for (const v of this.tables) {
        const label = getValidSuggestion(v);
        if (label) {
          sugs.push({
            label,
            kind: CodeEditorSuggestionItemKind.Property,
            detail: `(Table)`,
          });
        }
      }
    }

    if (this.measures) {
      for (const v of this.measures) {
        const label = getValidSuggestion(v);
        if (label) {
          sugs.push({
            label: `'${label}'`, // measure names are quoted
            kind: CodeEditorSuggestionItemKind.Property,
            detail: `(Measure)`,
          });
        }
      }
    }

    if (this.dimensions && this.state.measure) {
      const dims = this.dimensions[this.state.measure];
      if (dims) {
        for (const v of dims) {
          sugs.push({
            label: `'${v}'`, // measure names are quoted
            kind: CodeEditorSuggestionItemKind.Property,
            detail: `(Dimension)`,
          });
        }
      }
    }
    return sugs;
  };

  andTemplates(): Array<SelectableValue<string>> {
    if (this.templateSrv) {
      return this.templateSrv.getVariables().map(v => {
        const template = '${' + v.name + '}';
        return { label: template, value: template };
      });
    }
    return [];
  }

  getDatabases = async (query?: string) => {
    if (this.databases) {
      return Promise.resolve(this.databases);
    }
    const vals = await this.ds.getDatabases();
    this.databases = vals.map(name => {
      return { label: name, value: name };
    });
    if (this.templateSrv) {
      this.databases.push(...this.andTemplates());
      this.databases.push({
        label: '-- remove --',
        value: '',
      });
    }
    return this.databases;
  };

  getTables = async (query?: string) => {
    if (this.tables) {
      return Promise.resolve(this.tables);
    }
    if (!this.state.database) {
      return Promise.resolve([{ label: 'database not configured', value: '' }]);
    }
    return this.ds.getTables(this.state.database).then(vals => {
      this.tables = vals.map(name => {
        return { label: name, value: name };
      });
      if (this.templateSrv) {
        this.tables.push(...this.andTemplates());
        this.tables.push({
          label: '-- remove --',
          value: '',
        });
      }
      return this.tables;
    });
  };

  filterMeasures = async (query?: string) => {
    if (!this.measures) {
      return [];
    }
    if (!query) {
      return this.measures;
    }
    return this.measures.filter(f => {
      if (!f.value) {
        return;
      }
      return f.value.indexOf(query) >= 0;
    });
  };

  getMeasures = async (query?: string) => {
    if (this.measures) {
      return this.filterMeasures(query);
    }
    const { database, table } = this.state;
    if (!database) {
      return Promise.resolve([{ label: 'database not configured', value: '' }]);
    }
    if (!table) {
      return Promise.resolve([{ label: 'table not configured', value: '' }]);
    }
    return this.ds.getMeasureInfo(database, table).then(info => {
      const dims: KeyValue<string[]> = {};
      this.measures = info.map(v => {
        dims[v.name] = v.dimensions;
        return { label: `${v.name} (${v.type})`, value: v.name };
      });
      this.dimensions = dims;
      if (this.templateSrv) {
        this.measures.push(...this.andTemplates());
        this.measures.push({
          label: '-- remove --',
          value: '',
        });
      }
      return this.filterMeasures(query);
    });
  };
}

function getValidSuggestion(v: SelectableValue<string>): string | undefined {
  if (!v || !v.value) {
    return undefined;
  }
  const txt = v.value;
  if (txt.startsWith('$') || txt.startsWith('-')) {
    return undefined;
  }
  return txt;
}
