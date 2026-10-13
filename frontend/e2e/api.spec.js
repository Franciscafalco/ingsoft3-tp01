import { expect, test } from '@playwright/test'

const API = process.env.API_BASE_URL || 'http://localhost:8081'
const GASTOS = `${API}/api/gastos`
const gasto = (descripcion) => ({ monto: 123.45, categoria: 'comida', descripcion, fecha: '2020-01-01T00:00:00Z' })

async function listar (request) {
  const r = await request.get(GASTOS)
  expect(r.status()).toBe(200)
  return await r.json()
}

test('el alta guarda en la base de verdad, y el borrado la saca', async ({ request }) => {
  const descripcion = `api ${Date.now()}`
  const alta = await request.post(GASTOS, { data: gasto(descripcion) })
  expect(alta.status()).toBe(201)
  const creado = await alta.json()
  expect((await listar(request)).some(g => g.descripcion === descripcion)).toBe(true)
  expect((await request.delete(`${GASTOS}/${creado.id}`)).status()).toBe(204)
  expect((await listar(request)).some(g => g.descripcion === descripcion)).toBe(false)
})

test('un monto en cero lo rechaza la api, y no crea nada', async ({ request }) => {
  const antes = await listar(request)
  const alta = await request.post(GASTOS, { data: { ...gasto(`api ${Date.now()}`), monto: 0 } })
  expect(alta.status()).toBe(400)
  expect((await listar(request)).length).toBe(antes.length)
})

test('editar un gasto persiste el cambio en la base, y se limpia', async ({ request }) => {
  const descripcion = `api ${Date.now()} edit`
  const creado = await (await request.post(GASTOS, { data: gasto(descripcion) })).json()
  const put = await request.put(`${GASTOS}/${creado.id}`, {
    data: { ...gasto(descripcion), monto: 200, categoria: 'transporte', estado: 'pendiente' },
  })
  expect(put.status()).toBe(200)
  const editado = (await listar(request)).find(g => g.descripcion === descripcion)
  expect(editado.monto).toBe(200)
  expect(editado.categoria).toBe('transporte')
  expect((await request.delete(`${GASTOS}/${creado.id}`)).status()).toBe(204)
  expect((await listar(request)).some(g => g.descripcion === descripcion)).toBe(false)
})