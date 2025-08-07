import { test, expect } from '@grafana/plugin-e2e';

test.use({
  featureToggles: {
    alertingQueryAndExpressionsStepMode: false,
  },
});

test('should successfully create an alert rule', async ({
  alertRuleEditPage,
  page,
  readProvisionedDataSource,
  selectors,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-timestream-e2e.yaml', name: 'AWS Timestream E2E' });
  const queryA = await alertRuleEditPage.getQueryRow('A');
  await queryA.datasource.set(ds.name);
  await page.waitForFunction(() => window.monaco);
  await queryA.getByGrafanaSelector(selectors.components.CodeEditor.container).click();
  await page.keyboard.insertText(
    `select region, avg(measure_value::double) from $__database.$__table where time between from_milliseconds(1615395600000) and from_milliseconds(1615395900000) and measure_value::double > 1 group by region limit 10`
  );
  await expect(alertRuleEditPage.evaluate()).toBeOK();
});
