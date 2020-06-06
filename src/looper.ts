import { TimestreamCustomMeta, TimestreamQuery } from 'types';
import { DataSource, getNextTokenMeta } from 'DataSource';
import { DataQueryResponse, DataFrame, DataQueryRequest, LoadingState } from '@grafana/data';
import { Subscriber } from 'rxjs';

interface LooperArgs {
  subscriber: Subscriber<DataQueryResponse>;
  ds: DataSource;
  count: number;
  req: DataQueryRequest<TimestreamQuery>;
  rsp: DataQueryResponse;
  frame?: DataFrame;
}

// NOTE: this assumes one query for now!!!!
export async function keepChecking(args: LooperArgs): Promise<boolean> {
  const first = args.rsp.data[0] as DataFrame;
  let meta = first.meta?.custom as TimestreamCustomMeta;
  if (!meta || !meta.nextToken) {
    args.subscriber.complete();
    return true;
  }
  if (!meta.hasSeries && !args.frame) {
    args.frame = first;
  }

  if (first.length < 2) {
    // Wait a little bit...
    await new Promise(r => setTimeout(r, args.count * 1000));
  }

  const r2 = {
    ...args.req,
    targets: [
      {
        refId: first.refId,
        rawQuery: first.meta?.executedQueryString, // :(
        nextToken: meta.nextToken!,
      } as TimestreamQuery,
    ],
  };

  const sub = await args.ds.query(r2).toPromise();
  const ttt = getNextTokenMeta(sub)!;
  const done = !(ttt && ttt.nextToken);
  if (sub.state !== LoadingState.Error) {
    sub.state = done ? LoadingState.Done : LoadingState.Loading;
  }
  if (args.frame) {
    sub.key = meta.queryId; // replace the data...
    const append = sub.data[0] as DataFrame;
    if (append) {
      if (args.frame.length) {
        console.log('TODO.... append', append, args.frame);
      } else {
        args.frame = append; //
      }
    }
  }
  args.subscriber.next(sub);
  if (done) {
    args.subscriber.complete();
    return true;
  }

  // Check again
  return await keepChecking({
    ...args,
    count: args.count + 1,
    rsp: sub,
  });
}
