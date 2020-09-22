import {
  DataSourceInstanceSettings,
  DataQueryResponse,
  DataFrame,
  DataQueryRequest,
  MetricFindValue,
  ScopedVars,
  QueryResultMetaStat,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { TimestreamQuery, TimestreamOptions, TimestreamCustomMeta, MeasureInfo, DataType } from './types';
import { getRequestLooper, MultiRequestTracker } from 'requestLooper';

export class DataSource extends DataSourceWithBackend<TimestreamQuery, TimestreamOptions> {
  // Easy access for QueryEditor
  options: TimestreamOptions;

  constructor(instanceSettings: DataSourceInstanceSettings<TimestreamOptions>) {
    super(instanceSettings);
    this.options = instanceSettings.jsonData;
  }

  // This will support annotation queries for 7.2+
  annotations = {};

  async metricFindQuery(query: any, options?: any): Promise<MetricFindValue[]> {
    if (!query) {
      return Promise.resolve([]);
    }
    const q = getTemplateSrv().replace(query as string);
    return this.getStrings(q)
      .toPromise()
      .then(strings => {
        return strings.map(s => ({
          text: s,
        }));
      });
  }

  /**
   * Do not execute queries that do not exist yet
   */
  filterQuery(query: TimestreamQuery): boolean {
    return !!query.rawQuery;
  }

  getQueryDisplayText(query: TimestreamQuery): string {
    return query.rawQuery ?? '';
  }

  applyTemplateVariables(query: TimestreamQuery, scopedVars: ScopedVars): TimestreamQuery {
    if (!query.rawQuery) {
      return query;
    }

    const templateSrv = getTemplateSrv();
    return {
      ...query,
      database: templateSrv.replace(query.database || '', scopedVars),
      table: templateSrv.replace(query.table || '', scopedVars),
      measure: templateSrv.replace(query.measure || '', scopedVars),
      rawQuery: templateSrv.replace(query.rawQuery), // DO NOT include scopedVars! it uses $__interval_ms!!!!!
    };
  }

  query(request: DataQueryRequest<TimestreamQuery>): Observable<DataQueryResponse> {
    let tracker: TimestreamCustomMeta | undefined = undefined;
    let queryId: string | undefined = undefined;
    return getRequestLooper(request, {
      // Check for a "nextToken" in the response
      getNextQuery: (rsp: DataQueryResponse) => {
        if (rsp.data?.length) {
          const first = rsp.data[0] as DataFrame;
          const meta = first.meta?.custom as TimestreamCustomMeta;
          if (meta && meta.nextToken) {
            queryId = meta.queryId;

            return {
              refId: first.refId,
              rawQuery: first.meta?.executedQueryString,
              nextToken: meta.nextToken,
            } as TimestreamQuery;
          }
        }
        return undefined;
      },

      /**
       * The original request
       */
      query: (request: DataQueryRequest<TimestreamQuery>) => {
        return super.query(request);
      },

      /**
       * Process the results
       */
      process: (t: MultiRequestTracker, data: DataFrame[], isLast: boolean) => {
        const meta = data[0]?.meta?.custom as TimestreamCustomMeta;
        if (!meta) {
          return data; // NOOP
        }
        // Single request
        meta.fetchStartTime = t.fetchStartTime;
        meta.fetchEndTime = t.fetchEndTime;
        meta.fetchTime = t.fetchEndTime! - t.fetchStartTime!;

        if (tracker) {
          // Additional request
          if (!tracker.subs?.length) {
            const { subs, nextToken, queryId, ...rest } = tracker;
            (rest as any).requestNumber = 1;
            tracker.subs?.push(rest as TimestreamCustomMeta);
          }
          for (const m of tracker.subs!) {
            delete m.nextToken; // not useful in the
          }
          delete (meta as any).queryId;
          (meta as any).requestNumber = tracker.subs!.length + 1;

          tracker.subs!.push(meta);
          tracker.fetchEndTime = t.fetchEndTime;
          tracker.fetchTime = t.fetchEndTime! - tracker.fetchStartTime!;
          tracker.executionFinishTime = meta.executionFinishTime;

          data[0].meta!.custom = tracker;
        } else {
          // First request
          tracker = {
            ...t,
            ...meta,
            subs: [],
          } as TimestreamCustomMeta;
        }

        // Calculate stats
        if (isLast && tracker.executionStartTime && tracker.executionFinishTime) {
          delete tracker.nextToken;

          const tsTime = tracker.executionFinishTime - tracker.executionStartTime;
          if (tsTime > 0) {
            const stats: QueryResultMetaStat[] = [];
            if (tracker.subs && tracker.subs.length) {
              stats.push({
                displayName: 'HTTP request count',
                value: tracker.subs.length,
                unit: 'none',
              });
            }
            stats.push({
              displayName: 'Execution time (Grafana server ⇆ Timestream)',
              value: tsTime,
              unit: 'ms',
              decimals: 2,
            });
            if (tracker.fetchStartTime) {
              tracker.fetchEndTime = Date.now();
              const dsTime = tracker.fetchEndTime - tracker.fetchStartTime;
              tracker.fetchTime = dsTime - tsTime;
              if (dsTime > tsTime) {
                stats.push({
                  displayName: 'Fetch time (Browser ⇆ Grafana server w/o Timestream)',
                  value: tracker.fetchTime,
                  unit: 'ms',
                  decimals: 2,
                });
                stats.push({
                  displayName: 'Fetch overhead',
                  value: (tracker.fetchTime / dsTime) * 100,
                  unit: 'percent', // 0 - 100
                });
              }
            }
            data[0].meta!.stats = stats;
          }
        }
        return data;
      },

      /**
       * Callback that gets executed when unsubscribed
       */
      onCancel: (tracker: MultiRequestTracker) => {
        if (queryId) {
          console.log('Cancelling running timestream query');

          // tracker.killed = true;
          this.postResource(`cancel`, {
            queryId,
          })
            .then(v => {
              console.log('Timestream query Canceled:', v);
            })
            .catch(err => {
              err.isHandled = true; // avoid the popup
              console.log('error killing', err);
            });
        }
      },
    });
  }

  //----------------------------------------------
  // SCHEMA Style Functions
  //----------------------------------------------

  private getStrings(rawQuery: string): Observable<string[]> {
    return this.query(({
      targets: [
        {
          refId: 'GetStrings',
          rawQuery,
        },
      ],
    } as unknown) as DataQueryRequest).pipe(
      map(res => {
        const first = res.data[0] as DataFrame;
        const vals = first.fields[0]?.values;
        if (vals) {
          return vals.toArray(); //
        }
        return [];
      })
    );
  }

  async getDatabases(like?: string): Promise<string[]> {
    return this.getStrings('SHOW DATABASES').toPromise();
  }

  async getTables(db: string): Promise<string[]> {
    if (!db) {
      return [];
    }
    return this.getStrings(`SHOW TABLES FROM ${db}`).toPromise();
  }

  async getMeasureInfo(db: string, table: string): Promise<MeasureInfo[]> {
    if (!db || !table) {
      return [];
    }
    return this.query(({
      targets: [
        {
          refId: 'X',
          rawQuery: `SHOW MEASURES FROM ${db}.${table}`,
        },
      ],
    } as unknown) as DataQueryRequest)
      .toPromise()
      .then(res => {
        const rsp: MeasureInfo[] = [];
        const first = res.data[0] as DataFrame;
        if (!first || !first.length) {
          return rsp;
        }
        const name = first.fields[0]?.values;
        const type = first.fields[1]?.values;
        const dims = first.fields[2]?.values;

        for (let i = 0; i < first.length; i++) {
          const dimensions = (JSON.parse(dims.get(i)) as any[]).map(row => {
            return row.dimension_name;
          });
          rsp.push({
            name: name.get(i) as string,
            type: type.get(i) as DataType,
            dimensions,
          });
        }
        return rsp;
      });
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
