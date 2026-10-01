# 🚀 Git-To-Feed Microservice

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blue)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![Pattern](https://img.shields.io/badge/Pattern-Package--by--Feature-brightgreen)](#architecture--package-by-feature-layout)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

`git-to-feed` is an automated, event-driven Go microservice designed under **Clean Architecture** and **Package-by-Feature** principles. It listens to GitHub repository Webhooks (Pull Requests and Releases), enriches the context, formats the payload using LLM models (Google Gemini AI) guided by curated **Few-Shot Examples**, and notifies an interactive **Discord Bot** for Human-in-the-Loop approval before live publication to **LinkedIn**.

---

## 🌟 Architecture & Package-by-Feature Layout

The codebase strictly decouples business domain logic from infrastructure, storage, and transport layers.

```text
git-to-feed/
├── cmd/
│   └── api/                  # Main HTTP API Entrypoint
├── internal/
│   ├── approval/             # Approval Feature Module (Entities, Use Cases, Repositories, HTTP Handlers)
│   ├── pipeline/             # Content Generation Pipeline (ContextBuilder, PromptBuilder, GeminiLLMClient, Anti-Cringe, Few-Shot)
│   ├── publisher/            # LinkedIn Publisher Feature Module (Entities, Use Cases, HTTP Client, Mocks)
│   ├── webhook/              # GitHub Ingestion Feature Module (Payload Parser, HMAC Validation)
│   └── platform/             # Infrastructure & Platform Adapters
│       ├── config/           # Environment Configuration & godotenv Loader
│       ├── discord/          # Interactive Discord Bot (discordgo WebSocket & Embed Buttons)
│       ├── http/             # HTTP Server Router & Endpoint Wire-up
│       └── storage/          # Database Drivers (PostgreSQL/Supabase & SQL Auto-Migrations Engine)
├── migrations/               # SQL Database Schema Migrations (.sql)
├── .github/                  # GitHub Workflows, Templates & Community Guidelines
├── Dockerfile                # Multi-stage Production Docker Build
├── docker-compose.yml        # Orchestration Config
└── ROADMAP_EVOLUTION.md      # Future Architectural Roadmap & Evolution Specifications
```

---

## ⚙️ Features

- **Event-Driven GitHub Webhook Ingestion:** Validates incoming payloads using HMAC-SHA256 signatures (`X-Hub-Signature-256`). Supports `pull_request` (closed/merged) and `release` (published) events.
- **Universal Multi-Repository Compatibility:** Repository-agnostic design. A single deployed instance can ingest webhooks from 1, 10, or 100+ GitHub repositories.
- **LLM Content Transformation:** Integrates Google Gemini API (`gemini-3.5-flash-lite`, `gemini-3.8-flash`) via the official `google.golang.org/genai` SDK.
- **Curated Few-Shot Engineering:** Employs archetype-specific templates (`RELEASE`, `FEATURE`, `REFACTOR`) inspired by real-world technical posts (Alibaba Qwen, Google Gemini Skills, Anthropic Claude).
- **Anti-Cringe & Sanitization Filters:** Removes corporate buzzword cliches, cleans HTML comments/checkboxes, and caps excessive emoji density.
- **Interactive Discord Bot (Human-in-the-Loop):** Sends rich embeds with `🚀 Approve & Publish` and `❌ Reject` buttons directly to a designated Discord channel.
- **LinkedIn Publishing Engine:** Communicates with LinkedIn API v2 (`ugcPosts`) to share approved content live.
- **Multi-Storage Persistence:** Auto-detects `DATABASE_URL` for Supabase/PostgreSQL with automatic SQL migration runner, falling back to Thread-Safe InMemory repositories for local development and unit tests.

---

## 🛠️ Environment Configuration

Create a `.env` file in the root directory (refer to `.env.example`):

```ini
PORT=":8080"
GITHUB_WEBHOOK_SECRET="your-github-webhook-secret"
LINKEDIN_ACCESS_TOKEN="your-linkedin-access-token"
LINKEDIN_AUTHOR_URN="urn:li:person:your-author-urn"
DATABASE_URL="postgresql://user:password@host:5432/dbname?sslmode=require"
DISCORD_BOT_TOKEN="your-discord-bot-token"
DISCORD_CHANNEL_ID="your-discord-channel-id"
GEMINI_API_KEY="your-google-gemini-api-key"
GEMINI_MODEL_NAME="gemini-3.5-flash-lite"
```

---

## 🚀 Local Development

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- Docker & Docker Compose (Optional)

### Running Locally
```bash
# Clone the repository
git clone https://github.com/hawksxo/git-to-feed.git
cd git-to-feed

# Download dependencies
go mod download

# Run the API server
go run cmd/api/main.go
```

### Running Tests & Static Check
```bash
# Run static analysis
go vet ./...

# Run all unit and integration tests
go test -v ./...
```

---

## 🐳 Docker Deployment

Build and run using Docker Compose:

```bash
docker-compose up --build -d
```

Or using Docker directly:

```bash
docker build -t git-to-feed:latest .
docker run -d -p 8080:8080 --env-file .env git-to-feed:latest
```

---

## 📖 API Endpoints Summary

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Application health check (`200 OK`) |
| `POST` | `/api/v1/webhooks/github` | Ingests and processes GitHub Webhook events |
| `GET` | `/api/v1/approvals/pending` | Lists all draft posts awaiting approval |
| `POST` | `/api/v1/approvals/{uuid}/approve` | Approves and publishes a draft post to LinkedIn |
| `POST` | `/api/v1/approvals/{uuid}/reject` | Rejects and discards a pending draft post |
| `PUT` | `/api/v1/approvals/{uuid}/edit` | Edits content and approves post for publishing |
| `GET` | `/api/v1/examples` | Lists active Few-Shot examples |
| `POST` | `/api/v1/examples` | Creates a new Few-Shot example entity |

---

## 📄 Roadmap & Future Evolution

See [ROADMAP_EVOLUTION.md](file:///C:/Users/Posty/Documents/Backend/git-to-feed/ROADMAP_EVOLUTION.md) for detailed architectural specifications regarding event hierarchy, batching engines, and multi-tenant expansion.

---

## 📜 License

Distributed under the MIT License.
