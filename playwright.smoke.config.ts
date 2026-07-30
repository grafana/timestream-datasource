import { defineConfig } from '@playwright/test';
import type { PluginOptions } from '@grafana/plugin-e2e';
import baseConfig from './playwright.config';

/**
 * Smoke test config used by fork PRs.
 *
 * Fork PRs don't have access to the Vault secrets (AWS credentials) that the
 * full e2e suite needs, so instead they run the secret-free smoke tests in
 * ./tests/smoke, which only validate that the config, query and variable
 * editors load correctly.
 */
export default defineConfig<PluginOptions>({
  ...baseConfig,
  testDir: './tests/smoke',
});
