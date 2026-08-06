import { test, expect } from '@grafana/plugin-e2e';

// Disable the new dashboard layouts so plugin-e2e's addPanel() helper uses the
// stable "Add panel" flow instead of the sidebar/edit-pane path.
//
// Grafana 13.2.0 (the nightly build in CI's e2e matrix) renamed the sidebar
// test ids, e.g. "edit pane configure panel button" -> "sidebar configure panel
// button". @grafana/plugin-e2e@3.10.0 pins @grafana/e2e-selectors@13.1.0, which
// only knows the old name, so addPanel() waits for a test id that no longer
// exists and panelEditPage setup times out. No stable e2e-selectors release
// carries the new name yet (only Grafana's nightly prerelease does).
//
// Remove this override once we upgrade to a @grafana/plugin-e2e release whose
// bundled @grafana/e2e-selectors includes the Grafana 13.2.0 sidebar selectors.
test.use({ featureToggles: { dashboardNewLayouts: false } });

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
