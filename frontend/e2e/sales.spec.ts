import { expect, test, type Page } from '@playwright/test'

async function login(
  page: Page,
  tenant: { slug: string; email: string; password: string },
) {
  await page.goto('/login')
  await page.getByRole('textbox', { name: 'Empresa' }).fill(tenant.slug)
  await page.getByRole('textbox', { name: 'E-mail' }).fill(tenant.email)
  await page.getByLabel('Senha').fill(tenant.password)
  await page.getByRole('button', { name: 'Entrar' }).click()
  await page.waitForURL('/')
}

async function createDraftSale(page: Page, customerName: string, sku: string) {
  await page.getByRole('button', { name: 'Nova venda' }).click()
  const dialog = page.getByRole('dialog')
  await expect(dialog.getByLabel('Cliente *').locator('option', { hasText: customerName })).toHaveCount(1, { timeout: 15000 })
  await dialog.getByLabel('Cliente *').selectOption({ label: customerName })
  const productSelect = dialog.locator('fieldset select').first()
  await expect(productSelect.locator('option', { hasText: sku })).toHaveCount(1, { timeout: 15000 })
  const productValue = await productSelect.locator('option', { hasText: sku }).getAttribute('value')
  await productSelect.selectOption(productValue!)
  await page.getByRole('button', { name: 'Salvar rascunho' }).click()
}

const acme = { slug: 'acme', email: 'admin@acme.com', password: 'admin123' }
const beta = { slug: 'beta', email: 'admin@beta.com', password: 'admin123' }

test.describe('sales', () => {
  test.beforeEach(async ({ page }) => {
    await login(page, acme)
    await page.goto('/sales')
    await expect(page.getByRole('heading', { name: 'Vendas' })).toBeVisible()
  })

  test('create and confirm sale', async ({ page }) => {
    test.setTimeout(60_000)
    const customerName = `Cliente Venda ${Date.now()}`
    await page.goto('/customers')
    await page.getByRole('button', { name: 'Novo cliente' }).click()
    await page.getByLabel('Nome *').fill(customerName)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByRole('cell', { name: customerName })).toBeVisible()

    const sku = `SALE-SKU-${Date.now()}`
    await page.goto('/products')
    await page.getByRole('button', { name: 'Novo produto' }).click()
    await page.getByLabel('Nome *').fill(`Produto ${sku}`)
    await page.getByLabel('SKU *').fill(sku)
    await page.getByLabel('Preço *').fill('10.00')
    await page.getByLabel('Estoque').fill('20')
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByRole('cell', { name: sku, exact: true })).toBeVisible()

    await page.goto('/sales')
    await createDraftSale(page, customerName, sku)
    await expect(page.getByRole('cell', { name: customerName })).toBeVisible()
    await expect(page.locator('.badge.draft').first()).toBeVisible()

    await page.getByRole('button', { name: 'Confirmar' }).first().click()
    await expect(page.locator('.badge.confirmed').first()).toBeVisible()
  })
})

test.describe('sales tenant isolation', () => {
  test('tenant B cannot see tenant A sale', async ({ browser }) => {
    test.setTimeout(60_000)
    const customerName = `Iso Sale ${Date.now()}`
    const acmeCtx = await browser.newContext()
    const betaCtx = await browser.newContext()
    const acmePage = await acmeCtx.newPage()
    const betaPage = await betaCtx.newPage()

    await login(acmePage, acme)
    await acmePage.goto('/customers')
    await acmePage.getByRole('button', { name: 'Novo cliente' }).click()
    await acmePage.getByLabel('Nome *').fill(customerName)
    await acmePage.getByRole('button', { name: 'Salvar' }).click()

    const sku = `ISO-${Date.now()}`
    await acmePage.goto('/products')
    await acmePage.getByRole('button', { name: 'Novo produto' }).click()
    await acmePage.getByLabel('Nome *').fill(`P ${sku}`)
    await acmePage.getByLabel('SKU *').fill(sku)
    await acmePage.getByLabel('Preço *').fill('5.00')
    await acmePage.getByLabel('Estoque').fill('5')
    await acmePage.getByRole('button', { name: 'Salvar' }).click()

    await acmePage.goto('/sales')
    await createDraftSale(acmePage, customerName, sku)
    await expect(acmePage.getByText(customerName)).toBeVisible()

    await login(betaPage, beta)
    await betaPage.goto('/sales')
    await betaPage.getByPlaceholder('Buscar por cliente').fill(customerName)
    await expect(betaPage.getByText(customerName)).toHaveCount(0)

    await acmeCtx.close()
    await betaCtx.close()
  })
})
