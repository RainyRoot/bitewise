# BiteWise

A full-stack nutrition and meal planning platform with a Go REST API backend and a React Native (Expo) mobile app. Track your food intake, plan meals, manage recipes, monitor grocery prices, and build healthy habits — all from your phone.

## Features

- **Meal Planning** — Generate weekly meal plans based on your preferences and dietary needs
- **Food Tracking** — Log meals and water intake with barcode scanning support (OpenFoodFacts integration)
- **Recipe Management** — Create, search, and favorite recipes; share them with friends via unique links
- **Shopping Lists** — Auto-generate shopping lists from your active meal plan
- **Pantry Management** — Track what you have at home and find recipes that use your leftovers
- **Price Tracking** — Log grocery prices, compare stores, and monitor spending trends
- **Nutrition Diary** — Keep a daily food diary with monthly overviews
- **Social Features** — Add friends, share recipes, and compete on the leaderboard
- **Achievements & Streaks** — Stay motivated with gamified tracking milestones
- **Seasonal Calendar** — Discover what produce is in season
- **Data Export** — Export all your data as CSV or JSON at any time
- **Localization** — English and German language support in the mobile app

## Quick Start

```bash
# Clone the repo
git clone https://github.com/RainyRoot/bitewise.git
cd bitewise

# Run database migrations
export JWT_SECRET="your-secret-key"
make migrate

# Start both backend and mobile dev servers
make dev
```

The API runs on `http://localhost:8080` and the Expo dev server starts alongside it.

See [INSTALL.md](INSTALL.md) for full setup instructions and [USAGE.md](USAGE.md) for detailed usage.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.26+, chi router, raw SQL |
| Database | SQLite (modernc.org/sqlite, pure Go) |
| Auth | JWT (golang-jwt) + bcrypt |
| Mobile | React Native (Expo 52), TypeScript |
| Nutrition Data | OpenFoodFacts API |
| Package Manager | pnpm (mobile) |

## Project Structure

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
│   └── pkg/httputil/     # JSON response helpers
├── mobile/
│   ├── app/              # Expo Router screens (tabs + detail views)
│   ├── hooks/            # React hooks (auth)
│   ├── i18n/             # Internationalization (en/de)
│   ├── services/         # API client + offline support
│   └── types/            # TypeScript type definitions
└── Makefile
```

## Docker

Pull and run the backend API directly:

```bash
docker pull ghcr.io/rainyroot/bitewise:latest

docker run -d \
  -p 8080:8080 \
  -e JWT_SECRET="your-secret-key" \
  -v bitewise-data:/data \
  -e DATABASE_PATH=/data/bitewise.db \
  ghcr.io/rainyroot/bitewise:latest
```

## API Overview

All endpoints live under `/api/v1`. Public routes:

- `POST /api/v1/auth/register` — Create account
- `POST /api/v1/auth/login` — Get JWT token
- `GET /api/v1/nutrition/barcode/{code}` — Barcode lookup
- `GET /api/v1/nutrition/search` — Food search
- `GET /api/v1/seasonal` — Seasonal produce calendar
- `GET /api/v1/shared/{code}` — View shared recipe

All other routes require a valid `Authorization: Bearer <token>` header.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Push and open a pull request
