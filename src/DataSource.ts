import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { TimestreamQuery, TimestreamOptions, MeasureInfo, DataType } from './types';

export class DataSource extends DataSourceWithBackend<TimestreamQuery, TimestreamOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<TimestreamOptions>) {
    super(instanceSettings);
  }

  getQueryDisplayText(query: TimestreamQuery): string {
    return query.rawQuery ?? '';
  }

  applyTemplateVariables(query: TimestreamQuery): TimestreamQuery {
    if (!query.rawQuery) {
      return query;
    }

    const local: ScopedVars = {};
    maybeAllVariable(local, 'database', query);
    maybeAllVariable(local, 'table', query);
    maybeAllVariable(local, 'measure', query);
    return {
      ...query,
      rawQuery: getTemplateSrv().replace(query.rawQuery, local),
    };
  }

  //----------------------------------------------
  // SCHEMA Style Functions
  //----------------------------------------------

  async getDatabases(like?: string): Promise<string[]> {
    return Promise.resolve(['db1', 'db2']);
  }

  async getTables(db: string, like?: string): Promise<string[]> {
    return Promise.resolve(['t1', 't2', 't3']);
  }

  async getMeasures(db: string, table: string): Promise<MeasureInfo[]> {
    return Promise.resolve([
      {
        name: 'm1',
        type: DataType.double,
        dimensions: { availability_zone: DataType.varchar },
      },
    ]);
  }
}

function maybeAllVariable(vars: ScopedVars, key: string, obj: any) {
  const value = obj[key];
  if (value && !value.startsWith('$')) {
    vars[key] = { text: key, value: value };
  }
}
