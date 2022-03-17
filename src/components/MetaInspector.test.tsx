import { render, screen } from '@testing-library/react';
import React from 'react';

import { mockDatasource } from '../__mocks__/datasource';
import { MetaInspector, Props } from './MetaInspector';

const props: Props = {
  datasource: mockDatasource,
  data: [],
};

describe('MetaInspector', () => {
  it('should return no data', async () => {
    render(<MetaInspector {...props} data={[]} />);

    const d = screen.getByText('No Data');
    expect(d).toBeInTheDocument();
  });

  it('should return metadata', async () => {
    render(
      <MetaInspector
        {...props}
        data={[
          {
            fields: [],
            length: 0,
            meta: {
              custom: {
                queryId: 'foo',
                nextToken: 'bar',
              },
            },
          },
        ]}
      />
    );

    expect(screen.getByText('Query ID')).toBeInTheDocument();
    expect(screen.getByText('foo')).toBeInTheDocument();
    expect(screen.getByText('Next Token')).toBeInTheDocument();
    expect(screen.getByText('bar')).toBeInTheDocument();
  });
});
