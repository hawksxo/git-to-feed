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

## 6. Registro de Deuda Técnica y Correcciones Arquitectónicas (10/11/2026)

A partir del análisis de Go Idiomático (principios SOLID aplicados a Go) y la verdadera naturaleza del proyecto como **Microservicio Orientado a Eventos**, se registran las siguientes simplificaciones esenciales:

1. **La Verdadera Identidad del Proyecto:**
   - `git-to-feed` no es una API REST con dashboard de usuario; es un **procesador reactivo de eventos**.
   - El servidor HTTP existe **únicamente** para recibir los Webhooks que GitHub empuja y para los checks de salud (`/health`).
2. **Eliminación de Endpoints Muertos / Fantasma (Complejidad Accidental):**
   - Los endpoints REST de `approval` (`GET /api/v1/approvals/pending`, `POST /{uuid}/reject`, etc.) son código muerto: la aprobación y descarte ocurren en vivo mediante **Discord (WebSocket)**. Deben eliminarse o relegarse si no tienen un frontend consumidor.
3. **Aplicación Pragmática de SOLID en Go:**
   - **S (Responsabilidad Única):** Separar el cómputo puro (`filter.go`, `cleaner.go`, validación HMAC) de la red y el disco.
   - **I y D (Interfaces Segregadas e Invertidas):**
     - La interfaz la define el **paquete que consume**, no el que implementa.
     - Interfaces atómicas de 1 o 2 métodos máximo (`io.Reader` style).
     - Eliminar interfaces prematuras donde solo existe un único struct concreto (regla: *"Accept interfaces, return structs"*).
4. **Simplificación del Paquete `generator` (en progreso):**
   - Consolidar los 19 archivos dispersos de `pipeline` en 5 responsabilidades transparentes: `model.go`, `cleaner.go`, `client.go`, `prompt.go` y `generator.go`.
