# 🚀 Architecture & Feature Evolution Roadmap (Git-To-Feed)

This document records the architectural roadmap, foundational platform modules, and high-impact features planned for the technical content generation pipeline.

---

# 🧱 PART I: FOUNDATIONAL PLATFORM MODULES (BASE)
> **Active Target Milestone:** [v1.1.0 - Foundation & Infrastructure Modules](https://github.com/hawksxo/git-to-feed/milestone/1)

## 1. ⚙️ Automated CI/CD Pipelines ([Issue #57](https://github.com/hawksxo/git-to-feed/issues/57)) - COMPLETED ✅
- **Scenario:** Manual verification tests and container-coupled database migrations introduce deployment race conditions and lack automated quality barriers.
- **Current State:** Implemented via GitHub Actions (`ci.yml` and `cd.yml`) with branch protection gates on `develop` (Completed in PR #67).
- **Architectural Invariant:** Strict separation of concerns. Decouple database migrations from container startup; enforce fail-fast automated testing before merging.
- **Deliverables:**
  - **CI Workflow (`.github/workflows/ci.yml`)**: Triggered on PRs to `develop`/`release`. Runs `go test -race ./...`, `golangci-lint`, and `govulncheck`. Blocks merge on failure.
  - **CD Workflow (`.github/workflows/cd.yml`)**: Triggered on release tag push (`v*.*.*`). Applies `supabase db push` via Supabase CLI, then invokes Render Deploy Hook only upon successful migration.

---

## 2. 📢 Multi-Channel Routing & Graceful Discord Fallback ([Issue #58](https://github.com/hawksxo/git-to-feed/issues/58)) - COMPLETED ✅
- **Scenario:** Single notification channel produces visual clutter, mixing major release drafts, minor PR drafts, and publication audit logs.
- **Current State:** Implemented multi-channel routing with fallback resolution in `internal/platform/discord` and `internal/platform/config` (Completed in PR #68).
- **Architectural Invariant:** Decouple category routing into application use cases; enforce backward-compatible graceful fallback if specific channel IDs are omitted.
- **Deliverables:**
  - Route release cards to `DISCORD_CHANNEL_RELEASES_ID` (`#drafts-releases`).
  - Route PR feature drafts to `DISCORD_CHANNEL_FEATURES_ID` (`#drafts-features`).
  - Route LinkedIn publication confirmations to `DISCORD_CHANNEL_AUDIT_ID` (`#published-feed`).
  - Dedicated operations channel (`#bot-control`) for administrative actions.

---

## 3. 🧹 Interactive Draft Lifecycle & Zero-Cost Housekeeping ([Issue #59](https://github.com/hawksxo/git-to-feed/issues/59)) - COMPLETED ✅
- **Scenario:** Unapproved drafts accumulate indefinitely in `PENDING` state in Supabase, cluttering database storage and queries with stale records.
- **Current State:** Implemented interactive one-click discard button in Discord with atomic optimistic concurrency guards and 14-day Lazy TTL boundary in PostgreSQL (Completed in PR #69).
- **Architectural Invariant:** 1-click frictionless developer experience; zero persistent background workers (zero CPU/memory bloat on free tiers).
- **Deliverables:**
  - **One-Click Discard Action (`🔴 Discard`)**: Discord button that atomically transitions post status to `REJECTED` (`reject_reason: manual_discard`) and deactivates the Discord embed.
  - **Administrative Slash Query (`/posts list [status]`)**: Interactive operations command in `#bot-control` with ephemeral responses.
  - **Lazy TTL Boundary**: Read-time exclusion or native Supabase `pg_cron` boundary marking drafts older than 14 days as `EXPIRED`.

---

## 4. 🐙 GitHub Platform Adapter & REST Client Suite ([Issue #60](https://github.com/hawksxo/git-to-feed/issues/60)) - COMPLETED ✅
- **Scenario:** Multiple downstream features (Compare API, diff stats, rate limit probes) require authenticated, resilient communication with GitHub REST API.
- **Current State:** Implemented dedicated platform adapter in `internal/platform/github` with rate limit inspection, PR file diffing, and commit range comparison (Completed in PR #70).
- **Architectural Invariant:** Encapsulate GitHub HTTP communication behind idiomatic Go interfaces with timeout boundaries and rate-limit tracking.
- **Deliverables:**
  - Package `internal/platform/github` with authenticated client (`GITHUB_TOKEN`).
  - `CompareCommits(ctx, owner, repo, base, head)`: Retrieves commit range diffs and associated PR numbers.
  - `GetPullRequestFiles(ctx, owner, repo, prNumber)`: Extracts touched files, layers, and addition/deletion line metrics.
  - `GetRateLimitStatus(ctx)`: Inspects remaining API quota (`X-RateLimit-Remaining`).

---

## 5. 🧠 Dynamic Few-Shot Prompt Builder & Supabase SSOT ([Issue #56](https://github.com/hawksxo/git-to-feed/issues/56))
- **Scenario:** `PromptBuilder` uses hardcoded few-shot examples inside binary code, bypassing the existing `few_shot_examples` database table and requiring full binary redeployments to update prompt engineering.
- **Current State:** Literal examples reside in `internal/pipeline/template.go`; `PostgresFewShotRepository` is not injected into the pipeline generation layer.
- **Architectural Invariant:** Single Source of Truth (SSOT) via database-driven prompt configuration with graceful in-memory static fallback for resilience against database downtime.
- **Deliverables:**
  - Database seed migration (`supabase/migrations/`) populating default few-shot archetypes.
  - Dependency injection of `FewShotRepository` into `PromptBuilder`.
  - Dynamic active archetype lookup (`FindActiveByArchetype`) with fail-safe fallback to static instructions.
  - Unit tests covering dynamic generation and database failure fallback boundaries.

---

# 🚀 PART II: FEATURE & INTELLIGENCE CAPABILITIES
> **Target Status:** Backlog (Scheduled for future minor milestones)

## 6. 🎯 Context Quality & Tri-Pillar Enrichment ([Issue #61](https://github.com/hawksxo/git-to-feed/issues/61))
- **Scenario:** Pull Requests or Releases with brief descriptions restrict Gemini AI output fidelity ("Garbage In, Garbage Out").
- **Current State:** Generation depends strictly on PR template markdown; `HumanNotes` domain entity is disconnected from external inputs.
- **Architectural Invariant:** Ground generation in real structural changes, developer guidance, and persistent repository identity.
- **Deliverables:**
  - **Pillar 1 (Diff Stats)**: Automated extraction of modified packages and diff stats via GitHub Client Adapter.
  - **Pillar 2 (PR Notes)**: Ingestion of `/feed-notes <text>` via GitHub `issue_comment` webhooks.
  - **Pillar 3 (Project Identity)**: Ingestion of a versioned `ARCHITECTURE.md` / `CONTEXT.md` document into base system instructions.

---

## 7. 👑 Event Hierarchy Resolution ([Issue #62](https://github.com/hawksxo/git-to-feed/issues/62))
- **Scenario:** A PR draft is generated in Discord, and shortly after a `Release Tag` is published containing the changes from that PR.
- **Current State:** Both events produce competing drafts in `approval_posts` (state `PENDING`).
- **Architectural Invariant:** GitHub Compare API as Single Source of Truth; atomic status invalidation.
- **Deliverables:**
  - Orchestrator queries GitHub `CompareCommits` when receiving a `RELEASE` event.
  - Matching pending PR drafts in `approval_posts` are transitioned atomically to `SUPERSEDED_BY_RELEASE`.
  - Discord PR embeds are updated with an informative badge advising publication of the parent release.

---

## 8. 🩺 Extended Health Probes & Dependency Observability ([Issue #63](https://github.com/hawksxo/git-to-feed/issues/63))
- **Scenario:** Silent failures occur when external dependencies (GitHub token, Gemini API key, Discord channels) degrade.
- **Current State:** Active probes limited to database connection and Discord Gateway WebSocket (`/health/discord-bot`).
- **Architectural Invariant:** Non-intrusive synthetic sub-endpoints reporting upstream health without leaking credentials.
- **Deliverables:**
  - `/health/github-api`: Validates token validity and remaining API budget.
  - `/health/gemini`: Validates Gemini API key and model readiness.
  - Aggregate status reporting within master `/health` endpoint.

---

## 9. 📜 Consolidated Version Digest & Adaptive Synthesis Engine ([Issue #64](https://github.com/hawksxo/git-to-feed/issues/64))
- **Scenario:** Major and minor releases lose depth when iterative patch histories (`v1.0.1`...`v1.0.7`) are fragmented across individual drafts.
- **Current State:** Version drafts only reflect immediate tag release notes without historical awareness of predecessor patches.
- **Architectural Invariant:** Coordinated triggers (automated milestone vs. manual on-demand) backed by database idempotency and reusable GitHub Compare infrastructure.
- **Deliverables:**
  - **Automated Milestone Synthesis**: Compiles all unmerged patch notes into a cohesive evolutionary post upon `vX.Y.0` releases.
  - **On-Demand Slash Command (`/digest base:<tag> head:<tag>`)**: Interactive generation in `#bot-control`.
  - **Narrative Prompting**: Thematic extraction of architecture pillars over raw commit logs.

---

## 10. 🔄 24h Daily Cadence Rule & Staging Event Queue ([Issue #65](https://github.com/hawksxo/git-to-feed/issues/65))
- **Scenario:** Multiple PRs and rapid patches published on the same day flood LinkedIn feeds with low-value atomic posts.
- **Current State:** System reactively generates an isolated draft per event.
- **Architectural Invariant:** Maximum one single publication per 24-hour cycle.
- **Deliverables:**
  - Staging table in PostgreSQL (`event_staging`) accumulating daily merged events.
  - 24-hour scheduled cadence worker or trigger synthesizing daily staged events into a single consolidated post.

---

## 11. 🛡️ External API Resilience, Circuit Breaker & Retry Mechanism ([Issue #66](https://github.com/hawksxo/git-to-feed/issues/66))
- **Scenario:** Transient external API rate limits (LinkedIn `HTTP 429`, Gemini timeouts) cause permanent pipeline drops.
- **Current State:** Errors result in immediate rejection without deferred retries.
- **Architectural Invariant:** Non-blocking exponential backoff with jitter and explicit recoverable error states.
- **Deliverables:**
  - Retry policy wrapper around LinkedIn and Gemini HTTP adapters.
  - `FAILED_RETRYABLE` post status and automated alert dispatching to `#bot-control`.