# Installation

## Prerequisites

- **Go** 1.26 or later — [golang.org/dl](https://golang.org/dl/)
- **Node.js** 18+ — [nodejs.org](https://nodejs.org/)
- **pnpm** — `npm install -g pnpm` (or via [pnpm.io](https://pnpm.io/installation))
- **Expo CLI** — installed automatically via pnpm
- **Make** — pre-installed on macOS/Linux; on Windows use WSL or `choco install make`

## Backend Setup

```bash
git clone https://github.com/RainyRoot/bitewise.git
cd bitewise
```

### 1. Configure environment variables

The backend requires a `JWT_SECRET` at minimum:

```bash
export JWT_SECRET="replace-with-a-strong-random-string"

# Optional — defaults shown:
export PORT=8080
export DATABASE_PATH="bitewise.db"
export JWT_EXPIRY_HOURS=72
```

You can also create a `.env` file in the project root (it is gitignored).

### 2. Run database migrations

```bash
make migrate
```

This creates the SQLite database file and applies all migration scripts from `backend/migrations/`.

### 3. Start the backend

```bash
make dev-backend
```

The API server starts on `http://localhost:8080`. Verify with:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### 4. Build a production binary (optional)

```bash
make build
# Output: backend/bin/server
```

## Mobile App Setup

```bash
cd mobile
export PATH="$HOME/.local/bin:$PATH"
pnpm install
```

### Start the Expo dev server

```bash
pnpm start
# or: make dev-mobile (from the project root)
```

Scan the QR code with the Expo Go app on your phone, or press `a` for Android emulator / `i` for iOS simulator.

> **Note:** The mobile app connects to `http://localhost:8080` by default. If running on a physical device, update the API base URL in `mobile/services/api.ts` to your machine's local IP.

## Run Everything at Once

```bash
make dev
```

This starts both the backend API and the Expo mobile dev server in parallel.

## Docker Setup

If you only need the backend API, Docker is the fastest way to get started:

```bash
# Pull the image
docker pull ghcr.io/rainyroot/bitewise:latest

# Run with a persistent data volume
docker run -d \
  --name bitewise \
  -p 8080:8080 \
  -e JWT_SECRET="your-secret-key" \
  -v bitewise-data:/data \
  -e DATABASE_PATH=/data/bitewise.db \
  ghcr.io/rainyroot/bitewise:latest
```

### Build the image locally

```bash
docker build -t bitewise .
docker run -d -p 8080:8080 -e JWT_SECRET="your-secret" bitewise
```

## Verify Installation

After starting the backend, run:

```bash
# Health check
curl http://localhost:8080/health

# Register a test user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123","name":"Test User"}'
```
