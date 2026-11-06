import { test, expect } from '@grafana/plugin-e2e';

/**
 * Smoke test: validates that the query editor renders without needing real
 * AWS credentials. It only checks that the code editor loads, it does not run
 * a query against the backend.
 */
test('query editor should load', async ({ page, panelEditPage, selectors }) => {
  await panelEditPage.datasource.set('AWS Timestream E2E');

  await page.waitForFunction(() => window.monaco);
  const editor = panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container);
  await expect(editor).toBeVisible();
});
