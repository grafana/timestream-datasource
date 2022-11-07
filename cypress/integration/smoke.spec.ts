import { e2e } from '@grafana/e2e';

import { selectors as timestreamSelectors } from '../../src/components/selectors';

const e2eSelectors = e2e.getSelectors(timestreamSelectors.components);

const queryVariable = 'query';

export const addDataSourceWithKey = (datasourceType: string, datasource: any): any => {
  return e2e.flows.addDataSource({
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
    e2eSelectors.QueryEditor.CodeEditor.container().type('{selectall} SHOW DATABASES');
  };

  e2e.flows.addPanel({
    matchScreenshot: false,
    queriesForm: () => {
      // The following section will verify that autocompletion is behaving as expected.
      // Throughout the composition of the SQL query, the autocompletion engine will provide appropriate suggestions.
      // In this test the first few suggestions are accepted by hitting enter which will create a basic query.
      // Increasing delay to allow tables names and columns names to be resolved async by the plugin
      e2eSelectors.QueryEditor.CodeEditor.container()
        .click({ force: true })
        .type(`s{enter}{enter}{enter}g{enter}d{enter}{enter}c{enter}`, { delay: 5000 });
      e2eSelectors.QueryEditor.CodeEditor.container().contains('SELECT * FROM "grafanaDB"."DevOps" GROUP BY cell');

      fillQuery(q);
      // Blur the editor to execute the query and wait
      cy.get('.panel-content').last().click();
      cy.get('.panel-loading');
      cy.get('.panel-loading', { timeout: 10000 }).should('not.exist');
      cy.contains('Data does not have a time field').should('exist');
    },
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
              constantValue: 'SHOW DATABASES',
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
