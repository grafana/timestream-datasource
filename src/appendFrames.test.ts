import { DataFrame, FieldType, MutableDataFrame, Vector } from '@grafana/data';

import { appendMatchingFrames } from './appendFrames';

const v: Vector<number> = {
  length: 1,
  get: () => 1,
  toArray: () => [1],
};

describe('appendMatchingFrames', () => {
  it('should do nothing with empty input', () => {
    expect(appendMatchingFrames([], [])).toEqual([]);
  });

  it('should return the unchanged input when it is a MutableDataFrames array', () => {
    const frame = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    expect(appendMatchingFrames([frame], [])).toEqual([frame]);
  });

  it('should return a MutableDataFrame version of input when input is a DataFrames array', () => {
    const frame: DataFrame = {
      length: 1,
      fields: [
        {
          name: 'foo',
          type: FieldType.number,
          config: {},
          values: v,
        },
      ],
    };
    const expected = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    expect(appendMatchingFrames([frame], [])).toEqual([expected]);
  });

  it('should append new frames to an empty array', () => {
    const frame = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    expect(appendMatchingFrames([], [frame])).toEqual([frame]);
  });

  it('should append new frames to an existing array', () => {
    const frame1 = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    const frame2 = new MutableDataFrame({ fields: [{ name: 'bar', values: v }] });
    expect(appendMatchingFrames([frame1], [frame2])).toEqual([frame1, frame2]);
  });

  it('should merge two frames values', () => {
    const frame1 = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    const v2: Vector<number> = {
      length: 1,
      get: () => 2,
      toArray: () => [2],
    };
    const frame2 = new MutableDataFrame({ fields: [{ name: 'foo', values: v2 }] });
    const mergedFrames = appendMatchingFrames([frame1], [frame2]);
    expect(mergedFrames).toHaveLength(1);
    const mergedValue: Vector<number> = {
      length: 2,
      get: () => 1,
      toArray: () => [1, 2],
    };
    expect(mergedFrames[0]).toEqual(new MutableDataFrame({ fields: [{ name: 'foo', values: mergedValue }] }));
  });
});
