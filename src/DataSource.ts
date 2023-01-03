import {
  DataFrame,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  getValueFormat,
  MetricFindValue,
  QueryResultMetaStat,
  ScopedVars,
  TimeRange,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { appendMatchingFrames } from 'appendFrames';
import { getRequestLooper, MultiRequestTracker } from 'requestLooper';
import { merge, Observable, of } from 'rxjs';
import { map } from 'rxjs/operators';

import { TimestreamCustomMeta, TimestreamOptions, TimestreamQuery } from './types';

let requestCounter = 100;
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
    return this.getStrings(q, options.range)
      .toPromise()
      .then((strings) => {
        return (strings || []).map((s) => ({
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

  private interpolateVariable = (value: string | string[]) => {
    if (typeof value === 'string') {
      return value;
    }

    const quotedValues = value.map((v) => {
      return this.quoteLiteral(v);
    });
    return quotedValues.join(',');
  };

  private quoteLiteral(value: string) {
    return "'" + value.replace(/'/g, "''") + "'";
  }

  applyTemplateVariables(query: TimestreamQuery, scopedVars: ScopedVars): TimestreamQuery {
    if (!query.rawQuery) {
      return query;
    }

    // create a copy of scopedVars without $__interval_ms for using with rawQuery
    // ${__interval*} should be escaped by the server, not the frontend
    const queryScopedVars = { ...scopedVars };
    delete queryScopedVars.__interval_ms;
    delete queryScopedVars.__interval;

    const templateSrv = getTemplateSrv();
    return {
      ...query,
      database: templateSrv.replace(query.database || '', scopedVars),
      table: templateSrv.replace(query.table || '', scopedVars),
      measure: templateSrv.replace(query.measure || '', scopedVars),
      rawQuery: templateSrv.replace(query.rawQuery, queryScopedVars, this.interpolateVariable),
    };
  }

  query(request: DataQueryRequest<TimestreamQuery>): Observable<DataQueryResponse> {
    const targets = request.targets;
    if (!targets.length) {
      return of({ data: [] });
    }
    if (targets.some((t) => t.waitForResult)) {
      // Defaults to the common behavior of waiting for all the queries to be finished
      // before rendering
      return super.query(request);
    }
    const all: Array<Observable<DataQueryResponse>> = [];
    for (let target of targets) {
      if (target.hide) {
        continue;
      }
      all.push(this.doSingle(target, request));
    }
    if (all.length === 1) {
      return all[0];
    }
    return merge(...all);
  }

  doSingle(target: TimestreamQuery, request: DataQueryRequest<TimestreamQuery>): Observable<DataQueryResponse> {
    let tracker: TimestreamCustomMeta | undefined = undefined;
    let queryId: string | undefined = undefined;
    let allData: DataFrame[] = [];
    return getRequestLooper(
      { ...request, targets: [target], requestId: `aws_ts_${requestCounter++}` },
      {
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
            return allData.length ? allData : data; // NOOP
          }
          const valueFormat = getValueFormat('decbytes');

          // Single request
          meta.fetchStartTime = t.fetchStartTime;
          meta.fetchEndTime = t.fetchEndTime;
          meta.fetchTime = t.fetchEndTime! - t.fetchStartTime!;

          if (meta.hasSeries || !allData.length) {
            for (const frame of data) {
              if (frame.fields.length > 0) {
                allData.push(frame);
              }
            }
          } else {
            if (data.length > 1) {
              console.log('non timeseries should have a single frame', data);
            }
            const append = data[0];
            if (append.length > 0) {
              allData = appendMatchingFrames(allData, data);
            }
          }

          // Empty results
          if (!allData[0]?.meta) {
            return data;
          }

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

            allData[0].meta!.custom = tracker;
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
              if (tracker.status.CumulativeBytesMetered) {
                const v = valueFormat(tracker.status.CumulativeBytesMetered, 2, 1024);
                stats.push({
                  displayName: 'Cumulative bytes metered',
                  value: Number(v.text),
                  unit: v.suffix?.trimLeft(),
                  decimals: 2,
                });
              }
              if (tracker.status.CumulativeBytesScanned) {
                const v = valueFormat(tracker.status.CumulativeBytesScanned, 2, 1024);
                stats.push({
                  displayName: 'Cumulative bytes scanned',
                  value: Number(v.text),
                  unit: v.suffix?.trimLeft(),
                  decimals: 2,
                });
              }
              allData[0].meta!.stats = stats;
            }
          }
          return allData;
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
              .then((v) => {
                console.log('Timestream query Canceled:', v);
              })
              .catch((err) => {
                err.isHandled = true; // avoid the popup
                console.log('error killing', err);
              });
          }
        },
      }
    );
  }

  //----------------------------------------------
  // SCHEMA Style Functions
  //----------------------------------------------

  private getStrings(rawQuery: string, range?: TimeRange): Observable<string[]> {
    return this.query({
      targets: [
        {
          refId: 'GetStrings',
          rawQuery,
        },
      ],
      range,
    } as unknown as DataQueryRequest).pipe(
      map((res) => {
        if (res.error) {
          const message = res.error.message ?? res.error.data?.message ?? 'Error getting variable';
          throw new Error(message);
        }
        const first = res.data[0] as DataFrame;
        if (!first || !first.length) {
          return [];
        }
        const vals = first.fields[0]?.values;
        if (vals) {
          return vals.toArray(); //
        }
        return [];
      })
    );
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
