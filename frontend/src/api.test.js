import { describe, it, expect, vi } from 'vitest'
import { crearGasto } from './api'

describe('crearGasto', () => {
  it('le pide a la API un POST a /api/gastos con el gasto en el body', async () => {
    const fetcher = vi.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: () => Promise.resolve({ id: 1, monto: 100 }),
    })

    await crearGasto({ monto: 100, categoria: 'comida' }, fetcher)

    expect(fetcher).toHaveBeenCalledWith('/api/gastos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ monto: 100, categoria: 'comida' }),
    })
  })
})