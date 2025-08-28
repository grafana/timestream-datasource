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

  // Seems there is a racecondition where the onBlur from the SQLEditor triggers a POST to '/api/v1/eval'
  // This causes flaky errors like Error: route.fulfill: Route is already handled!
  // In order to fix the flakiness this lets make sure we click outside the sql editor and wait for the first POST to '/api/v1/eval' before we evalute
  await page.mouse.click(0, 0); // forces a click outside the SqlEditor forcing the onBlur handler
  await page.waitForRequest(selectors.apis.Alerting.eval); // wait for the onBlur handlers call to '/api/v1/eval'

  await expect(alertRuleEditPage.evaluate()).toBeOK();
});
