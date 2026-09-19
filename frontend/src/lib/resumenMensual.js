export function resumenMensual(gastos, mes, anio) {
  if (!Array.isArray(gastos)) {
    return { total: 0, cantidad: 0, mayor: null, nivel: 'sin-datos' }
  }
  const delMes = gastos.filter((g) => {
    if (!g || !g.fecha) {
      return false
    }
    const f = new Date(g.fecha)
    return f.getMonth() + 1 === mes && f.getFullYear() === anio
  })
  if (delMes.length === 0) {
    return { total: 0, cantidad: 0, mayor: null, nivel: 'sin-datos' }
  }
  let total = 0
  let mayor = delMes[0]
  for (const g of delMes) {
    total += g.monto
    if (g.monto > mayor.monto) {
      mayor = g
    }
  }
  let nivel = 'bajo'
  if (total >= 100000) {
    nivel = 'critico'
  } else if (total >= 50000) {
    nivel = 'alto'
  } else if (total >= 10000) {
    nivel = 'medio'
  }
  return { total, cantidad: delMes.length, mayor, nivel }
}