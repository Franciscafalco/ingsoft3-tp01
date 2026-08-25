# TP1 — Git colaborativo

## Por qué Git no pudo resolver el conflicto solo

Las ramas `feature/titulo-a` y `feature/titulo-b` modificaron la misma primera línea del README. Git puede fusionar automáticamente cuando los cambios están en lugares distintos del archivo, pero acá los dos tocaron exactamente la misma línea, elegir cuál vale es una decisión de contenido, no algo mecánico que Git pueda deducir. Por eso me dejó los marcadores y me pidió que decidiera yo.

Para que nunca hubiera aparecido, las dos ramas no tendrían que haber tocado la misma línea a la vez, si la rama B se hubiera creado después de mergear A, habría visto el cambio de A y no habría conflicto. En general, ramas cortas hacen que los conflictos sean chicos o que no existan.

## Problemas encontrados y cómo los resolví

- Escribí un mensaje de commit distinto al que quería. Lo corregí con `git commit --amend -m "..."`, que reescribe el mensaje del último commit. Aprendí que amend no edita el commit viejo sino que crea uno nuevo con otro hash.

- Intenté pushear a `main` y GitHub lo rechazó, porque ya la había protegido. El commit me había quedado parado en `main` local. Lo resolví moviéndolo a una rama con `git switch -c feature/url-repo-en-readme`, dejando `main` local igual al remoto con `git reset --hard origin/main`, y subiendo la rama para abrir el PR. El rechazo no fue un error: fue la protección funcionando y empujándome al flujo correcto (rama → PR).

## Declaración de uso de IA

Usé IA (Claude) para guiarme en el TP: entender el flujo de ramas y PRs, corregir el mensaje de un commit y resolver el rechazo del push a `main`.

Verifiqué cada indicación ejecutando yo misma los comandos y confirmando el resultado, revisé el estado con `git status` y `git log`, comprobé en GitHub que los PRs quedaran mergeados, que la rama protegida rechazara el push, y que la release v1.0.0 apareciera publicada. También contrasté los pasos con la guía de la cátedra.

# TP2 — Contenedores

## Qué app elegí y por qué

**Gestor de gastos personales**: backend en Go (Gin + GORM), frontend en React + Vite, base de datos PostgreSQL. Es una app nueva, chica, escrita para este TP y pensada para acompañarme el resto del semestre.

La contrasté contra los 5 criterios de `elegir-app.md`:

1. **¿Corre hoy?** Sí: la probé end-to-end en mi máquina (backend con `go run`, frontend con `npm run dev`, Postgres en un contenedor aparte) antes de escribir un solo Dockerfile.
2. **¿Conozco los comandos de compilación?** Sí: `go build` para el backend, `npm run build` para el frontend (Vite emite a `dist/`).
3. **¿Dónde se configura la conexión a la base?** Por variables de entorno (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`) leídas en `backend/internal/db/db.go`, con valores por defecto solo para desarrollo. Nada de la conexión está fijo en el código.
4. **¿Tiene reglas de negocio para testear?** Se la diseñé a propósito con el TP5 en mente (8 tests de backend, 4 de frontend). Reglas del backend (`backend/internal/models/gasto.go`):
   - Validación: el monto debe ser mayor a 0.
   - Validación: la fecha no puede ser futura.
   - Validación: la categoría debe pertenecer a una lista fija.
   - Cálculo: total general y total por categoría (`backend/internal/service/resumen.go`).
   - Transición de estado: `pendiente → pagado` permitida, `pagado → pendiente` prohibida.
   - Restricción: no se puede eliminar un gasto en estado `pagado`.

   Comportamientos del frontend:
   - El formulario no deja enviar con monto ≤ 0 o fecha futura.
   - El total se recalcula solo al agregar/eliminar/cambiar el estado de un gasto.
   - El botón "Eliminar" se deshabilita si el gasto está `pagado`.

5. **¿La entiendo lo suficiente para modificarla?** Sí: la escribí yo (con ayuda de IA, ver más abajo) y puedo señalar en qué archivo va cada regla.

Tamaño acotado a propósito: dos pantallas (listado con formulario, y resumen por categoría), sin dependencias además de Postgres.

## Decisiones de contenerización

- **Backend**: Dockerfile multi-stage. Etapa de build sobre `golang:1.25-alpine` (compila con
  `CGO_ENABLED=0` para un binario estático), etapa final sobre `alpine:3.20` que solo copia el
  binario. Imagen final: **68.4 MB** contra **329 MB** de la imagen de build (con el SDK completo).
- **Frontend**: Dockerfile multi-stage. Etapa de build sobre `node:22-alpine` (`npm ci` + `npm run
build`), etapa final sobre `nginx:alpine` que sirve los estáticos y hace de proxy: todo lo que
  llega a `/api/` lo reenvía a `http://backend:8080` dentro de la red de compose. Así el frontend
  llama a rutas relativas (`/api/...`) y no hace falta configurar CORS.
- **Qué persiste y qué no**: el volumen nombrado `db_data` (montado en
  `/var/lib/postgresql/data`) es lo único que sobrevive a `docker compose down`. Los contenedores en
  sí son descartables: los recreé varias veces sin perder los gastos cargados, y con
  `docker compose down -v` verifiqué que el volumen (y los datos) se borran.
- **Red y descubrimiento**: los servicios se hablan por nombre (`db`, `backend`) gracias al DNS
  interno de compose; el `backend` usa `Host=db` en vez de una IP fija.
- **`depends_on` + `healthcheck`**: el backend espera a que Postgres esté _healthy_
  (`pg_isready`), no solo a que el contenedor haya arrancado — sin esto, el backend intentaría
  conectarse antes de que Postgres acepte conexiones.
- **Secretos**: la contraseña de la base viaja en `DB_PASSWORD`, cargada desde un `.env` que está en
  `.gitignore`. Lo que se commitea es `.env.example`, con un valor de ejemplo.
- **Puertos publicados**: `frontend` en `3000:80` y `backend` en `8081:8080` (para poder pegarle con
  curl/Postman directo). Adentro de la red de compose el backend sigue siendo `backend:8080`
  siempre: lo único que cambia es el puerto que ve mi máquina desde afuera.

## Problemas encontrados y cómo los resolví

- **Gin necesita Go 1.25, y yo tenía Go 1.24.1 instalado.** `go get` fijó `go.mod` en `go 1.25.0` porque `gin-gonic/gin v1.12.0` lo exige. Probé bajarlo a `1.24` para no depender de un toolchain nuevo, pero el build falló (`requires go >= 1.25.0`). Lo resolví usando `golang:1.25-alpine` como base de la etapa de build del Dockerfile — no hace falta que el Go de mi máquina y el de la imagen coincidan, y confirmé que esa tag existe en Docker Hub antes de comprometerme.
- **El puerto 8080 de mi máquina ya estaba ocupado** por un proceso `httpd` (un Apache/XAMPP instalado aparte, sin relación con este proyecto). En vez de matar un proceso que no era mío, cambié el mapeo del backend en `docker-compose.yml` a `8081:8080`: el puerto publicado hacia el host es una decisión mía, no afecta la comunicación interna entre contenedores.
- **`decisiones.md` y `evidencias.md` del TP1 habían quedado sin subir.** Existían en mi carpeta local, pero nunca habían llegado a `main` por PR. Los subí en un PR aparte, previo al de este TP2, para no mezclar la entrega del TP1 con el código nuevo. De paso, encontré que mi `main` local tenía dos commits que no existían en el remoto (un commit hecho sin querer directo en local, más un merge): comparé el contenido con `git diff origin/main main` y, al ser idéntico, hice `git reset --hard origin/main` para dejar la rama local prolija sin perder nada. Esto se repitió una segunda vez más adelante (mismo patrón, mismo diagnóstico y misma solución).
- **OneDrive rompía Go y Docker de forma intermitente.** Mi carpeta de trabajo vivía dentro de `OneDrive\Escritorio\...`, y tanto `go build`/`go mod tidy` como `docker build` fallaban sin patrón fijo: a veces "cannot find package" para dependencias que estaban bien instaladas, a veces `docker build` directamente no podía leer el `Dockerfile` (`invalid file request`). Cerrar OneDrive ayudó para Go pero no alcanzó para Docker. Confirmé la causa copiando el proyecto a una carpeta fuera de OneDrive: ahí compiló a la primera, tanto con Go como con Docker. La solución de fondo fue mover la carpeta de trabajo a `C:\dev\ingsoft3-tp01` — el repositorio remoto en GitHub no cambia, solo dónde vive mi copia local.
- **Un `go mod tidy` corrido mientras Go ya estaba fallando por lo anterior vació `go.mod`/`go.sum`** en vez de mostrar un error claro (`warning: "all" matched no packages"`). Lo resolví con `git restore go.mod go.sum`, porque el contenido correcto ya estaba commiteado.
- **Al verificar que las imágenes de ghcr.io quedaran públicas, tipeé mal el logout** (`docker logout ghrc.io` en vez de `ghcr.io`). Como no cerró la sesión real, el primer intento de "bajar la imagen sin sesión" no probaba nada — seguía logueada. Repetí el logout bien escrito y ahí sí confirmé que el `pull` funcionaba sin credenciales.

## Publicación en el registry

Imágenes `gastos-backend` y `gastos-frontend` publicadas en `ghcr.io/franciscafalco/`, tag `v0.1.0`,
visibilidad **Public** (verificado bajándolas sin sesión). Las construí en una PC con arquitectura
Intel/AMD (amd64); alguien con otra arquitectura (por ejemplo un Mac con chip ARM) necesitaría
re-buildear en vez de usar estas imágenes tal cual — en el TP7 esto se resuelve con `docker buildx`.

## Declaración de uso de IA

Usé IA (Claude) para prácticamente todo el desarrollo de este TP: definir el alcance de la app (incluyendo las reglas de negocio, pensadas junto con la IA para cubrir los requisitos del TP5), escribir el código del backend y del frontend, los Dockerfiles multi-stage, el `docker-compose.yml`/`docker-compose.registry.yml`, el `README.md`, y diagnosticar el problema de OneDrive con Go y Docker (incluida la prueba de copiar el proyecto a otra carpeta para confirmar la causa antes de mover nada).

Cada decisión importante (qué app hacer, con qué stack, qué reglas de negocio, cómo organizar las ramas y PRs, mover la carpeta de trabajo) la charlé y aprobé antes de que se implementara, no dejé que la IA decidiera sola.
Verifiqué el resultado corriendo yo misma, en la misma sesión, todo el sistema: probé cada regla de negocio con `curl` (monto inválido, fecha futura, categoría inválida, transición de estado prohibida, restricción de borrado), usé la app en el navegador, hice la prueba de persistencia completa (`down` conserva los datos, `down -v` los borra), repetí los pasos de `README.md` desde una corrida limpia (`docker compose down -v` + `.env` nuevo) para confirmar que el arranque documentado funciona tal cual está escrito, y comprobé con mis propias manos (`docker login`, `docker push`, `docker logout`, `docker pull`) que las imágenes publicadas en ghcr.io son realmente públicas.
