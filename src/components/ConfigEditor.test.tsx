import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import { select } from 'react-select-event';

import { mockDatasourceOptions } from '../__mocks__/datasource';
import { ConfigEditor, applySavedDatasource } from './ConfigEditor';
import { selectors } from './selectors';

const resourceName = 'foo';
const mockPut = jest.fn();

jest.mock('@grafana/aws-sdk', () => {
  return {
    ...(jest.requireActual('@grafana/aws-sdk') as any),
    ConnectionConfig: function ConnectionConfig(props: any) {
      return <button onClick={() => props.onOptionsChange(props.options)}>Change connection</button>;
    },
  };
});
jest.mock('@grafana/runtime', () => {
  return {
    ...(jest.requireActual('@grafana/runtime') as any),
    getBackendSrv: () => ({
      put: mockPut,
      post: jest.fn().mockResolvedValue([resourceName]),
      get: jest.fn().mockResolvedValue([resourceName]),
    }),
  };
});
const props = mockDatasourceOptions;

type resourceType = 'defaultDatabase' | 'defaultTable' | 'defaultMeasure';

describe('ConfigEditor', () => {
  it('merges the server-minted external ID from the save response', () => {
    const current = mockDatasourceOptions.options;
    const saved = {
      ...current,
      version: (current.version ?? 1) + 1,
      jsonData: {
        ...current.jsonData,
        grafanaExternalId: '5285-uid-abcdef0123456789',
      },
    };

    const next = applySavedDatasource(current, saved);
    expect(next.version).toBe(saved.version);
    expect(next.jsonData.defaultRegion).toBe(current.jsonData.defaultRegion);
    expect(next.jsonData.grafanaExternalId).toBe('5285-uid-abcdef0123456789');
  });

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

  it('saves through the uid route before loading databases and passes the merged options up', async () => {
    mockPut.mockResolvedValueOnce({
      datasource: {
        ...props.options,
        version: 2,
        jsonData: { ...props.options.jsonData, grafanaExternalId: '5285-1-abcdef0123456789' },
      },
    });
    const onOptionsChange = jest.fn();
    render(<ConfigEditor {...props} onOptionsChange={onOptionsChange} />);

    fireEvent.click(screen.getByRole('button', { name: 'Change connection' }));
    expect(mockPut).not.toHaveBeenCalled();

    await select(screen.getByLabelText(selectors.components.ConfigEditor.defaultDatabase.input), resourceName, {
      container: document.body,
    });

    expect(mockPut).toHaveBeenCalledTimes(1);
    expect(mockPut).toHaveBeenCalledWith('/api/datasources/uid/1', props.options);
    await waitFor(() =>
      expect(onOptionsChange).toHaveBeenCalledWith(
        expect.objectContaining({
          version: 2,
          jsonData: expect.objectContaining({
            defaultRegion: 'us-east-2',
            grafanaExternalId: '5285-1-abcdef0123456789',
          }),
        })
      )
    );
  });
});
