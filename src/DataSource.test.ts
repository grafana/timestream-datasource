import { ScopedVars } from '@grafana/data';
import * as runtime from '@grafana/runtime';

import { mockDatasource, mockQuery } from './__mocks__/datasource';

describe('DataSource', () => {
  describe('applyTemplateVariables', () => {
    const scopedVars: ScopedVars = {
      $simple: { value: 'foo' },
      $multiple: { value: ['foo', 'bar'] },
      __interval_ms: { value: 5000000 },
      __interval: { value: 50000 },
    };
    // simplified version of getTemplateSrv().replace
    const replaceMock = jest
      .fn()
      .mockImplementation((target?: string, scopedVars?: ScopedVars, format?: string | Function) => {
        let res = target ?? '';
        if (scopedVars && typeof format === 'function') {
          Object.keys(scopedVars).forEach((v) => (res = res.replace(v, format(scopedVars[v]?.value))));
        }
        return res;
      });
    beforeEach(() => {
      jest.spyOn(runtime, 'getTemplateSrv').mockImplementation(() => ({
        getVariables: jest.fn(),
        replace: replaceMock,
        containsTemplate: jest.fn(),
        updateTimeRange: jest.fn(),
      }));
    });

    it('should replace a simple var', () => {
      const res = mockDatasource.applyTemplateVariables(
        { ...mockQuery, rawQuery: 'select * from $simple' },
        scopedVars
      );
      expect(res.rawQuery).toEqual('select * from foo');
    });

    it('should replace a multiple var', () => {
      const res = mockDatasource.applyTemplateVariables(
        { ...mockQuery, rawQuery: 'select * from foo where var in ($multiple)' },
        scopedVars
      );
      expect(res.rawQuery).toEqual(`select * from foo where var in ('foo','bar')`);
    });

    it('should replace __interval interpolated variables with their original string', () => {
      replaceMock.mockClear();
      mockDatasource.applyTemplateVariables(
        { ...mockQuery, rawQuery: 'select $__interval_ms, $__interval' },
        {
          __interval_ms: { value: 5000000 },
          __interval: { value: 50000 },
        }
      );
      // check rawQuery.replace is called with correct interval value
      expect(replaceMock.mock.calls[3][1].__interval).toEqual({ value: '$__interval' });
      expect(replaceMock.mock.calls[3][1].__interval_ms).toEqual({ value: '$__interval_ms' });
    });

    it('should return number variables', () => {
      replaceMock.mockClear();
      mockDatasource.applyTemplateVariables(
        { ...mockQuery, rawQuery: 'select $__from' },
        {
          __from: { value: 3000 },
        }
      );
      // check rawQuery.replace is called with correct interval value
      expect(replaceMock.mock.calls[3][1].__from).toEqual({ value: 3000 });
    });
  });
});
