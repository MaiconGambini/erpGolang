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

test.describe('products', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, acme)
    await page.goto('/products')
    await expect(page.getByRole('heading', { name: 'Produtos' })).toBeVisible()
  })

  test('create product appears in list', async ({ page }) => {
    const name = `Produto E2E ${Date.now()}`
    const sku = `SKU-${Date.now()}`
    await page.getByRole('button', { name: 'Novo produto' }).click()
    await page.getByLabel('Nome *').fill(name)
    await page.getByLabel('SKU *').fill(sku)
    await page.getByLabel('Preço *').fill('19.90')
    await page.getByLabel('Estoque').fill('10')
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByText(name)).toBeVisible()
    await expect(page.getByText(sku)).toBeVisible()
  })

  test('edit product updates list', async ({ page }) => {
    const name = `Edit Prod ${Date.now()}`
    const sku = `EDIT-${Date.now()}`
    await page.getByRole('button', { name: 'Novo produto' }).click()
    await page.getByLabel('Nome *').fill(name)
    await page.getByLabel('SKU *').fill(sku)
    await page.getByLabel('Preço *').fill('9.99')
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByText(name)).toBeVisible()

    await page.getByRole('button', { name: 'Editar' }).first().click()
    const updated = `${name} Updated`
    await page.getByLabel('Nome *').fill(updated)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByText(updated)).toBeVisible()
  })
})

test.describe('products tenant isolation', () => {
  test('tenant B cannot see tenant A product', async ({ browser }) => {
    const name = `Isolated Prod ${Date.now()}`
    const sku = `ISO-${Date.now()}`
    const acmeCtx = await browser.newContext()
    const betaCtx = await browser.newContext()
    const acmePage = await acmeCtx.newPage()
    const betaPage = await betaCtx.newPage()

    await login(acmePage, acme)
    await acmePage.goto('/products')
    await acmePage.getByRole('button', { name: 'Novo produto' }).click()
    await acmePage.getByLabel('Nome *').fill(name)
    await acmePage.getByLabel('SKU *').fill(sku)
    await acmePage.getByLabel('Preço *').fill('15.00')
    await acmePage.getByRole('button', { name: 'Salvar' }).click()
    await expect(acmePage.getByText(name)).toBeVisible()

    await login(betaPage, beta)
    await betaPage.goto('/products')
    await betaPage.getByPlaceholder('Buscar por nome ou SKU').fill(name)
    await expect(betaPage.getByText(name)).toHaveCount(0)

    await acmeCtx.close()
    await betaCtx.close()
  })
})
