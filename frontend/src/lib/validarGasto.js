export function hoy() {
  return new Date().toISOString().slice(0, 10)
}

export function validarMonto(monto) {
  if (monto === '') return { valido: true }
  return {
    valido: Number(monto) > 0,
    error: 'El monto debe ser mayor a 0',
  }
}

export function validarFecha(fecha) {
  const esFutura = fecha > hoy()
  return {
    valido: !esFutura,
    error: 'La fecha no puede ser futura',
  }
}