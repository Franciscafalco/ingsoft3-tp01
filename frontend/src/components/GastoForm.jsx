import { useState } from 'react'
import { CATEGORIAS } from '../constants'

const hoy = () => new Date().toISOString().slice(0, 10)

const vacio = {
  monto: '',
  categoria: CATEGORIAS[0],
  descripcion: '',
  fecha: hoy(),
}

export default function GastoForm({ onCrear }) {
  const [form, setForm] = useState(vacio)
  const [error, setError] = useState(null)
  const [enviando, setEnviando] = useState(false)

  const montoInvalido = form.monto !== '' && Number(form.monto) <= 0
  const fechaFutura = form.fecha > hoy()
  const puedeEnviar = form.monto !== '' && !montoInvalido && !fechaFutura && !enviando

  function actualizarCampo(campo, valor) {
    setForm((prev) => ({ ...prev, [campo]: valor }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    if (!puedeEnviar) return

    setEnviando(true)
    setError(null)
    try {
      await onCrear({
        ...form,
        monto: Number(form.monto),
        fecha: new Date(form.fecha).toISOString(),
      })
      setForm(vacio)
    } catch (err) {
      setError(err.message)
    } finally {
      setEnviando(false)
    }
  }

  return (
    <form className="gasto-form" onSubmit={handleSubmit}>
      <h2>Nuevo gasto</h2>

      <label>
        Monto
        <input
          type="number"
          step="0.01"
          value={form.monto}
          onChange={(e) => actualizarCampo('monto', e.target.value)}
          placeholder="0.00"
        />
        {montoInvalido && <span className="campo-error">El monto debe ser mayor a 0</span>}
      </label>

      <label>
        Categoría
        <select
          value={form.categoria}
          onChange={(e) => actualizarCampo('categoria', e.target.value)}
        >
          {CATEGORIAS.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
      </label>

      <label>
        Descripción
        <input
          type="text"
          value={form.descripcion}
          onChange={(e) => actualizarCampo('descripcion', e.target.value)}
          placeholder="Opcional"
        />
      </label>

      <label>
        Fecha
        <input
          type="date"
          value={form.fecha}
          max={hoy()}
          onChange={(e) => actualizarCampo('fecha', e.target.value)}
        />
        {fechaFutura && <span className="campo-error">La fecha no puede ser futura</span>}
      </label>

      <button type="submit" disabled={!puedeEnviar}>
        {enviando ? 'Guardando…' : 'Agregar gasto'}
      </button>

      {error && <p className="error">{error}</p>}
    </form>
  )
}
