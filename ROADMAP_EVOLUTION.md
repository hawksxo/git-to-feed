# 🚀 Architecture & Feature Evolution Roadmap (Git-To-Feed)

This document records the architectural analysis and future enhancements identified for the technical content generation pipeline.

---

## 1. 🔄 Frequency Management & Deduplication (Spam Prevention)
- **Scenario:** Multiple Pull Requests (`FEATURE` / `REFACTOR`) are merged on the same day.
- **Current State:** The system generates one draft per PR. The interactive Discord bot (Human-in-the-Loop) allows the developer to approve only the highest-impact post or schedule them over time.
- **Future Evolution (Batching Engine):** Implement a aggregator ("Daily Digest") to group multiple PRs merged within a 6-hour window into a single consolidated publication.

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

## 🌐 Multi-Repository Support (Universal GitHub Webhook)

- **Does it work only with `git-to-feed` or any repository?**
  - **UNIVERSAL:** The `git-to-feed` microservice is completely repo-agnostic and multi-repository compatible.
  - Dynamically extracts `repository.full_name`, `pull_request.html_url`, `release.html_url`, etc., from the standard GitHub payload.
  - You can configure this same endpoint (`/api/v1/webhooks/github`) across **10, 20, or 100 different repositories** in your GitHub account or organization.
  - Sharing the same secret key (`GITHUB_WEBHOOK_SECRET`) across repositories will trigger instant ingestion for all of them.
  - **Future SaaS Multi-Tenant Enhancement:** Allow dynamic HMAC secrets per repository or organization token to separate Discord/LinkedIn configurations per project.
