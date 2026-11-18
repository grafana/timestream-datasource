import { test, expect } from '@grafana/plugin-e2e';

test('should successfully provision a data source', async ({ gotoDataSourceConfigPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-timestream-e2e.yaml', name: 'AWS Timestream E2E' });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  await expect(configPage.saveAndTest()).toBeOK();
});

test('should return an error when invalid credentials are used', async ({
  createDataSourceConfigPage,
  page,
  readProvisionedDataSource,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-timestream-e2e.yaml', name: 'AWS Timestream E2E' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });
  await page.getByLabel(/^Authentication Provider/).fill('Access & secret key');
  await page.keyboard.press('Enter');
  await page.getByLabel('Access Key ID').fill('invalid_access_id');
  await page.getByLabel('Secret Access Key').fill('invalid_secret_key');
  await page.getByLabel('Default Region').fill('us-east-2');
  await page.keyboard.press('Enter');
  await expect(configPage.saveAndTest()).not.toBeOK();
});

test('shows the server-minted external ID after save', async ({
  createDataSourceConfigPage,
  page,
  readProvisionedDataSource,
}) => {
  // Plugin CI starts a local Grafana. Minting happens on datasourcese2e once
  // awsAssumeRolePerDatasourceExternalId is rolled out there.
  test.skip(!process.env.GRAFANA_URL, 'Requires datasourcese2e.grafana-dev.net');

  const ds = await readProvisionedDataSource({ fileName: 'aws-timestream-e2e.yaml', name: 'AWS Timestream E2E' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });
  const uid = configPage.datasource.uid;

  await page.getByLabel(/^Authentication Provider/).fill('Grafana Assume Role');
  await page.keyboard.press('Enter');
  await expect(page.locator('#grafanaExternalId')).toHaveValue('Save the data source to generate an external ID');

  const save = page.waitForResponse(
    (response) =>
      response.request().method() === 'PUT' && response.url().includes(`/api/datasources/uid/${uid}`) && response.ok()
  );
  await page.getByRole('button', { name: 'Save & test' }).click();
  const saved = await (await save).json();
  const minted = String(saved.datasource.jsonData.grafanaExternalId);

  expect(minted).toMatch(new RegExp(`-${uid}-[0-9a-f]{16}$`));
  expect(saved.datasource.jsonData.usePerDatasourceExternalId).toBe(true);
  await expect(page.locator('#grafanaExternalId')).toHaveValue(minted);
});
