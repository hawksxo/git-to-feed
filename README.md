<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![Go Version][go-shield]][go-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/hawksxo/git-to-feed">
    <img src="https://raw.githubusercontent.com/github/explore/80688e429a7d4ef2fca1e82350fe8e3517d3494d/topics/go/go.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Git-To-Feed</h3>

  <p align="center">
    Automated Event-Driven Microservice converting GitHub Webhooks into LinkedIn posts via Gemini AI & Discord Human-in-the-Loop approval.
    <br />
    <a href="https://github.com/hawksxo/git-to-feed"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/hawksxo/git-to-feed/issues">Report Bug</a>
    &middot;
    <a href="https://github.com/hawksxo/git-to-feed/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
        <li><a href="#architecture--package-by-feature-layout">Architecture & Package-by-Feature Layout</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#environment-configuration">Environment Configuration</a></li>
        <li><a href="#installation--running">Installation & Running</a></li>
      </ul>
    </li>
    <li><a href="#usage--api-endpoints">Usage & API Endpoints</a></li>
    <li><a href="#docker-deployment">Docker Deployment</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

`git-to-feed` is a production-ready, event-driven Go microservice designed under **Clean Architecture** and **Package-by-Feature** principles. It listens to GitHub repository Webhooks (Pull Requests and Releases), enriches event contexts, formats payloads using LLMs (Google Gemini AI) guided by curated **Few-Shot Examples**, and notifies an interactive **Discord Bot** for Human-in-the-Loop approval before live publication to **LinkedIn**.

### System Architecture Flow

```mermaid
graph TD
    GH["GitHub Webhook"] -->|HMAC Verified| WH["Webhook Handler"]
    WH --> PIPE["Pipeline Gemini AI"]
    PIPE --> BOT["Discord Bot"]
    BOT -->|Approve/Reject| APP["Approval Engine"]
    APP -->|Published| LI["LinkedIn API"]
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Architecture & Package-by-Feature Layout

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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

* [![Go][Go-shield]][Go-url]
* [![Google Gemini][Gemini-shield]][Gemini-url]
* [![Discord][Discord-shield]][Discord-url]
* [![PostgreSQL][Postgres-shield]][Postgres-url]
* [![Docker][Docker-shield]][Docker-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

Follow these steps to set up `git-to-feed` locally for development and testing.

### Prerequisites

* **Go 1.22+**: [Install Go](https://go.dev/dl/)
* **Docker & Docker Compose** (Optional for containerization)

### Environment Configuration

Create a `.env` file in the root directory (see `.env.example`):

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

### Installation & Running

1. Clone the repository:
   ```sh
   git clone https://github.com/hawksxo/git-to-feed.git
   cd git-to-feed
   ```
2. Download dependencies:
   ```sh
   go mod download
   ```
3. Run the application:
   ```sh
   go run cmd/api/main.go
   ```
4. Run tests & static analysis:
   ```sh
   go vet ./...
   go test -v ./...
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE -->
## Usage & API Endpoints

The microservice exposes clean REST HTTP endpoints:

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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- DOCKER DEPLOYMENT -->
## Docker Deployment

Build and run using Docker Compose:

```sh
docker-compose up --build -d
```

Or using Docker CLI directly:

```sh
docker build -t git-to-feed:latest .
docker run -d -p 8080:8080 --env-file .env git-to-feed:latest
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap

- [x] Universal Multi-Repository Webhook Ingestion (`pull_request`, `release`)
- [x] Google Gemini AI Integration with Curated Few-Shot Engineering
- [x] Anti-Cringe & Payload Sanitization Filters
- [x] Interactive Discord Bot for Human-in-the-Loop Approvals
- [x] Automatic PostgreSQL Schema Migration Engine
- [ ] Event Batching & Aggregation Engine (Frequency Windows)
- [ ] Multi-Tenant Repository Rules & Custom Prompts

See [ROADMAP_EVOLUTION.md](ROADMAP_EVOLUTION.md) for full architectural roadmap specifications.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'feat: add AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

hawksxo - [@hawksxo](https://github.com/hawksxo)

Project Link: [https://github.com/hawksxo/git-to-feed](https://github.com/hawksxo/git-to-feed)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/hawksxo/git-to-feed.svg?style=for-the-badge
[contributors-url]: https://github.com/hawksxo/git-to-feed/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/hawksxo/git-to-feed.svg?style=for-the-badge
[forks-url]: https://github.com/hawksxo/git-to-feed/network/members
[stars-shield]: https://img.shields.io/github/stars/hawksxo/git-to-feed.svg?style=for-the-badge
[stars-url]: https://github.com/hawksxo/git-to-feed/stargazers
[issues-shield]: https://img.shields.io/github/issues/hawksxo/git-to-feed.svg?style=for-the-badge
[issues-url]: https://github.com/hawksxo/git-to-feed/issues
[license-shield]: https://img.shields.io/github/license/hawksxo/git-to-feed.svg?style=for-the-badge
[license-url]: https://github.com/hawksxo/git-to-feed/blob/main/LICENSE
[go-shield]: https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go
[go-url]: https://go.dev/

[Go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[Gemini-shield]: https://img.shields.io/badge/Google%20Gemini-8E75B5?style=for-the-badge&logo=googlegemini&logoColor=white
[Gemini-url]: https://ai.google.dev/
[Discord-shield]: https://img.shields.io/badge/Discord%20Bot-5865F2?style=for-the-badge&logo=discord&logoColor=white
[Discord-url]: https://discord.com/
[Postgres-shield]: https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white
[Postgres-url]: https://www.postgresql.org/
[Docker-shield]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
