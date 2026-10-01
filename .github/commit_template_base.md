# 📋 CONVENTION AND BASE COMMIT TEMPLATE

This template establishes the strict structure for all commits in the repository.

---

## 📌 General Format

```bash
git commit -m "<type>: :<gitmoji_code>: <concise description>

- <Bullet point 1 specifying the detailed change>
- <Bullet point 2 specifying the modified file or component>
- <Bullet point 3 specifying the test or validation performed>
```

---

## 🏷️ Commit Types (`<type>`)

| Type (`<type>`) | Description |
| :--- | :--- |
| **`feat`** | A new feature |
| **`fix`** | A bug fix |
| **`docs`** | Documentation only changes |
| **`style`** | Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc.) |
| **`refactor`** | A code change that neither fixes a bug nor adds a feature |
| **`perf`** | A code change that improves performance |
| **`test`** | Adding missing tests or correcting existing tests |
| **`build`** | Changes that affect the build system or external dependencies |
| **`ci`** | Changes to our CI configuration files and scripts |
| **`chore`** | Other changes that don't modify src or test files |
| **`revert`** | Reverts a previous commit |

---

## 🎨 Complete Catalog of Gitmojis (`:<gitmoji_code>:`)

| Gitmoji Code | Emoji | Description |
| :--- | :---: | :--- |
| `:art:` | 🎨 | Improve structure / format of the code |
| `:zap:` | ⚡️ | Improve performance |
| `:fire:` | 🔥 | Remove code or files |
| `:bug:` | 🐛 | Fix a bug |
| `:ambulance:` | 🚑️ | Critical hotfix |
| `:sparkles:` | ✨ | Introduce new features |
| `:memo:` | 📝 | Add or update documentation |
| `:rocket:` | 🚀 | Deploy stuff |
| `:lipstick:` | 💄 | Add or update the UI and style files |
| `:tada:` | 🎉 | Begin a project |
| `:white_check_mark:` | ✅ | Add, update, or pass tests |
| `:lock:` | 🔒️ | Fix security issues |
| `:closed_lock_with_key:` | 🔐 | Add or update secrets |
| `:bookmark:` | 🔖 | Release / Version tags |
| `:rotating_light:` | 🚨 | Fix compiler / linter warnings |
| `:construction:` | 🚧 | Work in progress |
| `:green_heart:` | 💚 | Fix CI Build |
| `:arrow_down:` | ⬇️ | Downgrade dependencies |
| `:arrow_up:` | ⬆️ | Upgrade dependencies |
| `:pushpin:` | 📌 | Pin dependencies to specific versions |
| `:construction_worker:` | 👷 | Add or update CI build system |
| `:chart_with_upwards_trend:` | 📈 | Add or update analytics or track code |
| `:recycle:` | ♻️ | Refactor code |
| `:heavy_plus_sign:` | ➕ | Add a dependency |
| `:heavy_minus_sign:` | ➖ | Remove a dependency |
| `:wrench:` | 🔧 | Add or update configuration files |
| `:hammer:` | 🔨 | Add or update development scripts |
| `:globe_with_meridians:` | 🌐 | Internationalization and localization |
| `:pencil2:` | ✏️ | Fix typos |
| `:poop:` | 💩 | Write bad code that needs to be improved |
| `:rewind:` | ⏪️ | Revert changes |
| `:twisted_rightwards_arrows:` | 🔀 | Merge branches |
| `:package:` | 📦️ | Add or update compiled files or packages |
| `:alien:` | 👽️ | Update code due to external API changes |
| `:truck:` | 🚚 | Move or rename resources (e.g: files, paths, routes, etc.) |
| `:page_facing_up:` | 📄 | Add or update license |
| `:boom:` | 💥 | Introduce breaking changes |
| `:bento:` | 🍱 | Add or update assets |
| `:wheelchair:` | ♿️ | Improve accessibility |
| `:bulb:` | 💡 | Add or update comments in source code |
| `:beers:` | 🍻 | Write code drunkenly |
| `:speech_balloon:` | 💬 | Add or update text and literals |
| `:card_file_box:` | 🗃️ | Perform database related changes |
| `:loud_sound:` | 🔊 | Add or update logs |
| `:mute:` | 🔇 | Remove logs |
| `:busts_in_silhouette:` | 👥 | Add or update contributor(s) |
| `:children_crossing:` | 🚸 | Improve user experience / usability |
| `:building_construction:` | 🏗️ | Make architectural changes |
| `:iphone:` | 📱 | Work on responsive design |
| `:clown_face:` | 🤡 | Mock things |
| `:egg:` | 🥚 | Add or update an easter egg |
| `:see_no_evil:` | 🙈 | Add or update a .gitignore file |
| `:camera_flash:` | 📸 | Add or update snapshots |
| `:alembic:` | ⚗️ | Perform experiments |
| `:mag:` | 🔍️ | Improve SEO |
| `:label:` | 🏷️ | Add or update types |
| `:seedling:` | 🌱 | Add or update seed files |
| `:triangular_flag_on_post:` | 🚩 | Add, update, or remove feature flags |
| `:goal_net:` | 🥅 | Catch errors |
| `:dizzy:` | 💫 | Add or update animations and transitions |
| `:wastebasket:` | 🗑️ | Deprecate code that needs to be cleaned up |
| `:passport_control:` | 🛂 | Work on code related to authorization, roles and permissions |
| `:adhesive_bandage:` | 🩹 | Simple fix for a non-critical issue |
| `:monocle_face:` | 🧐 | Data exploration/inspection |
| `:coffin:` | ⚰️ | Remove dead code |
| `:test_tube:` | 🧪 | Add a failing test |
| `:necktie:` | 👔 | Add or update business logic |
| `:stethoscope:` | 🪵 | Add or update healthcheck |
| `:bricks:` | 🧱 | Infrastructure related changes |
| `:technologist:` | 🧑‍💻 | Improve developer experience |
| `:money_with_wings:` | 💸 | Add sponsorships or money related infrastructure |
| `:thread:` | 🧵 | Add or update code related to multithreading or concurrency |
| `:safety_vest:` | 🦺 | Add or update code related to validation |

---

## 💡 Complete Example

```bash
git commit -m "feat: :necktie: implement ContextBuilder assembler service

- Create ContextBuilder constructor for FromRelease and FromPullRequest events
- Integrate automatic Markdown sanitization and optional human note injection
- Add unit tests and null-pointer resilience validation in context_builder_test.go
```
