import { e2e } from '@grafana/e2e';

import { selectors as timestreamSelectors } from '../../src/components/selectors';

const e2eSelectors = e2e.getSelectors(timestreamSelectors.components);

const query = 'SHOW DATABASES';
const queryVariable = 'query';

export const addDataSourceWithKey = (datasourceType: string, datasource: any): any => {
  return e2e.flows.addDataSource({
    checkHealth: false,
    expectedAlertMessage: 'Connection success',
    form: () => {
      e2eSelectors.ConfigEditor.AuthenticationProvider.input().type('Access & secret key').type('{enter}');
      e2eSelectors.ConfigEditor.AccessKey.input().type(datasource.secureJsonData.accessKey);
      e2eSelectors.ConfigEditor.SecretKey.input().type(datasource.secureJsonData.secretKey);
      e2eSelectors.ConfigEditor.DefaultRegion.input()
        .click({ force: true })
        .type(datasource.jsonData.defaultRegion)
        .type('{enter}');
      // Databases
      e2eSelectors.ConfigEditor.defaultDatabase.input().click({ force: true });
      // wait for it to load
      e2e()
        .get(`[data-testid="${timestreamSelectors.components.ConfigEditor.defaultDatabase.wrapper}"]`)
        .contains(datasource.jsonData.defaultDatabase);
      e2e()
        .get(`[data-testid="${timestreamSelectors.components.ConfigEditor.defaultDatabase.wrapper}"]`)
        .contains(datasource.jsonData.defaultDatabase);
      e2eSelectors.ConfigEditor.defaultDatabase.input().type(datasource.jsonData.defaultDatabase).type('{enter}');
      // Tables
      e2eSelectors.ConfigEditor.defaultTable.input().click({ force: true });
      // wait for it to load
      e2e()
        .get(`[data-testid="${timestreamSelectors.components.ConfigEditor.defaultTable.wrapper}"]`)
        .contains(datasource.jsonData.defaultTable);
      e2eSelectors.ConfigEditor.defaultTable.input().type(datasource.jsonData.defaultTable).type('{enter}');
      // Measures
      e2eSelectors.ConfigEditor.defaultMeasure.input().click({ force: true });
      // wait for it to load
      e2e()
        .get(`[data-testid="${timestreamSelectors.components.ConfigEditor.defaultMeasure.wrapper}"]`)
        .contains(datasource.jsonData.defaultMeasure);
      e2eSelectors.ConfigEditor.defaultMeasure.input().type(datasource.jsonData.defaultMeasure).type('{enter}');
    },
    type: datasourceType,
  });
};

const addTablePanel = (q: string) => {
  const fillQuery = (query: string) => {
    // Wait for the selectors to load
    e2e()
      .get(`[data-testid="${timestreamSelectors.components.ConfigEditor.defaultMeasure.wrapper}"]`)
      .contains('cpu_hi');
    e2eSelectors.QueryEditor.CodeEditor.container().type(query);
  };

  e2e.flows.addPanel({
    matchScreenshot: true,
    visualizationName: e2e.flows.VISUALIZATION_TABLE,
    queriesForm: () => {
      fillQuery(q);
      // Blur the editor to execute the query and wait
      cy.get('.panel-content').last().click();
      cy.get('.panel-loading');
      cy.get('.panel-loading', { timeout: 10000 }).should('not.exist');
    },
  });

  e2e.flows.explore({
    matchScreenshot: false,
    timeRange: {
      from: '2001-01-31 19:00:00',
      to: '2016-01-31 19:00:00',
    },
    queriesForm: () => fillQuery(query),
  });
};

e2e.scenario({
  describeName: 'Smoke tests',
  itName: 'Login, create data source, dashboard with variable and panel',
  scenario: () => {
    e2e()
      .readProvisions(['datasources/aws-timestream.yaml'])
      .then(([provision]) => {
        const datasource = provision.datasources[0];
        return addDataSourceWithKey('Amazon Timestream', datasource);
      })
      .then(() => {
        e2e.flows.addDashboard({
          timeRange: {
            from: '2001-01-31 19:00:00',
            to: '2016-01-31 19:00:00',
          },
          variables: [
            {
              constantValue: query,
              label: 'Template Variable',
              name: queryVariable,
              type: e2e.flows.VARIABLE_TYPE_CONSTANT,
            },
          ],
        });
        addTablePanel('$' + queryVariable);
      });
  },
});
