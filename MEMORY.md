# MEMORY.md - Git-To-Feed: Contexto, Arquitectura y Hoja de Ruta

## 1. Resumen y Propósito del Servicio
`git-to-feed` es un microservicio backend orientado a eventos escrito en Go diseñado para automatizar la presencia técnica en redes sociales profesionales (como LinkedIn) a partir de la actividad de desarrollo en repositorios de GitHub.

### Flujo Operativo Principal:
1. **Ingesta de Webhooks:** Escucha eventos de GitHub (`pull_request`, `release`, `issue_comment`). Valida firmas HMAC (SHA-256).
2. **Filtrado y Enriquecimiento de Contexto (Tri-Pillar Context):**
   - Extrae métricas de diff (`DiffStats`: archivos modificados, adiciones, supresiones).
   - Extrae notas humanas de los comentarios del PR (`/feed-notes`).
   - Inyecta la identidad y arquitectura del repositorio.
3. **Generación con IA (Gemini API):** Construye prompts dinámicos con ejemplos few-shot (almacenados en PostgreSQL) para redactar posts atractivos de ingeniería.
4. **Human-in-the-Loop Moderation (Discord Bot):** Envía tarjetas interactivas de aprobación a canales de Discord con botones de `Approve & Publish` o `Discard`.
5. **Publicación Externa:** Tras la aprobación, publica el post formateado a través de la API oficial de LinkedIn y almacena el registro en base de datos.

---

## 2. Hito Reciente Culminado: Ticket #62 / PR #79 (06/11/2026)

### Problema Resuelto:
Cuando se mergeaba un PR y poco después se cortaba un Release que contenía dicho PR, ambos eventos generaban borradores competidores en `approval_posts` en estado `PENDING`, produciendo ruido redundante en Discord y riesgo de publicaciones fragmentadas.

### Solución Implementada:
- **Jerarquía de Eventos:** Los eventos de tipo `RELEASE` tienen precedencia sobre PRs individuales.
- **Single Source of Truth:** Integración de la API de GitHub `CompareCommits(base, head)` mediante el puerto `ReleaseComparator`.
- **Transición Atómica en Base de Datos:** Los borradores de PR absorbidos pasan a `SUPERSEDED_BY_RELEASE` registrando el tag en la columna `superseded_by` (migración `000006_add_superseded_by_to_approval_posts.up.sql`).
- **Notificación en Discord:** Método `SendSupersededNotification` emite un aviso en Discord informando la invalidación del borrador a favor del release global.
- **Estado en Git:** PR #79 fusionado con estrategia `--merge` a `develop`. Suite de tests pasando al 100% con `-race` y 0 advertencias de `golangci-lint`.

---

## 3. Diagnóstico Técnico y Crisis de Autoría (Punto de Inflexión)

Durante la implementación del Ticket #62 se evidenció un problema crítico en el proceso de desarrollo:
- **La Ilusión de Competencia:** El proyecto creció aceleradamente adoptando patrones enterprise (Clean Architecture estricta, interfaces para cada interacción, mocks anidados, Docker cached runners, linter implacable).
- **La Brecha:** El autor dominaba el flujo conceptual del negocio, pero no los fundamentos primarios del lenguaje Go ni la navegación orgánica del código generado con asistencia de IA.
- **La Ley de Gall:** *"Un sistema complejo que funciona invariablemente ha evolucionado a partir de un sistema simple que también funcionaba."* `git-to-feed` intentó nacer en el "piso 10" sin que el autor dominara la mezcla del cemento.

---

## 4. Decisión Estratégica Acordada

En lugar de crear un repositorio nuevo desde cero o destruir el proyecto:

1. **Preservar el Repositorio Público:** Mantener el repo `hawksxo/git-to-feed`, su visibilidad, estrellas e historial de commits intacto en `develop`.
2. **Estrategia de Refactorización In-Place:**
   - Trabajar sobre ramas dedicadas de refactorización (ej. `refactor/package-by-feature-layout`).
   - Aplicar la filosofía aprendida en `lab-go` y Go Idiomático: código directo, simple, idiomático y 100% comprendido por el autor.
   - Aplanar y desacoplar abstracciones innecesarias capa por capa (empezando por el manejo HTTP y deserialización de webhooks) antes de introducir patrones avanzados.
3. **Rol de la IA en Adelante:** Guía socrática y explicativa; el autor es quien razona y escribe el código para consolidar la memoria muscular y el criterio profesional real.

---

## 5. Política Lingüística y de Gobernanza del Repositorio (08/11/2026)

Para eliminar la fricción cognitiva de la traducción forzada y garantizar que el autor comprenda y domine cada palabra del proyecto con su nivel de inglés actual, se establece la siguiente política por áreas:

1. **Gestión de Proyecto y Comunicación (100% en Español):**
   - **Issues y Milestones:** Títulos, descripciones, análisis y checklists escritos nativamente en español.
   - **Labels / Etiquetas:** Nombres y descripciones en español para navegación intuitiva del autor.
   - **Pull Requests y Revisiones:** Explicaciones de cambios, contexto y debates técnicos en español.
   - **Commits:** Conventional Commits en español (ej: `refactor: :truck: mover config y storage...`).
2. **Código Fuente (Inglés Técnico Idiomático):**
   - Nombres de paquetes, structs, interfaces, funciones y variables en inglés conciso de biblioteca estándar (`net/http`, `crypto/hmac`, `io`). Permite mantener compatibilidad con el ecosistema y la documentación oficial de Go.
3. **Colaboración Internacional Asistida:**
   - Para colaboradores externos que no hablen español, el workflow oficial `dromara/issues-translate-action@v2.8` traduce automáticamente issues y comentarios nuevos al inglés.
4. **Documentación Pública (Estrategia Bilingüe):**
   - Documentos de cara al público general (`README.md`, `CONTRIBUTING.md`, `ROADMAP_EVOLUTION.md`) adoptarán una estructura bilingüe (Español nativo como base + versión en Inglés para alcance global).

---

## 6. Culminación Milestone v2.0.0: Package-by-Feature (Issue #80 / PR #81 - 12/11/2026)

Se completó la migración integral desde Clean Architecture sobre-ingenierizada hacia la arquitectura idiomática **Package-by-Feature**, aplanando el microservicio en características verticales y altamente cohesivas:

### Estructura y Módulos de la v2.0.0:
1. **`internal/config` & `internal/storage`:**
   - Extraídos a la raíz de `internal/`.
   - Soporte directo de migraciones SQL para PostgreSQL/Supabase y ping de salud.
2. **`internal/webhook`:**
   - `parser.go`: Validación de firmas HMAC SHA-256 y parseo seguro de payloads.
   - `filter.go`: Evaluación pura de criterios de entrega (ramas, prefijos convencionales, tags).
   - `model.go`: Estructuras de datos del evento GitHub.
   - `github_client.go`: Cliente de la API de GitHub cohesionado directamente dentro del feature receptor para inspección de diffs y rate limits.
3. **`internal/generator`:**
   - Migración desde el antiguo `internal/pipeline` (reducción de 19 archivos a 5).
   - Integración con el SDK oficial de Google GenAI (`google.golang.org/genai`) y modelo `gemini-3.5-flash-lite`.
   - `prompt.go` y `cleaner.go`: Plantillas estáticas por arquetipo (`release`, `feature`, `refactor`) y filtros anti-cringe de texto.
4. **`internal/approval`:**
   - `service.go`: Ciclo de vida estricto de borradores (`PENDING`, `APPROVED`, `REJECTED`, `EXPIRED`, `SUPERSEDED_BY_RELEASE`) con TTL de 14 días.
   - `store.go`: Almacenamiento PostgreSQL (`PostgresStore`) con fallback `InMemoryStore` para tests y ejecución local.
   - `bot.go`: Cliente Discord WebSocket interactivo con botones de publicación y canal de auditoría. Eliminados los endpoints REST muertos de aprobación.
5. **`internal/publisher`:**
   - `client.go`: Cliente HTTP nativo para la API UGC Posts de LinkedIn con mock desacoplado.
   - `service.go`: Transición y auditoría de publicaciones (`PENDING`, `PUBLISHED`, `FAILED`).
   - `store.go`: Persistencia PostgreSQL (`published_posts`) y fallback `InMemoryStore`.
6. **Eliminación Total de `internal/platform/`:**
   - Removido el paquete genérico `platform` (HTTP, servicios y wrappers obsoletos).
7. **`cmd/api` (Unificación del Microservicio):**
   - `server.go`: Router único (`NewRouter`) que orquesta linealmente:
     `GitHub Webhook -> Generator -> Approval (Discord Bot) -> Publisher (LinkedIn)`
   - `main_test.go`: Suite completa E2E validando todo el flujo en memoria.
8. **Entrega Formal:**
   - Rama base de integración: **`next`** (baseline de la versión 2.0.0).
   - Pull Request: **[PR #81](https://github.com/hawksxo/git-to-feed/pull/81)** vinculado al Issue #80.
