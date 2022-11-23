import '@testing-library/jest-dom';

import * as runtime from '@grafana/runtime';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import { select } from 'react-select-event';
import * as experimental from '@grafana/experimental';

import { mockDatasource, mockQuery } from '../__mocks__/datasource';
import { QueryEditor } from './QueryEditor';
import { sampleQueries } from './samples';
import { selectors } from './selectors';

jest
  .spyOn(runtime, 'getTemplateSrv')
  .mockImplementation(() => ({
    containsTemplate: jest.fn(),
    getVariables: jest.fn().mockReturnValue([]),
    replace: jest.fn(),
    updateTimeRange: jest.fn()})
  )

jest.mock('@grafana/experimental', () => ({
  ...jest.requireActual<typeof experimental>('@grafana/experimental'),
  SQLEditor: function SQLEditor() {
    return <></>;
  },
}));

const ds = mockDatasource;
const q = mockQuery;

const databases = ['db1', 'db2'];
const tables = ['t1', 't2'];
const measures = ['cpu', 'mem'];
const dimensions = ['region', 'zone'];

beforeEach(() => {
  ds.getResource = jest.fn().mockResolvedValue(databases);
  ds.postResource = jest.fn((r: string, body?: any) => {
    switch (r) {
      case 'tables':
        return Promise.resolve(tables);
      case 'measures':
        return Promise.resolve(measures);
      case 'dimensions':
        return Promise.resolve(dimensions);
    }
    return Promise.resolve([]);
  });
});

const props = {
  datasource: ds,
  query: q,
  onChange: jest.fn(),
  onRunQuery: jest.fn(),
};

describe('QueryEditor', () => {
  it('should request databases and execute the query', async () => {
    const onChange = jest.fn();
    render(<QueryEditor {...props} onChange={onChange} query={{ ...props.query, database: databases[0] }} />);

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.defaultDatabase.input);
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, databases[1], { container: document.body });

    expect(ds.getResource).toHaveBeenCalledWith('databases');
    expect(onChange).toHaveBeenCalledWith({
      ...q,
      database: databases[1],
    });
  });

  it('should request tables and execute the query', async () => {
    const onChange = jest.fn();
    render(<QueryEditor {...props} onChange={onChange} query={{ ...props.query, database: databases[0] }} />);

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.defaultTable.input);
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, tables[1], { container: document.body });

    expect(ds.postResource).toHaveBeenCalledWith('tables', { database: databases[0] });
    expect(onChange).toHaveBeenCalledWith({
      ...q,
      database: databases[0],
      table: tables[1],
    });
  });

  it('should request measures and execute the query', async () => {
    const onChange = jest.fn();
    render(
      <QueryEditor
        {...props}
        onChange={onChange}
        query={{ ...props.query, database: databases[0], table: tables[0] }}
      />
    );

    const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.defaultMeasure.input);
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, measures[1], { container: document.body });

    expect(ds.postResource).toHaveBeenCalledWith('measures', { database: databases[0], table: tables[0] });
    expect(onChange).toHaveBeenCalledWith({
      ...q,
      database: databases[0],
      table: tables[0],
      measure: measures[1],
    });
  });

  it('run a query if it has database, table and measure set', async () => {
    const onRunQuery = jest.fn();
    render(
      <QueryEditor
        {...props}
        onRunQuery={onRunQuery}
        query={{ ...props.query, database: databases[0], table: tables[0], measure: measures[0] }}
      />
    );

    await waitFor(() => expect(ds.postResource).toHaveBeenCalledTimes(2));
    // Measure field is set
    expect(screen.getByText(measures[0])).toBeInTheDocument();

    expect(onRunQuery).toHaveBeenCalled();
  });

  it('should enable switch to wait for all queries', async () => {
    const onChange = jest.fn();
    render(<QueryEditor {...props} onChange={onChange} />);
    await waitFor(() => expect(ds.getResource).toHaveBeenCalledTimes(1));

    const toggleButton = screen.getByLabelText('Wait for all queries');
    expect(toggleButton).toBeInTheDocument();

    fireEvent.click(toggleButton);
    expect(onChange).toHaveBeenCalledWith({
      ...q,
      waitForResult: true,
    });
  });

  it('should set the code of a sample', async () => {
    const onChange = jest.fn();
    render(<QueryEditor {...props} onChange={onChange} />);

    const selectEl = screen.getByLabelText('Query');
    expect(selectEl).toBeInTheDocument();

    await select(selectEl, sampleQueries[0].label!, { container: document.body });

    expect(onChange).toHaveBeenCalledWith({
      ...q,
      rawQuery: sampleQueries[0].value,
    });
  });
});
