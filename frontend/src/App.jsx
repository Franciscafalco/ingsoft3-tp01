import { useEffect, useMemo, useState } from 'react'
import GastoForm from './components/GastoForm'
import GastosTable from './components/GastosTable'
import Resumen from './components/Resumen'
import { actualizarGasto, crearGasto, eliminarGasto, listarGastos, obtenerResumen } from './api'
import './App.css'

const formatoMoneda = new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' })

export default function App() {
  const [gastos, setGastos] = useState([])
  const [resumen, setResumen] = useState(null)
  const [vista, setVista] = useState('gastos')
  const [error, setError] = useState(null)

  const totalGeneral = useMemo(
    () => gastos.reduce((acc, g) => acc + g.monto, 0),
    [gastos],
  )

  async function cargarGastos() {
    try {
      setGastos(await listarGastos())
    } catch (err) {
      setError(err.message)
    }
  }

  async function cargarResumen() {
    try {
      setResumen(await obtenerResumen())
    } catch (err) {
      setError(err.message)
    }
  }

  useEffect(() => {
    cargarGastos()
  }, [])

  useEffect(() => {
    if (vista === 'resumen') cargarResumen()
  }, [vista])

  async function handleCrear(gasto) {
    setError(null)
    await crearGasto(gasto)
    await cargarGastos()
    if (vista === 'resumen') await cargarResumen()
  }

  async function handleMarcarPagado(gasto) {
    setError(null)
    try {
      await actualizarGasto(gasto.id, { ...gasto, estado: 'pagado' })
      await cargarGastos()
      if (vista === 'resumen') await cargarResumen()
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleEliminar(gasto) {
    setError(null)
    try {
      await eliminarGasto(gasto.id)
      await cargarGastos()
      if (vista === 'resumen') await cargarResumen()
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div className="app">
      <header>
        <h1>Gestor de gastos personales</h1>
        <nav className="tabs">
          <button
            type="button"
            className={vista === 'gastos' ? 'activo' : ''}
            onClick={() => setVista('gastos')}
          >
            Gastos
          </button>
          <button
            type="button"
            className={vista === 'resumen' ? 'activo' : ''}
            onClick={() => setVista('resumen')}
          >
            Resumen
          </button>
        </nav>
      </header>

      {error && <p className="error error-global">{error}</p>}

      {vista === 'gastos' ? (
        <main className="contenido-gastos">
          <GastoForm onCrear={handleCrear} />
          <section>
            <p className="total-general">Total: {formatoMoneda.format(totalGeneral)}</p>
            <GastosTable
              gastos={gastos}
              onMarcarPagado={handleMarcarPagado}
              onEliminar={handleEliminar}
            />
          </section>
        </main>
      ) : (
        <main>
          <Resumen resumen={resumen} />
        </main>
      )}
    </div>
  )
}
