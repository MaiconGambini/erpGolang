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

test.describe('suppliers', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, acme)
    await page.goto('/suppliers')
    await expect(page.getByRole('heading', { name: 'Fornecedores' })).toBeVisible()
  })

  test('create supplier appears in list', async ({ page }) => {
    const name = `Fornecedor E2E ${Date.now()}`
    await page.getByRole('button', { name: 'Novo fornecedor' }).click()
    await page.getByLabel('Nome *').fill(name)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByRole('dialog')).toHaveCount(0)
    await page.getByPlaceholder('Buscar por nome ou documento').fill(name)
    await expect(page.getByRole('cell', { name })).toBeVisible()
  })

  test('edit supplier updates list', async ({ page }) => {
    const name = `Edit Forn ${Date.now()}`
    await page.getByRole('button', { name: 'Novo fornecedor' }).click()
    await page.getByLabel('Nome *').fill(name)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByRole('dialog')).toHaveCount(0)
    await page.getByPlaceholder('Buscar por nome ou documento').fill(name)
    await expect(page.getByRole('cell', { name })).toBeVisible()

    await page.getByRole('button', { name: 'Editar' }).first().click()
    const updated = `${name} Updated`
    await page.getByLabel('Nome *').fill(updated)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByText(updated)).toBeVisible()
  })
})

test.describe('suppliers tenant isolation', () => {
  test('tenant B cannot see tenant A supplier', async ({ browser }) => {
    const name = `Isolated Forn ${Date.now()}`
    const acmeCtx = await browser.newContext()
    const betaCtx = await browser.newContext()
    const acmePage = await acmeCtx.newPage()
    const betaPage = await betaCtx.newPage()

    await login(acmePage, acme)
    await acmePage.goto('/suppliers')
    await acmePage.getByRole('button', { name: 'Novo fornecedor' }).click()
    await acmePage.getByLabel('Nome *').fill(name)
    await acmePage.getByRole('button', { name: 'Salvar' }).click()
    await expect(acmePage.getByRole('dialog')).toHaveCount(0)
    await acmePage.getByPlaceholder('Buscar por nome ou documento').fill(name)
    await expect(acmePage.getByRole('cell', { name })).toBeVisible()

    await login(betaPage, beta)
    await betaPage.goto('/suppliers')
    await betaPage.getByPlaceholder('Buscar por nome ou documento').fill(name)
    await expect(betaPage.getByText(name)).toHaveCount(0)

    await acmeCtx.close()
    await betaCtx.close()
  })
})
