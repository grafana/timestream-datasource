import { MutableDataFrame, createDataFrame } from '@grafana/data';

import { appendMatchingFrames } from './appendFrames';

describe('appendMatchingFrames', () => {
  it('should do nothing with empty input', () => {
    expect(appendMatchingFrames([], [])).toEqual([]);
  });

  it('should return the unchanged input when it is a MutableDataFrames array', () => {
    const v = [1];
    const frame = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    expect(appendMatchingFrames([frame], [])).toEqual([frame]);
  });

  it('should return a MutableDataFrame version of input when input is a DataFrames array', () => {
    const v = [1];
    const frame = createDataFrame({ fields: [{ name: 'foo', values: v }] });
    const expected = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    const result = appendMatchingFrames([frame], []);

    expect(JSON.stringify(result)).toEqual(JSON.stringify([expected]));
    expect(result[0]).toBeInstanceOf(MutableDataFrame);
  });

  it('should append new frames to an empty array', () => {
    const v = [1];
    const frame = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    expect(appendMatchingFrames([], [frame])).toEqual([frame]);
  });

  it('should append new frames to an existing array', () => {
    const v = [1];
    const frame1 = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    const frame2 = new MutableDataFrame({ fields: [{ name: 'bar', values: v }] });
    expect(appendMatchingFrames([frame1], [frame2])).toEqual([frame1, frame2]);
  });

  it('should merge two frames values', () => {
    const v = [1];
    const frame1 = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    const v2 = [2];
    const frame2 = new MutableDataFrame({ fields: [{ name: 'foo', values: v2 }] });
    const mergedFrames = appendMatchingFrames([frame1], [frame2]);
    expect(mergedFrames).toHaveLength(1);
    expect(mergedFrames[0].fields[0].name).toEqual('foo');
    expect(mergedFrames[0].fields[0].values).toEqual(expect.arrayContaining([1, 2]));
  });
});
