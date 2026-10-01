# 🚀 Architecture & Feature Evolution Roadmap (Git-To-Feed)

This document records the architectural analysis, technical debt items, and future enhancements identified for the technical content generation pipeline.

---

## 1. 🔄 Frequency Management & Deduplication (Spam Prevention)
- **Scenario:** Multiple Pull Requests (`FEATURE` / `REFACTOR`) are merged on the same day.
- **Current State:** The system generates one draft per PR. The interactive Discord bot (Human-in-the-Loop) allows the developer to approve only the highest-impact post or schedule them over time.
- **Future Evolution (Batching Engine):** Implement an aggregator ("Daily Digest") to group multiple PRs merged within a 6-hour window into a single consolidated publication.

---

## 2. 👑 Event Hierarchy (Release vs. Pull Request)
- **Scenario:** A PR is approved and 5 minutes later a `Release Tag` is published (which includes that PR).
- **Current State:** Both generate independent drafts in Discord.
- **Future Evolution (Event Priority Matrix):** A `RELEASE` event holds higher hierarchy. The orchestrator will detect if there are recent PR drafts included in the release and mark them as `SUPERSEDED_BY_RELEASE`, notifying Discord that publishing only the global release note is recommended.

---

## 3. 🎯 Context Quality & Enrichment (Garbage In, Garbage Out)
- **Scenario:** PRs with brief or vague titles and descriptions.
- **Current State:** GitHub templates (`PULL_REQUEST_TEMPLATE` and `RELEASE_TEMPLATE`) enforce and guide good input documentation.
- **Future Evolution (Context Enrichment):** 
  - Query `diff_stats` or modified file lists via GitHub API if the PR body is concise.
  - Support supplementary human notes via Discord/API (`HumanNotes`) to prompt or instruct Gemini specifically.

---

## 4. ⚙️ Automated CI/CD Pipelines (GitHub Actions & Deployment Automation)
- **Scenario:** Manual execution of verification tests, Docker builds, and cloud service deployments.
- **Current State:** Manual release branch workflow and environment configuration in Render/Supabase.
- **Future Evolution (CI/CD Pipeline):**
  - **CI Workflow (`.github/workflows/ci.yml`)**: Automate `go vet ./...`, `go test -v ./...`, and static analysis execution on every PR targeting `develop` or `release`.
  - **Docker Build & Security Scan (`.github/workflows/docker.yml`)**: Build multi-stage Docker images, verify `.dockerignore` context, and run vulnerability scanning (Trivy/Grype).
  - **CD Workflow (`.github/workflows/cd.yml`)**: Trigger automatic deployment hooks to Render Web Service and execute Supabase SQL migrations automatically when changes are pushed to the `release` branch.


---

## 5. 📢 Multi-Channel Routing & Archetype Segmentation
- **Scenario:** High volume of webhook events cluttering a single Discord channel, making it hard to locate high-priority releases vs. standard PR drafts.
- **Current State:** All draft notifications are dispatched to a single unified Discord channel (`DISCORD_CHANNEL_ID`).
- **Future Evolution (Multi-Channel Router):**
  - **`DISCORD_CHANNEL_RELEASES_ID` (`#drafts-releases`)**: Dedicated high-priority channel receiving only major/minor release tags (`vX.Y.Z`).
  - **`DISCORD_CHANNEL_FEATURES_ID` (`#drafts-features`)**: Channel dedicated to pull request drafts and minor refactors.
  - **`DISCORD_CHANNEL_AUDIT_ID` (`#published-feed`)**: Audit log channel recording successful publication confirmations with direct LinkedIn URN post links.

---

## 6. 🧹 Interactive Draft Lifecycle & Purge Engine (Housekeeping)
- **Scenario:** Pending drafts accumulate over time when developers decide not to publish specific PRs or releases.
- **Current State:** Post approval cards remain static in Discord until manual action is taken.
- **Future Evolution (Draft Management Suite):**
  - **Interactive Discard Action (`🔴 Discard / Delete`)**: Button on Discord embed cards allowing immediate status transition to `REJECTED` in Supabase and automatic deletion/archival of the Discord message.
  - **Discord Slash Command (`/list-pending`)**: Interactive query command to locate, review, or re-trigger approval cards for historical pending releases (e.g., `v1.0.0`).
  - **Automatic TTL & Stale Purge Job**: Scheduled worker marking pending drafts older than 14 days as `EXPIRED`.

---

## 7. 📜 Consolidated Version Digest & Patch Merging
- **Scenario:** Developer wants to publish a comprehensive launch post for a major version (`v1.0.0`), but multiple patch releases (`v1.0.1`...`v1.0.4`) have already been deployed with critical architectural fixes.
- **Current State:** Each version tag generates an isolated draft without knowledge of subsequent or preceding patch history.
- **Future Evolution (Digest Engine):**
  - **Incremental History Enrichment**: Gemini AI pipeline queries recent tag history when generating a major release draft, synthesizing the initial release architecture with battle-tested patch enhancements.
  - **Manual Digest Re-trigger (`/digest v1.0.0 --upto v1.0.4`)**: Slash command to re-synthesize and consolidate a historical major post with all intermediate patch release notes into a unified, high-impact LinkedIn post.



