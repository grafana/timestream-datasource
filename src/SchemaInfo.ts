import { TimestreamQuery } from './types';
import { DataSource } from './DataSource';
import { SelectableValue } from '@grafana/data';
import { CodeEditorSuggestionItem, CodeEditorSuggestionItemKind } from '@grafana/ui';
import { TemplateSrv } from '@grafana/runtime';

export class SchemaInfo {
  state: Partial<TimestreamQuery>;

  databases?: Array<SelectableValue<string>>;
  tables?: Array<SelectableValue<string>>;
  measures?: Array<SelectableValue<string>>;

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
    } else if (state.table) {
      this.measures = undefined;
    } else if (state.measure) {
      // nothing right now?
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
    const always = [
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
    ];
    if (!this.templateSrv) {
      return always;
    }

    return [
      ...this.templateSrv.getVariables().map(variable => {
        const label = '${' + variable.name + '}';
        let val = this.templateSrv!.replace(label);
        if (val === label) {
          val = '';
        }
        return {
          label,
          kind: CodeEditorSuggestionItemKind.Text,
          //origin: VariableOrigin.Template,
          detail: `(Template Variable) ${val}`,
        };
      }),
      ...always,
    ];
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
      this.measures = info.map(v => {
        return { label: `${v.name} (${v.type})`, value: v.name };
      });
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
