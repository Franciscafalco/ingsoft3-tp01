## Por qué Git no pudo resolver el conflicto solo

Las ramas `feature/titulo-a` y `feature/titulo-b` modificaron la misma primera línea del README. Git puede fusionar automáticamente cuando los cambios están en lugares distintos del archivo, pero acá los dos tocaron exactamente la misma línea, elegir cuál vale es una decisión de contenido, no algo mecánico que Git pueda deducir. Por eso me dejó los marcadores y me pidió que decidiera yo.

Para que nunca hubiera aparecido, las dos ramas no tendrían que haber tocado la misma línea a la vez, si la rama B se hubiera creado después de mergear A, habría visto el cambio de A y no habría conflicto. En general, ramas cortas hacen que los conflictos sean chicos o que no existan.

## Problemas encontrados y cómo los resolví

- Escribí un mensaje de commit distinto al que quería. Lo corregí con `git commit --amend -m "..."`, que reescribe el mensaje del último commit. Aprendí que amend no edita el commit viejo sino que crea uno nuevo con otro hash.

- Intenté pushear a `main` y GitHub lo rechazó, porque ya la había protegido. El commit me había quedado parado en `main` local. Lo resolví moviéndolo a una rama con `git switch -c feature/url-repo-en-readme`, dejando `main` local igual al remoto con `git reset --hard origin/main`, y subiendo la rama para abrir el PR. El rechazo no fue un error: fue la protección funcionando y empujándome al flujo correcto (rama → PR).

## Declaración de uso de IA

Usé IA (Claude) para guiarme en el TP: entender el flujo de ramas y PRs, corregir el mensaje de un commit y resolver el rechazo del push a `main`.

Verifiqué cada indicación ejecutando yo misma los comandos y confirmando el resultado, revisé el estado con `git status` y `git log`, comprobé en GitHub que los PRs quedaran mergeados, que la rama protegida rechazara el push, y que la release v1.0.0 apareciera publicada. También contrasté los pasos con la guía de la cátedra.
