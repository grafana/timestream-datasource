import {
  DataSourceInstanceSettings,
  ScopedVars,
  DataQueryResponse,
  DataFrame,
  LoadingState,
  DataQueryRequest,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { Observable } from 'rxjs';

import { TimestreamQuery, TimestreamOptions, TimestreamCustomMeta } from './types';
import { keepChecking } from 'looper';

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
    maybeSetVariable(local, 'database', query);
    maybeSetVariable(local, 'table', query);
    maybeSetVariable(local, 'measure', query);
    return {
      ...query,
      rawQuery: getTemplateSrv().replace(query.rawQuery, local),
    };
  }

  query(request: DataQueryRequest<TimestreamQuery>): Observable<DataQueryResponse> {
    return new Observable<DataQueryResponse>(subscriber => {
      super
        .query(request)
        .toPromise()
        .then(rsp => {
          const meta = getNextTokenMeta(rsp);
          if (meta) {
            rsp.state = LoadingState.Loading; // streaming?
            if (!meta.hasSeries) {
              rsp.key = meta.queryId;
            }
            keepChecking({
              subscriber,
              req: request,
              rsp,
              count: 1,
              ds: this,
            });
          }
          subscriber.next(rsp);
          if (!meta) {
            subscriber.complete(); // done
          }
        });

      return () => {
        console.log('unsubscribe.. timestream cancel?', this);
      };
    });
  }

  //----------------------------------------------
  // SCHEMA Style Functions
  //----------------------------------------------
  private async getStrings(rawQuery: string): Promise<string[]> {
    return this.query(({
      targets: [
        {
          refId: 'X',
          rawQuery,
        },
      ],
    } as unknown) as DataQueryRequest)
      .toPromise()
      .then(res => {
        const first = res.data[0] as DataFrame;
        const vals = first.fields[0]?.values;
        if (vals) {
          return vals.toArray(); //
        }
        return [];
      });
  }

  async getDatabases(like?: string): Promise<string[]> {
    return this.getStrings('SHOW DATABASES');
  }

  async getTables(db: string, like?: string): Promise<string[]> {
    if (!db) {
      return [];
    }
    return this.getStrings(`SHOW TABLES FROM ${db}`);
  }

  async getMeasures(db: string, table: string): Promise<string[]> {
    if (!db || !table) {
      return [];
    }
    return this.getStrings(`SHOW MEASURES FROM ${db}.${table}`);
  }
}

function maybeSetVariable(vars: ScopedVars, key: string, obj: any) {
  const value = obj[key];
  if (value && !value.startsWith('$')) {
    vars[key] = { text: key, value: value };
  }
}

export function getNextTokenMeta(rsp: DataQueryResponse): TimestreamCustomMeta | undefined {
  if (rsp.data?.length) {
    const first = rsp.data[0] as DataFrame;
    const meta = first.meta?.custom as TimestreamCustomMeta;
    if (meta && meta.nextToken) {
      return meta;
    }
  }
  return undefined;
}
