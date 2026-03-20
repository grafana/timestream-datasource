import { expect, test } from '@grafana/plugin-e2e';
import { gte } from 'semver';

test('should successfully create an annotation', async ({ annotationEditPage, grafanaVersion, page, selectors }) => {
  await annotationEditPage.datasource.set('AWS Timestream E2E');
  await page.waitForFunction(() => window.monaco);
  await page.getByLabel('Query', { exact: true }).fill('First 10 rows');
  await expect(annotationEditPage.getByGrafanaSelector(selectors.components.Select.option)).toContainText([
    'First 10 rows',
  ]);
  await page.keyboard.press('Enter');
  await expect(annotationEditPage.runQuery()).toBeOK();
  if (gte(grafanaVersion, '11.0.0')) {
    // grafana_example_table has 6 columns: time, measure_name, measure_value::bigint, host, datacenter, app
    await expect(annotationEditPage).toHaveAlert('success', { hasText: '10 events (from 6 fields)' });
  }
});
