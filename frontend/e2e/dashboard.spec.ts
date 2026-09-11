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

test.describe('dashboard', () => {
  test('shows live KPI counts after login', async ({ page }) => {
    await login(page, acme)

    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()

    const kpiValues = page.locator('section.grid article strong:not(.skeleton)')
    await expect(kpiValues).toHaveCount(6)

    for (let i = 0; i < 5; i++) {
      const text = await kpiValues.nth(i).textContent()
      expect(text).toMatch(/^\d[\d.]*$/)
    }

    const revenueText = await kpiValues.nth(5).textContent()
    expect(revenueText).toMatch(/^R\$\s/)
  })

  test('shows chart sections after login', async ({ page }) => {
    await login(page, acme)

    await expect(page.getByRole('heading', { name: 'Vendas por dia' })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Produtos mais vendidos' })).toBeVisible()
    await expect(page.locator('section.charts canvas')).toHaveCount(2, { timeout: 15_000 })
  })

  test('active customers increments after create', async ({ page }) => {
    await login(page, acme)

    const kpi = () => page.locator('section.grid article').first().locator('strong:not(.skeleton)')
    await expect(kpi()).toHaveText(/\d/, { timeout: 10_000 })
    const before = Number((await kpi().textContent() ?? '0').replace(/\./g, ''))

    await page.goto('/customers')
    await expect(page.getByRole('heading', { name: 'Clientes' })).toBeVisible()

    const name = `Dashboard KPI ${Date.now()}`
    await page.getByRole('button', { name: 'Novo cliente' }).click()
    await page.getByLabel('Nome *').fill(name)
    await page.getByRole('button', { name: 'Salvar' }).click()
    await expect(page.getByRole('dialog')).toHaveCount(0)
    await page.getByPlaceholder('Buscar por nome, documento ou e-mail').fill(name)
    await expect(page.getByRole('cell', { name })).toBeVisible()

    await page.goto('/')
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
    await expect(kpi()).toHaveText(/\d/, { timeout: 10_000 })
    const after = Number((await kpi().textContent() ?? '0').replace(/\./g, ''))
    expect(after).toBeGreaterThanOrEqual(before + 1)
  })
})
