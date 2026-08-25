const formatoMoneda = new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' })

export default function GastosTable({ gastos, onMarcarPagado, onEliminar }) {
  if (gastos.length === 0) {
    return <p className="vacio">Todavía no cargaste ningún gasto.</p>
  }

  return (
    <table className="gastos-table">
      <thead>
        <tr>
          <th>Fecha</th>
          <th>Categoría</th>
          <th>Descripción</th>
          <th>Monto</th>
          <th>Estado</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {gastos.map((g) => (
          <tr key={g.id} className={g.estado === 'pagado' ? 'fila-pagada' : ''}>
            <td>{new Date(g.fecha).toLocaleDateString('es-AR')}</td>
            <td>{g.categoria}</td>
            <td>{g.descripcion || '—'}</td>
            <td>{formatoMoneda.format(g.monto)}</td>
            <td>
              <span className={`badge badge-${g.estado}`}>{g.estado}</span>
            </td>
            <td className="acciones">
              {g.estado === 'pendiente' && (
                <button type="button" onClick={() => onMarcarPagado(g)}>
                  Marcar pagado
                </button>
              )}
              <button
                type="button"
                disabled={g.estado === 'pagado'}
                title={g.estado === 'pagado' ? 'No se puede eliminar un gasto pagado' : undefined}
                onClick={() => onEliminar(g)}
              >
                Eliminar
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
