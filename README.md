# Proyecto IngSoft3 - versión A

Repositorio para los TPs de la materia Ingeniería de Software III.

## App del semestre: Gestor de gastos personales

CRUD de gastos personales con reglas de negocio (validaciones, cálculo de totales, transición de
estado y restricción de borrado). Ver el detalle de las reglas en [`decisiones.md`](decisiones.md).

**Stack:**
- Backend: Go + Gin + GORM
- Frontend: React + Vite, servido por nginx en producción
- Base de datos: PostgreSQL

## Instalación

```bash
git clone <url-del-repo>
```

## Levantar el sistema completo (con Docker)

Prerequisito: [Docker](https://docs.docker.com/get-docker/) instalado y corriendo.

```bash
cp .env.example .env
# (opcional) editá .env y poné la contraseña que quieras para la base
docker compose up -d --build
```

Esto levanta tres contenedores: `db` (PostgreSQL), `backend` (API en Go) y `frontend` (React
servido por nginx, con proxy hacia el backend en `/api`).

- Frontend: http://localhost:3000
- Backend (API directa, para curl/Postman): http://localhost:8081
  > El backend se publica en el puerto **8081** del host (no 8080) porque en esta máquina el 8080
  > ya lo usa otro proceso local (un httpd). Dentro de la red de Docker el backend sigue
  > escuchando en el 8080 de siempre; solo cambia el puerto publicado hacia afuera.

Verificar que levantó bien:

```bash
curl -s http://localhost:8081/health
# {"status":"ok"}
```

Para bajar el sistema:

```bash
docker compose down       # apaga los contenedores, conserva los datos de la BD
docker compose down -v    # apaga y además borra el volumen de datos
```

## Levantar usando las imágenes publicadas (sin buildear localmente)

```bash
cp .env.example .env
docker compose -f docker-compose.registry.yml up -d
```

Baja las imágenes ya publicadas en ghcr.io en vez de construirlas desde el código.

## Desarrollo local (sin Docker)

**Backend** (necesita un PostgreSQL accesible, por ejemplo levantado con Docker):

```bash
cd backend
DB_HOST=localhost DB_PASSWORD=postgres go run .
```

**Frontend** (el servidor de Vite proxea `/api` hacia `http://localhost:8080`, ver
`vite.config.js`):

```bash
cd frontend
npm install
npm run dev
```

## Endpoints de la API

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/health` | Healthcheck |
| GET | `/api/gastos` | Lista los gastos (`?categoria=` para filtrar) |
| POST | `/api/gastos` | Crea un gasto |
| PUT | `/api/gastos/:id` | Edita un gasto (incluye cambio de estado) |
| DELETE | `/api/gastos/:id` | Elimina un gasto (no permitido si está `pagado`) |
| GET | `/api/gastos/resumen` | Total general y total por categoría |

