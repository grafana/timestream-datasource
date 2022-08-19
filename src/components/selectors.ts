import { E2ESelectors } from '@grafana/e2e-selectors';

export const Components = {
  ConfigEditor: {
    AuthenticationProvider: {
      input: 'Authentication Provider',
    },
    SecretKey: {
      input: 'Secret Access Key',
    },
    AccessKey: {
      input: 'Access Key ID',
    },
    DefaultRegion: {
      input: 'Default Region',
    },
    defaultDatabase: {
      input: 'Database',
      wrapper: 'data-testid onloaddatabase',
    },
    defaultTable: {
      input: 'Table',
      wrapper: 'data-testid onloadtable',
    },
    defaultMeasure: {
      input: 'Measure',
      wrapper: 'data-testid onloadmeasure',
    },
  },
  QueryEditor: {
    CodeEditor: {
      container: 'Code editor container',
    },
  },
};

export const selectors: { components: E2ESelectors<typeof Components> } = {
  components: Components,
};
