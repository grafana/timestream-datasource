import { expect, test } from '@grafana/plugin-e2e';

test('should successfully create a variable', async ({ variableEditPage, page, selectors }) => {
  await variableEditPage.datasource.set('AWS Timestream E2E');
  const editor = page.getByTestId(
    selectors.pages.Dashboard.Settings.Variables.Edit.QueryVariable.queryOptionsQueryInput
  );
  await expect(editor).toBeVisible();
  await editor.click();
  await page.keyboard.insertText('SHOW TABLES FROM grafanaDB');
  const queryDataRequest = variableEditPage.waitForQueryDataRequest();
  await variableEditPage.runQuery();
  await queryDataRequest;
  await expect(variableEditPage).toDisplayPreviews(['DevOps', 'IoT']);
});
