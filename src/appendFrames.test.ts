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

  it('should return input MutableDataFrames', () => {
    const frame = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    expect(appendMatchingFrames([frame], [])).toEqual([frame]);
  });

  it('should return input DataFrame', () => {
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

  it('should append new frames', () => {
    const frame = new MutableDataFrame({ fields: [{ name: 'foo', values: v }] });
    expect(appendMatchingFrames([], [frame])).toEqual([frame]);
  });
});
