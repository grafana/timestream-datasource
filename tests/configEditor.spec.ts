import { test, expect } from '@grafana/plugin-e2e';
import { lt } from 'semver';

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
  grafanaVersion,
  page,
  readProvisionedDataSource,
}) => {
  // Minting shipped in Grafana 13.2. Older images in the CI matrix cannot produce this ID.
  test.skip(lt(grafanaVersion, '13.2.0'), 'Per-datasource external ID minting requires Grafana 13.2 or newer');

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

  // Grafana 13.3+ appends a 16-char hex suffix. The published nightly image predates that and mints {stack}-{uid}.
  expect(minted).toMatch(new RegExp(`-${uid}(-[0-9a-f]{16})?$`));
  expect(saved.datasource.jsonData.usePerDatasourceExternalId).toBe(true);
  await expect(page.locator('#grafanaExternalId')).toHaveValue(minted);
});
