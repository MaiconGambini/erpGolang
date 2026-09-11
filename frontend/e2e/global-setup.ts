import { execSync } from 'node:child_process'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const backendDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../../backend')

const env = {
  ...process.env,
  APP_ENV: process.env.APP_ENV ?? 'development',
  DATABASE_URL:
    process.env.DATABASE_URL ??
    'postgres://goerp:goerp@localhost:5434/goerp?sslmode=disable',
  REDIS_URL: process.env.REDIS_URL ?? 'redis://localhost:6381/0',
  JWT_ACCESS_SECRET: process.env.JWT_ACCESS_SECRET ?? 'test-access',
  JWT_REFRESH_SECRET: process.env.JWT_REFRESH_SECRET ?? 'test-refresh',
}

export default async function globalSetup() {
  execSync('go run ./cmd/migrate', { cwd: backendDir, env, stdio: 'inherit' })
  execSync('go run ./cmd/seed', { cwd: backendDir, env, stdio: 'inherit' })
}
