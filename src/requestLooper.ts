import { DataQuery, DataQueryRequest, DataQueryResponse, LoadingState, DataFrame } from '@grafana/data';
import { Observable, Subscription } from 'rxjs';

export interface MultiRequestTracker {
  fetchStartTime?: number; // The frontend clock
  fetchEndTime?: number; // The frontend clock
}

export interface RequestLoopOptions<TQuery extends DataQuery = DataQuery> {
  /**
   * If the response needs an additional request to execute, return it here
   */
  getNextQuery: (rsp: DataQueryResponse) => TQuery | undefined;

  /**
   * The datasource execute method
   */
  query: (req: DataQueryRequest<TQuery>) => Observable<DataQueryResponse>;

  /**
   * Process the results
   */
  process: (tracker: MultiRequestTracker, data: DataFrame[], isLast: boolean) => DataFrame[];

  /**
   * Callback that gets executed when unsubscribed
   */
  onCancel: (tracker: MultiRequestTracker) => void;
}

/**
 * Continue executing requests as long as `getNextQuery` returns a query
 */
export function getRequestLooper<T extends DataQuery = DataQuery>(
  req: DataQueryRequest<T>,
  options: RequestLoopOptions<T>
): Observable<DataQueryResponse> {
  return new Observable<DataQueryResponse>(subscriber => {
    let nextQuery: T | undefined = undefined;
    let subscription: Subscription | undefined = undefined;
    const tracker: MultiRequestTracker = {
      fetchStartTime: Date.now(),
      fetchEndTime: undefined,
    };
    let loadingState: LoadingState | undefined = LoadingState.Loading;
    let count = 1;

    // Single observer gets reused for each request
    const observer = {
      next: (rsp: DataQueryResponse) => {
        tracker.fetchEndTime = Date.now();
        loadingState = rsp.state;
        if (loadingState !== LoadingState.Error) {
          nextQuery = options.getNextQuery(rsp);
          loadingState = nextQuery ? LoadingState.Loading : LoadingState.Done;
        }

        const data = options.process(tracker, rsp.data, !!!nextQuery);
        subscriber.next({ ...rsp, data, state: loadingState });
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
        if (nextQuery) {
          tracker.fetchEndTime = undefined;
          tracker.fetchStartTime = Date.now();
          subscription = options
            .query({
              ...req,
              requestId: `${req.requestId}.${++count}`,
              startTime: tracker.fetchStartTime,
              targets: [nextQuery],
            })
            .subscribe(observer);
          nextQuery = undefined;
        } else {
          subscriber.complete();
        }
      },
    };

    // First request
    subscription = options.query(req).subscribe(observer);

    return () => {
      nextQuery = undefined;
      observer.complete();
      if (!tracker.fetchEndTime) {
        options.onCancel(tracker);
      }
    };
  });
}
