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

test.describe('csv export', () => {
  test('customers export downloads CSV', async ({ page }) => {
    await login(page, acme)
    await page.goto('/customers')
    await expect(page.getByRole('heading', { name: 'Clientes' })).toBeVisible()

    const downloadPromise = page.waitForEvent('download')
    await page.getByRole('button', { name: 'Exportar CSV' }).click()
    const download = await downloadPromise
    expect(download.suggestedFilename()).toMatch(/clientes\.csv$/i)
  })
})
