import { describe, it, expect } from 'vitest'
import { clasificarGasto } from './clasificarGasto'

describe('clasificarGasto', () => {
  it.each([
    ['un valor que no es número', 'abc', 100, 'invalido'],
    ['un monto en cero', 0, 100, 'invalido'],
    ['un monto 3 veces el promedio', 300, 100, 'muy-alto'],
    ['un monto 1.5 veces el promedio', 150, 100, 'alto'],
    ['un monto cerca del promedio', 60, 100, 'normal'],
    ['un monto muy por debajo del promedio', 10, 100, 'bajo'],
    ['sin promedio y un monto grande', 10000, 0, 'alto'],
    ['sin promedio y un monto chico', 500, 0, 'normal'],
  ])('%s => %s', (_caso, monto, promedio, esperado) => {
    expect(clasificarGasto(monto, promedio)).toBe(esperado)
  })
})