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

const viewer = { slug: 'acme', email: 'viewer@acme.com', password: 'admin123' }

test.describe('viewer role', () => {
  test('cannot create customers (write actions hidden)', async ({ page }) => {
    await login(page, viewer)
    await page.goto('/customers')
    await expect(page.getByRole('heading', { name: 'Clientes' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Novo cliente' })).toHaveCount(0)
    await expect(page.getByRole('button', { name: 'Excluir' })).toHaveCount(0)
  })

  test('cannot access admin users page', async ({ page }) => {
    await login(page, viewer)
    await page.goto('/users')
    await expect(page).toHaveURL('/')
  })

  test('cannot access admin audit page', async ({ page }) => {
    await login(page, viewer)
    await page.goto('/audit')
    await expect(page).toHaveURL('/')
  })

  test('does not see financial dashboard charts', async ({ page }) => {
    await login(page, viewer)
    await expect(page.getByRole('heading', { name: 'Vendas por dia' })).toHaveCount(0)
    await expect(page.getByRole('button', { name: 'Exportar PDF' })).toHaveCount(0)
  })
})
