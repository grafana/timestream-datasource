import { e2e } from '@grafana/e2e';
import { selectors as timestreamSelectors } from '../../src/components/selectors';

const e2eSelectors = e2e.getSelectors(timestreamSelectors.components);

export const addDataSourceWithKey = (
  datasourceType: string,
  accessKey: string,
  secretKey: string,
  region: string
): any => {
  return e2e.flows.addDataSource({
    checkHealth: false,
    expectedAlertMessage: 'Connection success',
    form: () => {
      setSelectValue('.aws-config-authType', 'Access & secret key');
      e2eSelectors.ConfigEditor.AccessKey.input().type(accessKey);
      e2eSelectors.ConfigEditor.SecretKey.input().type(secretKey);
      setSelectValue('.aws-config-defaultRegion', region);
    },
    type: datasourceType,
  });
};

const setSelectValue = (container: string, text: string) => {
  // return e2e.flows.selectOption({
  //   clickToOpen: true,
  //   optionText: text,
  //   container: e2e().get(container),
  // });

  // couldn't get above code to work for some reason. need to investigate that
  return e2e().get(container).parent().find(`input`).click({ force: true }).type(text).type('{enter}');
};

const addTablePanel = (query: string) => {
  const fillQuery = () => e2eSelectors.QueryEditor.CodeEditor.container().type(query);

  e2e.flows.addPanel({
    matchScreenshot: true,
    visualizationName: e2e.flows.VISUALIZATION_TABLE,
    queriesForm: () => fillQuery(),
  });

  e2e.flows.explore({
    matchScreenshot: false,
    timeRange: {
      from: '2001-01-31 19:00:00',
      to: '2016-01-31 19:00:00',
    },
    queriesForm: () => fillQuery(),
  });
};

e2e.scenario({
  describeName: 'Smoke tests',
  itName: 'Login, create data source, dashboard and panel',
  scenario: () => {
    e2e()
      .readProvisions(['datasources/aws-timestream.yaml'])
      .then(([provision]) => {
        const datasource = provision.datasources[0];
        return addDataSourceWithKey(
          'Amazon Timestream',
          datasource.secureJsonData.accessKey,
          datasource.secureJsonData.secretKey,
          datasource.jsonData.defaultRegion
        );
      })
      .then(() => {
        e2e.flows.addDashboard({
          timeRange: {
            from: '2001-01-31 19:00:00',
            to: '2016-01-31 19:00:00',
          },
        });
        addTablePanel(`SHOW DATABASES`);
      });
  },
});
