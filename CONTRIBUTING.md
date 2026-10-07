# Contributing to Git-To-Feed 🚀

Thank you for your interest in contributing to **`git-to-feed`**! We welcome contributions from developers of all skill levels. To ensure a smooth, high-quality collaboration process, please review and adhere to the following guidelines.

---

## 📜 Table of Contents
- [Contributing to Git-To-Feed 🚀](#contributing-to-git-to-feed-)
  - [📜 Table of Contents](#-table-of-contents)
  - [🤝 Code of Conduct](#-code-of-conduct)
  - [🌿 Git Flow \& Branching Strategy](#-git-flow--branching-strategy)
  - [💬 Conventional Commits \& Gitmoji Standard](#-conventional-commits--gitmoji-standard)
    - [Commit Syntax](#commit-syntax)
    - [Examples:](#examples)
  - [📋 Repository Templates](#-repository-templates)
  - [🏗️ Clean Architecture \& Quality Standards](#️-clean-architecture--quality-standards)
  - [🚀 How to Submit a Pull Request](#-how-to-submit-a-pull-request)
  - [🎯 Finding Issues to Work On](#-finding-issues-to-work-on)

---

## 🤝 Code of Conduct
By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please report unacceptable behavior to the repository maintainers.

---

## 🌿 Git Flow & Branching Strategy
We strictly follow **Git Flow**:
- **`main`**: Production-ready release code. Never push directly to `main`.
- **`release`**: Target branch for production deployments (Render / Supabase).
- **`develop`**: Integration branch for new features and bug fixes.
- **Feature Branches**: Branch from `develop` using the naming convention `feature/<short-description>` or `fix/<short-description>`.

---

## 💬 Conventional Commits & Gitmoji Standard
All commit messages **must be written in English** following the [Conventional Commits](https://www.conventionalcommits.org/) specification coupled with official **Gitmojis**.

### Commit Syntax
`<type>: :<gitmoji_code>: <concise description in lower case>`

### Examples:
- `feat: :sparkles: add support for release tag events`
- `fix: :bug: enforce canonical published action for release idempotency`
- `docs: :memo: fix mermaid diagram rendering syntax`
- `build: :whale: install ca-certificates in runtime image`

---

## 📋 Repository Templates
When creating issues or opening pull requests, you **must use our official templates** located in `.github/`:

1. **Pull Request Templates (`.github/PULL_REQUEST_TEMPLATE/`)**:
   - [`feature.md`](.github/PULL_REQUEST_TEMPLATE/feature.md) - For new features and enhancements.
   - [`bugfix.md`](.github/PULL_REQUEST_TEMPLATE/bugfix.md) - For bug fixes and patch resolutions.
   - [`refactor.md`](.github/PULL_REQUEST_TEMPLATE/refactor.md) - For code restructuring without behavior changes.
   - [`docs.md`](.github/PULL_REQUEST_TEMPLATE/docs.md) - For documentation updates.
   - [`hotfix.md`](.github/PULL_REQUEST_TEMPLATE/hotfix.md) - For critical production emergency fixes.

2. **Issue Templates (`.github/`)**:
   - Refer to [`issue_template_base.md`](.github/issue_template_base.md) for bug reports and feature requests.

3. **Commit Template Reference**:
   - Refer to [`commit_template_base.md`](.github/commit_template_base.md) for commit formatting rules.

---

## 🏗️ Clean Architecture & Quality Standards
- **Domain Decoupling**: Business logic inside `internal/` must remain pure Go. Domain entities must **never** import infrastructure packages (`lib/pq`, `gorm`, HTTP frameworks, etc.).
- **Static Analysis**: Run `go vet ./...` before submitting your PR.
- **Unit Testing**: Run `go test -v ./...` to verify all test suites pass.

---

## 🚀 How to Submit a Pull Request
1. Fork the repository and create your feature branch from `develop`.
2. Implement your changes following our Clean Architecture standards.
3. Verify that `go vet ./...` and `go test ./...` execute cleanly.
4. Commit your changes using Conventional Commits with Gitmojis in English.
5. Open a Pull Request targeting the **`develop`** branch using the appropriate template in `.github/PULL_REQUEST_TEMPLATE/`.

---

## 🎯 Finding Issues to Work On
Before starting development, check our open issues and milestone boards:
1. **Explore Active Milestones:** Check [Milestone v1.1.0](https://github.com/hawksxo/git-to-feed/milestone/1) to see prioritized tasks for the next release.
2. **Filter by Community Labels:**
   - [`good first issue`](https://github.com/hawksxo/git-to-feed/labels/good%20first%20issue): Well-scoped tasks ideal for newcomers that require no private infrastructure or secrets.
   - [`help wanted`](https://github.com/hawksxo/git-to-feed/labels/help%20wanted): Open tasks where community contributions and PRs are actively welcomed.
3. **Claim an Issue:** Leave a comment on the issue you wish to tackle so maintainers can assign it to you and avoid duplicate effort.
