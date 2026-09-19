import { describe, it, expect } from 'vitest'
import { validarMonto, validarFecha } from './validarGasto'

describe('validarMonto', () => {
  it.each([
    ['un monto positivo', '100', true],
    ['un monto en cero', '0', false],
    ['un monto negativo', '-50', false],
    ['un monto vacío se considera válido (todavía no se completó el campo)', '', true],
  ])('%s => válido: %s', (_caso, monto, esperado) => {
    const resultado = validarMonto(monto)
    expect(resultado.valido).toBe(esperado)
  })

  it('un monto inválido explica el error', () => {
    const resultado = validarMonto('-10')
    expect(resultado.valido).toBe(false)
    expect(resultado.error).toBe('El monto debe ser mayor a 0')
  })
})

describe('validarFecha', () => {
  it.each([
    ['una fecha pasada', '2020-01-01', true],
    ['una fecha futura', '2099-01-01', false],
  ])('%s => válido: %s', (_caso, fecha, esperado) => {
    const resultado = validarFecha(fecha)
    expect(resultado.valido).toBe(esperado)
  })

  it('una fecha futura explica el error', () => {
    const resultado = validarFecha('2099-01-01')
    expect(resultado.valido).toBe(false)
    expect(resultado.error).toBe('La fecha no puede ser futura')
  })
})