# 🚀 Architecture & Feature Evolution Roadmap (Git-To-Feed)

This document records the architectural roadmap, foundational platform modules, and high-impact features planned for the technical content generation pipeline.

---

# 🧱 PART I: FOUNDATIONAL PLATFORM MODULES (BASE)

## 1. ⚙️ Automated CI/CD Pipelines (Quality Gates, Supabase CD & Render Deploy)
- **Scenario:** Manual verification tests and container-coupled database migrations introduce deployment race conditions and lack automated quality barriers.
- **Current State:** Manual release branch workflow; zero automated PR branch protection gates.
- **Architectural Invariant:** Strict separation of concerns. Decouple database migrations from container startup; enforce fail-fast automated testing before merging.
- **Deliverables:**
  - **CI Workflow (`.github/workflows/ci.yml`)**: Triggered on PRs to `develop`/`release`. Runs `go test -race ./...`, `golangci-lint`, and `govulncheck`. Blocks merge on failure.
  - **CD Workflow (`.github/workflows/cd.yml`)**: Triggered on release tag push (`v*.*.*`). Applies `supabase db push` via Supabase CLI, then invokes Render Deploy Hook only upon successful migration.

---

## 2. 📢 Multi-Channel Routing & Graceful Discord Fallback
- **Scenario:** Single notification channel produces visual clutter, mixing major release drafts, minor PR drafts, and publication audit logs.
- **Current State:** All notifications target a single channel (`DISCORD_CHANNEL_ID`).
- **Architectural Invariant:** Decouple category routing into application use cases; enforce backward-compatible graceful fallback if specific channel IDs are omitted.
- **Deliverables:**
  - Route release cards to `DISCORD_CHANNEL_RELEASES_ID` (`#drafts-releases`).
  - Route PR feature drafts to `DISCORD_CHANNEL_FEATURES_ID` (`#drafts-features`).
  - Route LinkedIn publication confirmations to `DISCORD_CHANNEL_AUDIT_ID` (`#published-feed`).
  - Dedicated operations channel (`#bot-control`) for administrative actions.

---

## 3. 🧹 Interactive Draft Lifecycle & Zero-Cost Housekeeping
- **Scenario:** Unapproved drafts accumulate indefinitely in `PENDING` state in Supabase, cluttering database storage and queries with stale records.
- **Current State:** Approval cards lack an immediate discard action in Discord; pending drafts have no expiration boundary.
- **Architectural Invariant:** 1-click frictionless developer experience; zero persistent background workers (zero CPU/memory bloat on free tiers).
- **Deliverables:**
  - **One-Click Discard Action (`🔴 Discard`)**: Discord button that atomically transitions post status to `REJECTED` (`reject_reason: manual_discard`) and deactivates the Discord embed.
  - **Administrative Slash Query (`/posts list [status]`)**: Interactive operations command in `#bot-control` with ephemeral responses.
  - **Lazy TTL Boundary**: Read-time exclusion or native Supabase `pg_cron` boundary marking drafts older than 14 days as `EXPIRED`.

---

## 4. 🐙 GitHub Platform Adapter & REST Client Suite (`internal/platform/github`)
- **Scenario:** Multiple downstream features (Compare API, diff stats, rate limit probes) require authenticated, resilient communication with GitHub REST API.
- **Current State:** No dedicated GitHub API client exists; system only parses inbound webhook JSON payloads.
- **Architectural Invariant:** Encapsulate GitHub HTTP communication behind idiomatic Go interfaces with timeout boundaries and rate-limit tracking.
- **Deliverables:**
  - Package `internal/platform/github` with authenticated client (`GITHUB_TOKEN`).
  - `CompareCommits(ctx, owner, repo, base, head)`: Retrieves commit range diffs and associated PR numbers.
  - `GetPullRequestFiles(ctx, owner, repo, prNumber)`: Extracts touched files, layers, and addition/deletion line metrics.
  - `GetRateLimitStatus(ctx)`: Inspects remaining API quota (`X-RateLimit-Remaining`).

---

# 🚀 PART II: FEATURE & INTELLIGENCE CAPABILITIES

## 5. 🎯 Context Quality & Tri-Pillar Enrichment
- **Scenario:** Pull Requests or Releases with brief descriptions restrict Gemini AI output fidelity ("Garbage In, Garbage Out").
- **Current State:** Generation depends strictly on PR template markdown; `HumanNotes` domain entity is disconnected from external inputs.
- **Architectural Invariant:** Ground generation in real structural changes, developer guidance, and persistent repository identity.
- **Deliverables:**
  - **Pillar 1 (Diff Stats)**: Automated extraction of modified packages and diff stats via GitHub Client Adapter.
  - **Pillar 2 (PR Notes)**: Ingestion of `/feed-notes <text>` via GitHub `issue_comment` webhooks.
  - **Pillar 3 (Project Identity)**: Ingestion of a versioned `ARCHITECTURE.md` / `CONTEXT.md` document into base system instructions.

---

## 6. 👑 Event Hierarchy Resolution (Release vs. Pull Request)
- **Scenario:** A PR draft is generated in Discord, and shortly after a `Release Tag` is published containing the changes from that PR.
- **Current State:** Both events produce competing drafts in `approval_posts` (state `PENDING`).
- **Architectural Invariant:** GitHub Compare API as Single Source of Truth; atomic status invalidation.
- **Deliverables:**
  - Orchestrator queries GitHub `CompareCommits` when receiving a `RELEASE` event.
  - Matching pending PR drafts in `approval_posts` are transitioned atomically to `SUPERSEDED_BY_RELEASE`.
  - Discord PR embeds are updated with an informative badge advising publication of the parent release.

---

## 7. 🩺 Extended Health Probes & Dependency Observability
- **Scenario:** Silent failures occur when external dependencies (GitHub token, Gemini API key, Discord channels) degrade.
- **Current State:** Active probes limited to database connection and Discord Gateway WebSocket (`/health/discord-bot`).
- **Architectural Invariant:** Non-intrusive synthetic sub-endpoints reporting upstream health without leaking credentials.
- **Deliverables:**
  - `/health/github-api`: Validates token validity and remaining API budget.
  - `/health/gemini`: Validates Gemini API key and model readiness.
  - Aggregate status reporting within master `/health` endpoint.

---

## 8. 📜 Consolidated Version Digest & Adaptive Synthesis Engine
- **Scenario:** Major and minor releases lose depth when iterative patch histories (`v1.0.1`...`v1.0.7`) are fragmented across individual drafts.
- **Current State:** Version drafts only reflect immediate tag release notes without historical awareness of predecessor patches.
- **Architectural Invariant:** Coordinated triggers (automated milestone vs. manual on-demand) backed by database idempotency and reusable GitHub Compare infrastructure.
- **Deliverables:**
  - **Automated Milestone Synthesis**: Compiles all unmerged patch notes into a cohesive evolutionary post upon `vX.Y.0` releases.
  - **On-Demand Slash Command (`/digest base:<tag> head:<tag>`)**: Interactive generation in `#bot-control`.
  - **Narrative Prompting**: Thematic extraction of architecture pillars over raw commit logs.

---

## 9. 🔄 24h Daily Cadence Rule & Staging Event Queue
- **Scenario:** Multiple PRs and rapid patches published on the same day flood LinkedIn feeds with low-value atomic posts.
- **Current State:** System reactively generates an isolated draft per event.
- **Architectural Invariant:** Maximum one single publication per 24-hour cycle.
- **Deliverables:**
  - Staging table in PostgreSQL (`event_staging`) accumulating daily merged events.
  - 24-hour scheduled cadence worker or trigger synthesizing daily staged events into a single consolidated post.

---

## 10. 🛡️ External API Resilience, Circuit Breaker & Retry Mechanism
- **Scenario:** Transient external API rate limits (LinkedIn `HTTP 429`, Gemini timeouts) cause permanent pipeline drops.
- **Current State:** Errors result in immediate rejection without deferred retries.
- **Architectural Invariant:** Non-blocking exponential backoff with jitter and explicit recoverable error states.
- **Deliverables:**
  - Retry policy wrapper around LinkedIn and Gemini HTTP adapters.
  - `FAILED_RETRYABLE` post status and automated alert dispatching to `#bot-control`.