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

# TP3 — Planificación y trazabilidad

## Duración del sprint

Elegí **1 semana**. Tenemos una clase por semana en la materia, así que alinear el sprint a esa cadencia hace que cada clase sea, de hecho, el cierre de un sprint: reviso qué quedó hecho, ajusto el plan, y arranco el siguiente. Un sprint más largo (2 semanas) no tendría un punto de revisión natural en el medio.

## Límite de trabajo en progreso

Elegí **2**, siguiendo la regla de arranque de la guía: cantidad de personas + 1. Trabajando sola, eso da 2. El "+1" es la válvula para cuando una tarea queda esperando algo (por ejemplo, una revisión) y necesito poder avanzar en otra sin quedarme bloqueada — pero sin que el límite deje de limitar. Si con el tiempo veo que nunca lo alcanzo, es señal de que está demasiado alto y lo bajaría.

## Diagnóstico de la historia mal escrita

La historia de ejemplo (_"Como desarrollador quiero crear la tabla usuarios para guardar los datos"_) está mal escrita porque es una **tarea disfrazada de historia**: nadie "quiere" una tabla en la base de datos — eso es un detalle de implementación, no un incremento de valor observable para alguien. Le falta el beneficio real (el "para qué" que le importa a un usuario, no a la base de datos) y no es _Valiosa_ en términos de INVEST.

Cómo la reescribiría: la historia real sería algo como _"Como usuario quiero que mis datos se guarden de forma persistente para poder acceder a mi cuenta la próxima vez que entre"_ — con criterios de aceptación verificables (por ejemplo: los datos siguen ahí después de cerrar sesión y volver a entrar). "Crear la tabla usuarios" pasaría a ser una **tarea** dentro de esa historia, no la historia en sí.

## Problemas encontrados y cómo los resolví

- **Mergeé el PR de la tarea #11 sin poner `Closes #11` en la descripción**, así que el issue no se cerró solo al mergear — me di cuenta después. Como el enlace automático solo se procesa en el momento del merge, no se puede agregar después. Lo resolví cerrando el #11 a mano (el trabajo ya estaba terminado) y usé la tarea #12, con un PR nuevo bien enlazado (`Closes #12` en la descripción), para tener la trazabilidad automática que pide el TP.

## Declaración de uso de IA

Usé IA (Claude) para guiarme paso a paso en la configuración de GitHub Projects, para pensar el diagnóstico de la historia mal escrita, y para diagnosticar y corregir el problema de trazabilidad del PR #14 (no cerraba el issue #11 por faltarle `Closes #11` en la descripción).

Las decisiones (duración del sprint, número de WIP) las tomé yo, no la IA. Verifiqué cada paso mirando yo misma el estado del proyecto en GitHub (issues cerrados, jerarquía de sub-issues, board con el sprint asignado) y confirmé que los issues #11 y #12 quedaron cerrados y que el PR #15 cerró el #12 automáticamente.

# TP4 — CI: Pipelines as Code

## Estructura del pipeline

Dos jobs, `build-backend` y `build-frontend`, uno por cada Dockerfile del TP2. Van **en paralelo** (no hay `needs:` entre ellos) porque no dependen uno del otro: cada uno construye su propia imagen, en su propia máquina, y no hay ninguna razón para esperar a que termine el otro. Si mi app tuviera un solo Dockerfile, sería un solo job — acá tiene sentido separarlos porque son dos artefactos independientes.

Disparadores: `pull_request` (el que hace el trabajo real, verifica antes del merge) y `push` a `main` (deja la corrida que lee el badge, y es la que le deja cache disponible a cualquier PR nuevo).

## Qué cachea y qué pasa si desaparece

Se cachean las **capas de Docker** de cada Dockerfile: en el backend, `COPY go.mod go.sum` y `RUN go mod download` (no vuelve a bajar dependencias si no cambiaron) y hasta el propio `RUN go build`; en el frontend, `RUN npm ci` y `RUN npm run build`. Cada job tiene su propio `scope` (`backend` / `frontend`) en `cache-from`/`cache-to`, así no se pisan entre sí — lo comprobé corriendo el pipeline dos veces seguidas (un commit vacío después de que la primera corrida terminara) y viendo `CACHED` en esas capas la segunda vez, en los dos jobs.

Si el cache desaparece (la plataforma lo puede desalojar en cualquier momento), el pipeline tiene que funcionar **igual**, solo que más lento: reconstruye todo desde cero. Si fallara sin cache, no era un cache — era una dependencia escondida que dependía de que algo ya estuviera ahí.

## Por qué construye con mi Dockerfile y no compila por su cuenta

Porque ya tengo **una** definición de cómo se compila mi app: el Dockerfile del TP2. Si el pipeline compilara aparte con `go build` directo, tendría dos definiciones de build que tarde o temprano divergen — y estaría verificando algo distinto de lo que después se despliega. Usar el mismo Dockerfile en CI y en despliegue es la garantía de que "lo que se verificó es lo que se corre".

## Problemas encontrados y cómo los resolví

- **El `ci.yml` fallaba en todos los runs, sin excepción.** La causa era un `:` (dos puntos) dentro de un string sin comillas en un paso `run:` (`"Reporte de tests: pendiente..."`) — YAML interpreta ese `:` como si abriera un nuevo par clave-valor, y rompe el archivo entero. Lo confirmé reproduciendo el error con un parser de YAML real (`yq`, corrido en un contenedor) contra la versión vieja del archivo, y until until until— lo arreglé pasando ese comando a un bloque multilínea (`run: |`), que no tiene esa restricción.
- **Al job `build-backend` se le había quedado sin el paso de `setup-buildx-action` y las líneas de cache**, mientras que `build-frontend` sí las tenía — un desprolijo al copiar la sección de la guía. Sin eso, ese job nunca iba a mostrar `CACHED`. Lo agregué con su propio `scope: backend`, distinto del `frontend`, para que no compartan estante.
- **Configuré por error un Ruleset de GitHub en vez de editar la regla clásica de `main`** (la del TP1), mientras buscaba activar los status checks obligatorios — la propia guía del TP1 avisa de no usar "Go to rulesets". Lo detecté porque quedaron las dos protecciones superpuestas; borré el Ruleset nuevo y dejé la configuración de required status checks en la regla clásica de siempre.
- **`gh` y `bat` no se encontraban en la terminal después de instalarlos con `winget`.** Los dos estaban instalados correctamente en disco; el problema era que la terminal ya abierta tenía el PATH viejo en memoria. Se resolvió cerrando y abriendo una terminal nueva.

## Declaración de uso de IA

Usé IA (Claude) para diagnosticar el error de YAML (verificado corriendo `yq` contra el archivo real, no solo mirándolo), detectar que al job del backend le faltaban pasos comparándolo línea por línea con el del frontend, y para armar este mismo archivo.

Verifiqué todo yo misma: corrí `docker build ./backend` en mi máquina para confirmar que el import roto realmente rompía la compilación antes de subirlo, miré las corridas de Actions con mis propios ojos (los checks en rojo, después en verde, `CACHED` en los logs de las dos etapas), y confirmé con GitHub que no quedó ningún PR abierto y que el badge del README apunta bien.
