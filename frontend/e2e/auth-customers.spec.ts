import { expect, test } from '@playwright/test'

async function login(
  page: import('@playwright/test').Page,
  tenant: { slug: string; email: string; password: string },
) {
  await page.goto('/login')
  await page.getByRole('textbox', { name: 'Empresa' }).fill(tenant.slug)
  await page.getByRole('textbox', { name: 'E-mail' }).fill(tenant.email)
  await page.getByLabel('Senha').fill(tenant.password)
  await page.getByRole('button', { name: 'Entrar' }).click()
  await page.waitForURL('/')
}

const acme = { slug: 'acme', email: 'admin@acme.com', password: 'admin123' }
const beta = { slug: 'beta', email: 'admin@beta.com', password: 'admin123' }

test.describe('auth', () => {
  test('redirects unauthenticated users to login', async ({ page }) => {
    await page.goto('/customers')
    await expect(page).toHaveURL(/\/login/)
  })

  test('login and reach customers page', async ({ page }) => {
    await login(page, acme)
    await page.goto('/customers')
    await page.waitForURL('**/customers')
    await expect(page.getByRole('heading', { name: 'Clientes' })).toBeVisible()
  })
})

test.describe('customers', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, acme)
    await page.goto('/customers')
    await expect(page.getByRole('heading', { name: 'Clientes' })).toBeVisible()
  })

  test('create customer appears in list', async ({ page }) => {
    const name = `Cliente E2E ${Date.now()}`
    await page.getByRole('button', { name: 'Novo cliente' }).click()
    await page.getByLabel('Nome *').fill(name)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByText(name)).toBeVisible()
  })

  test('edit customer updates list', async ({ page }) => {
    const name = `Edit E2E ${Date.now()}`
    await page.getByRole('button', { name: 'Novo cliente' }).click()
    await page.getByLabel('Nome *').fill(name)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByText(name)).toBeVisible()

    await page.getByRole('button', { name: 'Editar' }).first().click()
    const updated = `${name} Updated`
    await page.getByLabel('Nome *').fill(updated)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByText(updated)).toBeVisible()
  })
})

test.describe('tenant isolation', () => {
  test('tenant B cannot see tenant A customer', async ({ browser }) => {
    const name = `Isolated ${Date.now()}`
    const acmeCtx = await browser.newContext()
    const betaCtx = await browser.newContext()
    const acmePage = await acmeCtx.newPage()
    const betaPage = await betaCtx.newPage()

    await login(acmePage, acme)
    await acmePage.goto('/customers')
    await acmePage.getByRole('button', { name: 'Novo cliente' }).click()
    await acmePage.getByLabel('Nome *').fill(name)
    await acmePage.getByRole('button', { name: 'Salvar' }).click()
    await expect(acmePage.getByText(name)).toBeVisible()

    await login(betaPage, beta)
    await betaPage.goto('/customers')
    await betaPage.getByPlaceholder('Buscar por nome, documento ou e-mail').fill(name)
    await expect(betaPage.getByText(name)).toHaveCount(0)

    await acmeCtx.close()
    await betaCtx.close()
  })
})
