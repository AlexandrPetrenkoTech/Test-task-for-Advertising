# Advertising Service

This project is a REST API for posting and retrieving classified ads, built with Go (Golang) and following Clean Architecture principles.

## 🚀 Features

- Create new advertisements with title, description, photo URLs, and price.
- Get ad by ID (basic fields or full info via `fields=true`).
- List ads with pagination (10 items per page) and sorting by price or creation date (ascending/descending).
- Configuration via `.env` and `config.yaml` using Viper.
- Graceful shutdown support.
- Fully containerized with Docker and Docker Compose.
- Unit testing and linting support with Makefile.

## 🗂 Project Structure

├── cmd             →  `cmd/`             Entry point (main.go)
├── pkg/            →  `pkg/`             Clean Architecture layers
│   ├── handler/                      HTTP handlers
│   ├── service/                      Business logic
│   ├── repository/                   Database (sqlx)
│   ├── model/                        Domain models
│   └── error_message/                Error definitions
├── configs/        →  `configs/`        Config loaders (.env + YAML)
├── migrations/     →  `migrations/`     SQL schema migrations
├── scripts/        →  `scripts/`        Test data SQL
├── test/           →  `test/`           Test DB init
├── e2e_test.go     →  `e2e_test.go`     End-to-end tests (goconvey)
├── Dockerfile      →  `Dockerfile`      Multi-stage build
├── docker-compose.yml →  `docker-compose.yml`  Compose setup
├── Makefile        →  `Makefile`        Build, lint, test, docker
├── .env            →  `.env`            Env vars overrides
└── README.md       →  `README.md`       This file


## 🛠 Technologies Used

- Go 1.21+
- Echo
- PostgreSQL
- sqlx
- Viper
- Testify + mock
- golangci-lint
- Docker & Docker Compose
- Make
- migrate CLI (github.com/golang-migrate/migrate)

## 📦 Getting Started

### Example `config.yaml` (in `configs/`)

```yaml
server:
  host: "0.0.0.0"
  port: 8080

db:
  host: "db"
  port: 5432
  user: "user"
  password: "password"
  name: "advertising"
  
.env
# PostgreSQL connection (overrides config.yaml)
DB_HOST=db
DB_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=advertising

# HTTP server
PORT=8080
Using Docker & Docker Compose
docker-compose up --build
Your application will be available at http://localhost:8080.

Database Migrations
make migrate-up     # apply DB schema migrations
make migrate-down   # rollback migrations

Run Tests and Lint
make test           # runs unit tests
make e2e-test       # runs end-to-end tests (goconvey)
make lint           # runs golangci-lint

Build and Run Locally
make build
./advertising
🧑‍💻 Author
Alexandr Petrenko

This project was built as a technical test assignment demonstrating backend development skills, clean code, testing, and containerization.