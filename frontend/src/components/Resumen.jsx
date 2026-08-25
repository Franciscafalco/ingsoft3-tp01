const formatoMoneda = new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' })

export default function Resumen({ resumen }) {
  if (!resumen) return <p className="vacio">Cargando resumen…</p>

  const categorias = Object.entries(resumen.porCategoria).sort((a, b) => b[1] - a[1])

  return (
    <div className="resumen">
      <h2>Resumen por categoría</h2>
      <p className="total-general">Total general: {formatoMoneda.format(resumen.totalGeneral)}</p>

      {categorias.length === 0 ? (
        <p className="vacio">Sin datos todavía.</p>
      ) : (
        <ul className="resumen-lista">
          {categorias.map(([categoria, total]) => (
            <li key={categoria}>
              <span>{categoria}</span>
              <span>{formatoMoneda.format(total)}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
