export function clasificarGasto(monto, promedio) {
  if (typeof monto !== 'number' || Number.isNaN(monto)) {
    return 'invalido'
  }
  if (monto <= 0) {
    return 'invalido'
  }
  if (promedio > 0) {
    const ratio = monto / promedio
    if (ratio >= 3) {
      return 'muy-alto'
    }
    if (ratio >= 1.5) {
      return 'alto'
    }
    if (ratio >= 0.5) {
      return 'normal'
    }
    return 'bajo'
  }
  if (monto >= 10000) {
    return 'alto'
  }
  return 'normal'
}