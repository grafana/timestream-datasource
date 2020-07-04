import { TimestreamQuery } from './types';
import { DataSource } from './DataSource';
import { SelectableValue } from '@grafana/data';

export class SchemaInfo {
  state: Partial<TimestreamQuery>;

  databases?: Array<SelectableValue<string>>;
  tables?: Array<SelectableValue<string>>;
  measures?: Array<SelectableValue<string>>;

  constructor(private ds: DataSource, q: Partial<TimestreamQuery>) {
    this.state = q || {};
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
    return (this.state = { ...this.state, ...state });
  }

  getDatabases = async (query?: string) => {
    if (this.databases) {
      return Promise.resolve(this.databases);
    }
    const vals = await this.ds.getDatabases();
    this.databases = vals.map(name => {
      return { label: name, value: name };
    });
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
      this.databases = vals.map(name => {
        return { label: name, value: name };
      });
      return this.databases;
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
      return this.filterMeasures(query);
    });
  };
}
