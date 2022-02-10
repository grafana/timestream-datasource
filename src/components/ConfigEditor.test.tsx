import { render, screen } from '@testing-library/react';
import React from 'react';
import { select } from 'react-select-event';

import { mockDatasourceOptions } from '../__mocks__/datasource';
import { ConfigEditor } from './ConfigEditor';
import { selectors } from './selectors';

const resourceName = 'foo';

jest.mock('@grafana/aws-sdk', () => {
  return {
    ...(jest.requireActual('@grafana/aws-sdk') as any),
    ConnectionConfig: function ConnectionConfig() {
      return <></>;
    },
  };
});
jest.mock('@grafana/runtime', () => {
  return {
    ...(jest.requireActual('@grafana/runtime') as any),
    getBackendSrv: () => ({
      put: jest.fn().mockResolvedValue({ datasource: {} }),
      post: jest.fn().mockResolvedValue([resourceName]),
      get: jest.fn().mockResolvedValue([resourceName]),
    }),
  };
});
const props = mockDatasourceOptions;

type resourceType = 'defaultDatabase' | 'defaultTable' | 'defaultMeasure';

describe('ConfigEditor', () => {
  const types: resourceType[] = ['defaultDatabase', 'defaultTable', 'defaultMeasure'];
  types.forEach((resource) => {
    it(`should save and request ${resource}s`, async () => {
      const onChange = jest.fn();
      render(<ConfigEditor {...props} onOptionsChange={onChange} />);

      const selectEl = screen.getByLabelText(selectors.components.ConfigEditor[resource].input);
      expect(selectEl).toBeInTheDocument();

      await select(selectEl, resourceName, { container: document.body });

      expect(onChange).toHaveBeenCalledWith({
        ...props.options,
        jsonData: { ...props.options.jsonData, [resource]: resourceName },
      });
    });
  });
});
