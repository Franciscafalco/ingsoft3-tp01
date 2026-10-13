import { expect, test } from '@playwright/test'

async function cargarGasto (page, descripcion, monto) {
  await page.getByLabel('Monto').fill(monto)
  await page.getByLabel('Descripción').fill(descripcion)
  await page.getByRole('button', { name: 'Agregar gasto' }).click()
}

test('crear un gasto lo muestra en la lista, y eliminarlo lo saca', async ({ page }) => {
  await page.goto('/')
  const descripcion = `e2e ${Date.now()}`
  await cargarGasto(page, descripcion, '123.45')
  await expect(page.getByText(descripcion)).toBeVisible()
  await page.getByRole('button', { name: `Eliminar ${descripcion}` }).click()
  await expect(page.getByText(descripcion)).toHaveCount(0)
})

test('un monto en cero muestra el error y no deja crear nada', async ({ page }) => {
  await page.goto('/')
  const descripcion = `e2e ${Date.now()} invalido`
  await page.getByLabel('Monto').fill('0')
  await page.getByLabel('Descripción').fill(descripcion)
  await expect(page.getByRole('alert')).toHaveText('El monto debe ser mayor a 0')
  await expect(page.getByRole('button', { name: 'Agregar gasto' })).toBeDisabled()
  await expect(page.getByText(descripcion)).toHaveCount(0)
})

test('el gasto sigue en la lista después de recargar la página', async ({ page }) => {
  await page.goto('/')
  const descripcion = `e2e ${Date.now()} persiste`
  await cargarGasto(page, descripcion, '77.5')
  await expect(page.getByText(descripcion)).toBeVisible()
  await page.reload()
  await expect(page.getByText(descripcion)).toBeVisible()
  await page.getByRole('button', { name: `Eliminar ${descripcion}` }).click()
  await expect(page.getByText(descripcion)).toHaveCount(0)
})