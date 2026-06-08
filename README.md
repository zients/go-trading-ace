# Go Trading Ace

Go Trading Ace is a DeFi rewards backend prototype built from scratch in November 2024 as a one-week take-home assignment. It demonstrates a Go service that listens to Uniswap USDC/WETH swap events, accumulates user trading volume, calculates campaign rewards, and exposes point history, task status, and leaderboard APIs.

This repository is positioned as a portfolio prototype. It focuses on showing the core backend flow for an on-chain trading rewards campaign, not on claiming production readiness.

## What It Demonstrates

- Go backend development with Gin and Fx dependency injection.
- Ethereum log subscription with go-ethereum and Infura WebSocket RPC.
- Uniswap V2 `Swap` event parsing for the USDC/WETH pool.
- Redis-backed volume accumulation and leaderboard storage.
- PostgreSQL persistence for campaign tasks and completed task histories.
- Docker Compose local infrastructure.
- Swagger API documentation and GitHub Actions CI.

## Architecture Flow

```text
Uniswap USDC/WETH Swap event
  -> Ethereum listener
  -> USDC amount extraction
  -> Campaign service
  -> Redis volume accumulation
  -> Reward calculation
  -> PostgreSQL task histories
  -> Campaign APIs
```

Redis is used for fast-changing campaign state: per-period swap totals, per-user accumulated volume, current task cache, and leaderboard sorted sets.

PostgreSQL is used for durable campaign records: task definitions and completed task histories.

## Main Components

- `main.go`: application wiring, database/Redis setup, Gin server startup.
- `services/ethereum_service.go`: Infura connection, Uniswap swap subscription, event parsing.
- `services/campaign_service.go`: campaign initialization, volume tracking, reward calculation, leaderboard reads.
- `repositories/`: PostgreSQL access for tasks and task histories.
- `helpers/redis_helper.go`: Redis key/value, hash, and sorted set operations.
- `controllers/` and `routes/`: HTTP API surface.
- `migrations/`: PostgreSQL schema setup.

## API Surface

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/` | Health/demo response |
| `GET` | `/campaign/start` | Initialize campaign tasks |
| `GET` | `/campaign/histories/:address` | Get point history for an address |
| `GET` | `/campaign/tasks/:address` | Get onboarding/share-pool task status for an address |
| `GET` | `/campaign/leaderboard/:taskName/:period` | Get leaderboard entries for a task period |
| `GET` | `/swagger/index.html` | Swagger UI |

## Requirements

- Go toolchain `1.26.4`
- Docker and Docker Compose
- `golang-migrate`
- Infura project ID for Ethereum mainnet WebSocket access

Install `golang-migrate` on macOS:

```bash
brew install golang-migrate
```

## Configuration

The application reads configuration from `config/config.yml`.

Example local Docker configuration:

```yaml
server:
  port: 8080
  request_timeout_seconds: 10

database:
  host: "postgres"
  user: "root"
  password: "root"
  port: 5432
  name: "trading-ace"
  sslmode: "disable"

redis:
  prefix: "trading-ace:"
  host: "redis"
  port: 6379

infura:
  key: "<your-infura-project-id>"
```

## Run Locally With Docker Compose

Clone the repository:

```bash
git clone https://github.com/zients/go-trading-ace.git
cd go-trading-ace
```

Start PostgreSQL, Redis, and the app:

```bash
docker compose up --build
```

Run database migrations from another terminal:

```bash
migrate -path=migrations \
  -database "postgres://root:root@localhost:5432/trading-ace?sslmode=disable" \
  up
```

Open Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

## Run Without Docker

For local execution outside Docker, update `config/config.yml` so PostgreSQL and Redis point to reachable local services, then run:

```bash
go mod download
go run main.go
```

## Tests

Run the test suite:

```bash
go test ./...
```

Run coverage:

```bash
go test ./... -cover
```

Run vulnerability analysis:

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

## Prototype Scope

This project was intentionally scoped as a one-week prototype. It demonstrates the core data flow, but it does not include every production hardening concern.

Implemented:

- Real Ethereum event subscription for the USDC/WETH Uniswap V2 pool.
- Campaign task initialization.
- Redis volume accumulation.
- PostgreSQL task history persistence.
- Leaderboard API backed by Redis sorted sets.
- Dockerized local services and CI tests.

Not production-ready yet:

- Event listener reconnect/backoff behavior.
- Historical backfill and event deduplication.
- Durable settlement worker state.
- Multi-instance locking for reward settlement.
- Automated database migrations during container startup.
- Secrets management beyond local YAML configuration.
- Multi-pool or multi-campaign configuration.

## Production Hardening Roadmap

If this prototype were extended into a production service, the next steps would be:

1. Replace the in-memory weekly ticker with a DB-driven settlement worker.
2. Track settlement status per task period to prevent duplicate reward distribution.
3. Store last processed block and support historical backfill.
4. Add reconnect/backoff handling for Ethereum subscriptions.
5. Add event deduplication based on transaction hash and log index.
6. Automate migrations as part of deployment.
7. Move secrets to environment variables or a secret manager.
8. Add integration tests covering API, PostgreSQL, and Redis together.
