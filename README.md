# BiteWise

A full-stack nutrition and meal planning platform with a Go REST API backend and a React Native (Expo) mobile app. Track your food intake, plan meals, manage recipes, monitor grocery prices, and build healthy habits — all from your phone.

## Features

- **Meal Planning** - Generate weekly meal plans based on your preferences and dietary needs
- **Food Tracking** - Log meals and water intake with barcode scanning support (OpenFoodFacts integration)
- **Recipe Management** - Create, search, and favorite recipes; share them with friends via unique links
- **Shopping Lists** - Auto-generate shopping lists from your active meal plan
- **Pantry Management** - Track what you have at home and find recipes that use your leftovers
- **Price Tracking** - Log grocery prices, compare stores, and monitor spending trends
- **Nutrition Diary** - Keep a daily food diary with monthly overviews
- **Social Features** - Add friends, share recipes, and compete on the leaderboard
- **Achievements & Streaks** - Stay motivated with gamified tracking milestones
- **Seasonal Calendar** - Discover what produce is in season
- **Data Export** - Export all your data as CSV or JSON at any time
- **Localization** - English and German language support in the mobile app

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.26+, chi router, raw SQL |
| Database | SQLite (modernc.org/sqlite, pure Go) |
| Auth | JWT (golang-jwt) + bcrypt |
| Mobile | React Native (Expo 52), TypeScript |
| Nutrition Data | OpenFoodFacts API |
| Package Manager | pnpm (mobile) |

---

## Running locally

### Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [Node.js](https://nodejs.org/) 20+ and [pnpm](https://pnpm.io/installation)
- [Expo Go](https://expo.dev/go) app on your phone (for testing on device)

### 1. Backend

```bash
cd backend

# Required: set a secret for JWT tokens (can be anything, keep it secret in production)
export JWT_SECRET=your-secret-here

# First run: set up the database
go run ./cmd/server/ -migrate

# Start the server (runs on port 8080)
go run ./cmd/server/
```

Optional environment variables:

| Variable | Default | Description |
|---|---|---|
| `JWT_SECRET` | *required* | Secret key for signing JWT tokens |
| `PORT` | `8080` | Port the server listens on |
| `DATABASE_PATH` | `bitewise.db` | Path to the SQLite database file |
| `JWT_EXPIRY_HOURS` | `72` | How long login tokens are valid |

### 2. Mobile app

```bash
cd mobile
pnpm install

# Tell the app where your backend is
# Replace with your machine's local IP (not localhost - your phone won't reach that)
export EXPO_PUBLIC_API_URL=http://192.168.x.x:8080

pnpm start
```

Then scan the QR code with Expo Go on your phone, or press `a` for Android emulator / `i` for iOS simulator.

> To find your local IP: `ip addr` on Linux, `ipconfig` on Windows, `ifconfig` on Mac.

### Both at once (convenience)

```bash
# From the repo root - starts backend and mobile together
make dev
```

Note: `make dev` uses `localhost` for the API URL, so the mobile app will only work in a simulator/emulator, not on a real device.

---

## Docker (backend only)

```bash
docker run -p 8080:8080 \
  -e JWT_SECRET=your-secret-here \
  -v $(pwd)/data:/app/data \
  -e DATABASE_PATH=/app/data/bitewise.db \
  ghcr.io/rainyroot/bitewise/backend:main
```

The database is stored at the path set by `DATABASE_PATH`. Mount a volume there so data survives container restarts.

First run with a fresh database, run migrations:

```bash
docker run --rm \
  -e JWT_SECRET=your-secret-here \
  -v $(pwd)/data:/app/data \
  -e DATABASE_PATH=/app/data/bitewise.db \
  ghcr.io/rainyroot/bitewise/backend:main \
  /app/server -migrate
```

### Docker Compose

```yaml
services:
  backend:
    image: ghcr.io/rainyroot/bitewise/backend:main
    ports:
      - "8080:8080"
    environment:
      JWT_SECRET: your-secret-here
      DATABASE_PATH: /app/data/bitewise.db
    volumes:
      - ./data:/app/data
```

---

## Releases

Pre-built backend binaries for Linux and Windows are attached to each [GitHub Release](https://github.com/RainyRoot/bitewise/releases).

Download the binary for your platform, then:

```bash
# Linux
chmod +x bitewise-server-linux-amd64
export JWT_SECRET=your-secret-here
./bitewise-server-linux-amd64 -migrate  # first run only
./bitewise-server-linux-amd64
```

```powershell
# Windows (PowerShell)
$env:JWT_SECRET = "your-secret-here"
.\bitewise-server-windows-amd64.exe -migrate  # first run only
.\bitewise-server-windows-amd64.exe
```

### Creating a new release

Tag a commit with a version number to trigger a release build:

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions will automatically build the binaries, create a release, and push a versioned Docker image.

---

## API Overview

All endpoints live under `/api/v1`. Public routes:

- `POST /api/v1/auth/register` - Create account
- `POST /api/v1/auth/login` - Get JWT token
- `GET /api/v1/nutrition/barcode/{code}` - Barcode lookup
- `GET /api/v1/nutrition/search` - Food search
- `GET /api/v1/seasonal` - Seasonal produce calendar
- `GET /api/v1/shared/{code}` - View shared recipe

All other routes require a valid `Authorization: Bearer <token>` header.

---

## Project structure

```
bitewise/
├── backend/
│   ├── cmd/server/       # Application entrypoint
│   ├── internal/
│   │   ├── config/       # Environment-based configuration
│   │   ├── domain/       # Pure domain types
│   │   ├── handler/      # HTTP handlers + middleware
│   │   ├── nutrition/    # OpenFoodFacts client, seasonal data
│   │   ├── repository/   # SQLite data access layer
│   │   ├── scraper/      # Chefkoch recipe scraper
│   │   └── service/      # Business logic
│   ├── migrations/       # SQL migration files
│   ├── pkg/httputil/     # JSON response helpers
│   └── Dockerfile
└── mobile/
    ├── app/              # Expo Router screens (tabs + detail views)
    ├── hooks/            # React hooks (auth)
    ├── i18n/             # Internationalization (en/de)
    ├── services/         # API client + offline support
    └── types/            # TypeScript type definitions
```

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Push and open a pull request
