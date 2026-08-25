const BASE_URL = '/api/gastos'

async function handleResponse(res) {
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `Error ${res.status}`)
  }
  if (res.status === 204) return null
  return res.json()
}

export function listarGastos() {
  return fetch(BASE_URL).then(handleResponse)
}

export function obtenerResumen() {
  return fetch(`${BASE_URL}/resumen`).then(handleResponse)
}

export function crearGasto(gasto) {
  return fetch(BASE_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(gasto),
  }).then(handleResponse)
}

export function actualizarGasto(id, gasto) {
  return fetch(`${BASE_URL}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(gasto),
  }).then(handleResponse)
}

export function eliminarGasto(id) {
  return fetch(`${BASE_URL}/${id}`, { method: 'DELETE' }).then(handleResponse)
}
