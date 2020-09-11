import {
  DataSourceInstanceSettings,
  DataQueryResponse,
  DataFrame,
  LoadingState,
  DataQueryRequest,
  MetricFindValue,
  ScopedVars,
  QueryResultMetaStat,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { Observable, Subscription } from 'rxjs';
import { map } from 'rxjs/operators';

import { TimestreamQuery, TimestreamOptions, TimestreamCustomMeta, MeasureInfo, DataType } from './types';

export class DataSource extends DataSourceWithBackend<TimestreamQuery, TimestreamOptions> {
  // Easy access for QueryEditor
  options: TimestreamOptions;

  constructor(instanceSettings: DataSourceInstanceSettings<TimestreamOptions>) {
    super(instanceSettings);
    this.options = instanceSettings.jsonData;
  }

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
    return new Observable<DataQueryResponse>(subscriber => {
      let nextRequest: DataQueryRequest<TimestreamQuery> | undefined = undefined;
      let subscription: Subscription | undefined = undefined;

      let count = 1;
      const tracker = {
        fetchStartTime: Date.now(),
      } as TimestreamCustomMeta;

      // Single observer gets reused
      const observer = {
        next: (rsp: DataQueryResponse) => {
          nextRequest = undefined;
          const meta = getNextTokenMeta(rsp);
          if (meta) {
            const nextQuery = this.updateLoadingFrameStatsAndGetNextQuery(rsp, tracker);
            nextRequest = {
              ...request,
              requestId: `${request.requestId}.${++count}`,
              startTime: Date.now(),
              targets: [nextQuery],
            };

            // The whole response is included in the results
            subscriber.next({ ...rsp, state: LoadingState.Loading });
          } else {
            tracker.fetchEndTime = Date.now();
            this.updateFinalFrameStats(rsp, tracker);
            let state = rsp.state;
            if (state !== LoadingState.Error) {
              state = LoadingState.Done;
            }
            subscriber.next({ ...rsp, state });
          }
        },
        error: (err: any) => {
          subscriber.error(err);
        },
        complete: () => {
          if (subscription) {
            subscription.unsubscribe();
            subscription = undefined;
          }

          // Let the previous request finish first
          if (nextRequest) {
            subscription = super.query(nextRequest).subscribe(observer);
            nextRequest = undefined;
          }
        },
      };

      // First request
      subscription = super.query(request).subscribe(observer);

      return () => {
        observer.complete();
        console.log('UNSUBSCRIBE', request.requestId);

        if (!tracker.fetchEndTime && tracker.queryId) {
          const elapsed = Date.now() - tracker.fetchStartTime!;
          console.log('Cancelling running timestream query', elapsed);

          // tracker.killed = true;
          this.postResource(`cancel`, {
            queryId: tracker.queryId,
          })
            .then(v => {
              console.log('Timestream query Canceled:', v);
            })
            .catch(err => {
              err.isHandled = true; // avoid the popup
              console.log('error killing', err);
            });
        }
      };
    });
  }

  private updateLoadingFrameStatsAndGetNextQuery = (response: DataQueryResponse, tracker: TimestreamCustomMeta) => {
    const frame = response.data[0] as DataFrame;
    const refId = frame.refId;
    const rawQuery = frame.meta?.executedQueryString;
    const meta = frame.meta?.custom as TimestreamCustomMeta;
    const nextToken = meta.nextToken!; // this stream exists because we know it is valid
    meta.fetchStartTime = tracker.fetchStartTime;
    meta.fetchEndTime = Date.now();
    meta.fetchTime = meta.fetchEndTime - tracker.fetchStartTime!;
    // Another request
    if (!tracker.subs) {
      tracker.queryId = meta.queryId;
      tracker.executionStartTime = meta.executionStartTime;
      tracker.subs = [];
    }
    tracker.subs.push(meta);
    frame.meta!.custom = tracker;

    // Cleanup the history a bit
    (meta as any).requestNumber = tracker.subs.length;
    delete meta.nextToken; // not useful in the response
    delete meta.queryId;

    return { refId, rawQuery, nextToken } as TimestreamQuery;
  };

  private updateFinalFrameStats = (response: DataQueryResponse, tracker: TimestreamCustomMeta) => {
    if (response.data?.length) {
      const first = response.data[0] as DataFrame;
      const custom = first.meta?.custom as TimestreamCustomMeta;
      custom.fetchStartTime = tracker.fetchStartTime;
      custom.fetchEndTime = Date.now();
      custom.fetchTime = custom.fetchEndTime - custom.fetchStartTime!;
      if (custom) {
        if (tracker.subs) {
          tracker.subs.push(custom);
          tracker.executionFinishTime = custom.executionFinishTime;
          tracker.hasSeries = custom.hasSeries;
        } else {
          tracker = {
            ...tracker,
            ...custom,
          };
        }
        delete tracker.nextToken;

        // Add stats
        if (tracker.executionStartTime && tracker.executionFinishTime) {
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
            first.meta!.stats = stats;
            first.meta!.custom = tracker;
          }
        }
      }
    }
  };

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
