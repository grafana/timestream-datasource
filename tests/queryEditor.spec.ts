import { test, expect } from '@grafana/plugin-e2e';

test('should return data when a valid query is successfully run', async ({ page, panelEditPage, selectors }) => {
  await panelEditPage.datasource.set('AWS Timestream E2E');
  await panelEditPage.timeRange.set({ from: 'now-1h', to: 'now' });

  // @grafana/plugin-e2e@3.12.0's setVisualization() assumes every Grafana
  // >=12.4.0 build ships the tab-based "All visualizations" picker, and hangs
  // waiting for that tab on Grafana 12.4.10 (an enterprise patch in CI's e2e
  // matrix), which still shows the pre-12.4 single-button picker. Fall back
  // to driving that picker directly when the built-in helper times out.
  //
  // Remove this once a @grafana/plugin-e2e release fixes the version
  // threshold for Grafana 12.4.x builds without the tab-based picker.
  try {
    await panelEditPage.setVisualization('Table');
  } catch {
    await panelEditPage.getByGrafanaSelector(selectors.components.PanelEditor.toggleVizPicker).click();
    await panelEditPage.getByGrafanaSelector(selectors.components.PluginVisualization.item('Table')).click();
  }

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
