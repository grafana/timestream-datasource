import { DataSourceInstanceSettings, ScopedVars, DataQueryResponse, DataFrame, LoadingState } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { TimestreamQuery, TimestreamOptions, MeasureInfo, DataType, TimestreamCustomMeta } from './types';

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

  // In flight data..
  pending = new Map<string, DataFrame>();

  processResponse(res: DataQueryResponse): Promise<DataQueryResponse> {
    if (res.data?.length) {
      let data = res.data[0] as DataFrame;
      const meta = data.meta?.custom as TimestreamCustomMeta;
      if (meta && meta.queryId) {
        const old = this.pending.get(meta.queryId);
        if (old) {
          console.log('TODO, append', old, data);
          res.data = [old]; // new array
        }

        if (meta.nextToken) {
          res.state = LoadingState.Streaming; // Spinner
          if (meta.hasSeries) {
            // nothing special since each will get a unique key
          } else {
            res.key = meta.queryId; // We need to append the rows explicitly
            this.pending.set(meta.queryId, data);
          }
          alert(`TODO... query: ${meta.nextToken}`);
        } else {
          this.pending.delete(meta.queryId);
        }
      }
    }
    return Promise.resolve(res);
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
