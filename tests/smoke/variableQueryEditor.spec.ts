import { expect, test } from '@grafana/plugin-e2e';

/**
 * Smoke test: validates that the variable query editor renders without needing
 * real AWS credentials. It only checks that the query input loads, it does not
 * run a query against the backend.
 */
test('variable query editor should load', async ({ variableEditPage, page, selectors }) => {
  await variableEditPage.datasource.set('AWS Timestream E2E');

  const editor = page.getByTestId(
    selectors.pages.Dashboard.Settings.Variables.Edit.QueryVariable.queryOptionsQueryInput
  );
  await expect(editor).toBeVisible();
});
