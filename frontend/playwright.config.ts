import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  workers: 1,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: 'list',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
  ],
  webServer: [
    {
      command: process.platform === 'win32'
        ? 'cd ../backend && .\\api.exe'
        : 'cd ../backend && go run ./cmd/api',
      url: 'http://localhost:8080/healthz',
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
      env: {
        DATABASE_URL: 'postgres://goerp:goerp@localhost:5432/goerp?sslmode=disable',
        REDIS_URL: 'redis://localhost:6379/0',
        JWT_ACCESS_SECRET: 'test-access',
        JWT_REFRESH_SECRET: 'test-refresh',
      },
    },
    {
      command: 'npm run dev',
      url: 'http://localhost:5173',
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],
})
