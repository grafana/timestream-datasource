import { test, expect } from '@grafana/plugin-e2e';

/**
 * Smoke test: validates that the config editor renders without needing real
 * AWS credentials. Fork PRs run these tests because they don't have access to
 * the Vault secrets required by the full e2e suite.
 */
test('config editor should load', async ({ createDataSourceConfigPage, page, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'aws-timestream-e2e.yaml', name: 'AWS Timestream E2E' });
  await createDataSourceConfigPage({ type: ds.type });

  const authProvider = page.getByLabel(/^Authentication Provider/);
  await expect(authProvider).toBeVisible();

  await authProvider.fill('Access & secret key');
  await page.keyboard.press('Enter');

  await expect(page.getByLabel('Access Key ID')).toBeVisible();
  await expect(page.getByLabel('Secret Access Key')).toBeVisible();
  await expect(page.getByLabel('Default Region')).toBeVisible();
});
