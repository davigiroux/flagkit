# FlagKit

A lightweight, self-hostable feature flag service.

Define flags, target them to users or environments, and evaluate them from any application via REST API or TypeScript SDK.

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌──────────┐
│  Dashboard   │────▶│   Go API     │────▶│ Postgres │
│  React/Vite  │     │  chi router  │     └──────────┘
└─────────────┘     │              │     ┌──────────┐
                    │              │────▶│  Redis   │
┌─────────────┐     │              │     └──────────┘
│  SDK (npm)   │────▶│  :8080       │
│  TypeScript  │     └──────────────┘
└─────────────┘
```

- **API** — Go, chi router, pgx, go-redis. Handles flag CRUD, evaluation, and audit logging.
- **Dashboard** — React, Vite, TailwindCSS, TanStack Query, dnd-kit. Manages flags and views audit logs.
- **SDK** — TypeScript, zero dependencies. Wraps the evaluation endpoint with in-memory caching.

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 22+ and pnpm
- Docker (for Postgres and Redis)

### 1. Clone and install

```bash
git clone https://github.com/davigiroux/flagkit.git
cd flagkit
pnpm install
```

### 2. Start Postgres and Redis

```bash
docker compose up -d
```

### 3. Start the API

```bash
cp api/.env.example api/.env
cd api && go run ./cmd/server
```

On first run, an API key is printed to the console:

```
========================================
  BOOTSTRAP API KEY (save this!):
  fk_abc123...
========================================
```

### 4. Start the Dashboard

```bash
pnpm --filter dashboard dev
```

Open http://localhost:5173, go to **Settings**, and paste your API key.

## API Reference

All endpoints (except `/health`) require `Authorization: Bearer <api_key>`.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check (no auth) |
| `GET` | `/flags` | List all flags |
| `POST` | `/flags` | Create a flag |
| `GET` | `/flags/:key` | Get a flag by key |
| `PATCH` | `/flags/:key` | Update a flag |
| `DELETE` | `/flags/:key` | Delete a flag |
| `POST` | `/flags/:key/toggle` | Toggle enabled state |
| `GET` | `/evaluate/:key` | Evaluate a flag |
| `GET` | `/audit` | List audit logs (paginated) |

### Create a flag

```bash
curl -X POST http://localhost:8080/flags \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "key": "new-checkout",
    "name": "New Checkout",
    "environment": "production",
    "rules": [
      { "type": "allowlist", "userIds": ["vip-1", "vip-2"] },
      { "type": "percentage", "value": 50 }
    ]
  }'
```

### Evaluate a flag

```bash
curl "http://localhost:8080/evaluate/new-checkout?user_id=user-123&environment=production" \
  -H "Authorization: Bearer $API_KEY"
```

Response:

```json
{
  "enabled": true,
  "reason": "rollout",
  "flagKey": "new-checkout",
  "evaluatedAt": "2026-03-14T10:00:00Z"
}
```

### Evaluation logic

Rules evaluate top-to-bottom. First match wins.

1. If `enabled` is `false` → `{ enabled: false, reason: "flag_disabled" }`
2. If rule is `allowlist` and `userId` is in the list → `{ enabled: true, reason: "allowlist" }`
3. If rule is `percentage` and `hash(flagKey + userId) % 100 < value` → `{ enabled: true, reason: "rollout" }`
4. No rule matched → `{ enabled: false, reason: "no_match" }`

Percentage rollout uses consistent hashing (FNV-32a) — the same user always gets the same result.

## TypeScript SDK

### Install

```bash
npm install flagkit-sdk
```

### Usage

```typescript
import { FlagKit } from 'flagkit-sdk'

const flags = new FlagKit({
  apiKey: 'fk_...',
  baseUrl: 'https://your-api.railway.app',
  ttl: 30_000, // cache TTL in ms (default: 30s)
})

const enabled = await flags.isEnabled('new-checkout', {
  userId: 'user-123',
  environment: 'production',
})

if (enabled) {
  // show new checkout
}
```

### Caching behavior

- Results are cached in memory with a configurable TTL
- On cache miss, the SDK calls the API and stores the result
- On API failure, it returns the last cached value (stale fallback)
- If no cached value exists and the API is down, it returns `false` (safe default)

## Data Model

### Flag

| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Primary key |
| `key` | string | Unique slug identifier |
| `name` | string | Display name |
| `description` | string | Optional description |
| `enabled` | boolean | Global on/off toggle |
| `environment` | enum | `production`, `staging`, or `development` |
| `rules` | JSON | Array of targeting rules |
| `createdAt` | timestamp | Creation time |
| `updatedAt` | timestamp | Last update time |

### Rule types

**Percentage rollout** — rolls out to a percentage of users based on consistent hashing:

```json
{ "type": "percentage", "value": 50 }
```

**User allowlist** — enables for specific user IDs:

```json
{ "type": "allowlist", "userIds": ["user-1", "user-2"] }
```

## Environment Variables

### API

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://flagkit:flagkit@localhost:5432/flagkit?sslmode=disable` | Postgres connection string |
| `REDIS_URL` | `redis://localhost:6379` | Redis connection string |
| `PORT` | `8080` | Server port |
| `CORS_ORIGINS` | `http://localhost:5173` | Comma-separated allowed origins |

### Dashboard

| Variable | Default | Description |
|----------|---------|-------------|
| `VITE_API_URL` | `http://localhost:8080` | API base URL |

## Development

```bash
# Run everything
docker compose up -d          # Postgres + Redis
cd api && go run ./cmd/server # API on :8080
pnpm --filter dashboard dev   # Dashboard on :5173

# Tests
cd api && go test ./...       # Go tests (hash, evaluation, audit diff)
pnpm --filter flagkit-sdk test # SDK tests

# Build
pnpm --filter dashboard build  # Dashboard production build
pnpm --filter flagkit-sdk build # SDK ESM + CJS build
```

## Deploy to Railway

The project includes Dockerfiles for both the API and Dashboard:

- **API**: `api/Dockerfile` — multi-stage Go build, runs migrations on startup
- **Dashboard**: `dashboard/Dockerfile` — Vite build served via nginx

Railway services needed:
1. **API** — from `api/Dockerfile`, env vars: `DATABASE_URL`, `REDIS_URL`, `PORT`, `CORS_ORIGINS`
2. **Dashboard** — from `dashboard/Dockerfile`, build arg: `VITE_API_URL`
3. **Postgres** — Railway managed
4. **Redis** — Railway managed

## License

MIT
