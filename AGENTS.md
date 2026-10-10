# 🤖 AGENT & TECH LEAD OPERATING DIRECTIVES (Git-To-Feed)

This document establishes non-negotiable operational guidelines, lessons learned, and execution invariants for AI agents and developers collaborating on the `git-to-feed` codebase.

---

## 1. 🌿 Git Flow & PR Merge Standard (Non-Negotiable)

- **Strict Merge Topology:** All feature Pull Requests targeting `develop` (or `main`/`release`) MUST be merged preserving the complete commit history and explicit merge commit topology:
  ```bash
  gh pr merge <PR_NUMBER> --merge --delete-branch
  ```
- **Prohibited Flags:** **NEVER** use `--squash` or `--rebase`. The repository strictly enforces full commit history retention and atomic commit granularity.
- **Branch Naming:** Always use conventional prefixes without issue numbers (`feat/<feature-name>`, `chore/<task-name>`). Issue linking is handled strictly in PR bodies via `Closes #XX`.

---

## 2. 🐳 Canonical Local Test Runner (Windows AppLocker / Docker Parity)

- **Host Limitation:** The Windows host enforces AppLocker / Application Control policies that block unsigned temporary Go test binaries generated in `AppData\Local\Temp`.
- **Enforced Execution Command:** **NEVER** run host-level `go test`. Always execute tests inside the cached Docker container:
  ```powershell
  docker run --rm -v "${PWD}:/app" -v go-cache:/root/.cache/go-build -v go-mod-cache:/go/pkg/mod -w /app golang:latest go test -race ./...
  ```

---

## 3. 🛡️ String Manipulation & Prefix Defensive Coding

- **Prohibition of Magic Slice Offsets:** Do not use hardcoded integer slicing on dynamic interaction strings (e.g. `customID[:8]` or `customID[8:]`).
- **Standard Library Idiom:** Always favor defensive standard library methods from `strings`:
  - Identification: `strings.HasPrefix(customID, "discard_")`
  - Extraction: `strings.TrimPrefix(customID, "discard_")`
- **Clean Transport DTOs:** Never print raw structs or pointer memory layouts (`%v` on domain entities) directly to end-user transport interfaces (Discord embeds, HTTP responses, client logs). Extract explicit identifiers (e.g., `postRecord.UUID`).

---

## 4. 🧪 Test Assertion Discipline & Time Mutation Protocol

- **Time-Sensitive Tests (TTLs / Expirations):** Constructors and use cases frequently stamp timestamps internally (`time.Now()`). When testing expiration boundaries:
  1. Instantiate the entity via the use case / factory.
  2. Mutate the timestamp directly on the persisted struct (`post.CreatedAt = time.Now().Add(-15 * 24 * time.Hour)`).
  3. Re-save/update the entity in the test repository.
  4. Invoke the target method and assert the explicit domain error (`ErrDraftExpired`).
- **Error Assertions:** Ensure test assertions verify that the required domain error is returned (`if err != ErrExpected { t.Errorf(...) }`), avoiding inverted conditions.

---

## 5. 📚 Post-Merge Roadmap & Documentation Synchronization Protocol

- **Execution Timing:** **NEVER** modify roadmap status or README checkboxes inside the feature branch before PR merge.
- **Dedicated Commit on `develop`:** Only AFTER the PR is merged into `develop`, update:
  - `README.md`: Check off completed milestone deliverable `[x]`.
  - `ROADMAP_EVOLUTION.md`: Update header to `COMPLETED ✅` and document exact PR reference (`Completed in PR #XX`).
  - Commit message format: `docs: :memo: mark <feature> as completed in roadmap and readme`.
