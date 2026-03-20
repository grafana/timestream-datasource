import { test, expect } from '@grafana/plugin-e2e';

test('should return data when a valid query is successfully run', async ({ page, panelEditPage, selectors }) => {
  await panelEditPage.datasource.set('AWS Timestream E2E');
  await panelEditPage.timeRange.set({ from: 'now-1h', to: 'now' });
  await panelEditPage.setVisualization('Table');

  await page.waitForFunction(() => window.monaco);
  const editor = panelEditPage.getByGrafanaSelector(selectors.components.CodeEditor.container);
  await editor.click();
  await page.keyboard.press('ControlOrMeta+A');
  await page.keyboard.insertText(
    `select time, measure_value::bigint from $__database.$__table where $__timeFilter and measure_value::bigint > 0 order by time asc limit 10`
  );

  await expect(panelEditPage.refreshPanel()).toBeOK();
  await expect(panelEditPage.panel.fieldNames).toHaveText(['time', 'measure_value::bigint']);
  await expect(panelEditPage.panel.data).toContainText([
    /\d{4}(-\d{2}){2} \d{2}(:\d{2}){2}\.\d{3}/ /* matches timestamp pattern e.g. '2025-01-15 09:03:36.654' */,
    /^\d+$/ /* matches BIGINT measure values */,
  ]);
});
