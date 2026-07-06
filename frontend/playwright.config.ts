import { defineConfig, devices } from '@playwright/test'

const e2ePort = 5174
const databaseUrl =
  process.env.DATABASE_URL ??
  'postgres://goerp:goerp@localhost:5434/goerp?sslmode=disable'
const redisUrl = process.env.REDIS_URL ?? 'redis://localhost:6381/0'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  workers: 1,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: 'list',
  use: {
    baseURL: `http://localhost:${e2ePort}`,
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
        DATABASE_URL: databaseUrl,
        REDIS_URL: redisUrl,
        JWT_ACCESS_SECRET: 'test-access',
        JWT_REFRESH_SECRET: 'test-refresh',
      },
    },
    {
      command: `npm run dev -- --port ${e2ePort} --strictPort`,
      url: `http://localhost:${e2ePort}`,
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],
})
