# Usage

## Authentication

All protected endpoints require a JWT token. Obtain one by registering and logging in:

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass","name":"Jane"}'

# Login — returns a JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass"}'
```

Use the token in subsequent requests:

```bash
export TOKEN="eyJhbG..."
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/profile
```

## Core Workflows

### Food Tracking

```bash
# Log a food entry
curl -X POST http://localhost:8080/api/v1/tracking/food \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Apple","calories":95,"protein":0.5,"carbs":25,"fat":0.3}'

# Log water intake
curl -X POST http://localhost:8080/api/v1/tracking/water \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount_ml":500}'

# View today's nutrition summary
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/tracking/summary
```

### Barcode Scanning

Look up food items by barcode (powered by OpenFoodFacts):

```bash
curl http://localhost:8080/api/v1/nutrition/barcode/4006381333931
```

### Recipe Management

```bash
# Create a recipe
curl -X POST http://localhost:8080/api/v1/recipes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Pasta Carbonara","servings":2,"ingredients":[...],"steps":[...]}'

# Search recipes
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/recipes?q=pasta"

# Favorite a recipe
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/recipes/1/favorite
```

### Meal Planning

```bash
# Generate a meal plan
curl -X POST http://localhost:8080/api/v1/meal-plans/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"days":7}'

# View the current active plan
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/meal-plans/current
```

### Shopping Lists

```bash
# Generate a shopping list from the active meal plan
curl -X POST http://localhost:8080/api/v1/shopping-lists \
  -H "Authorization: Bearer $TOKEN"

# View the current shopping list
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/shopping-lists/current

# Toggle an item as bought
curl -X PATCH -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/shopping-lists/items/42
```

### Pantry & Leftover Recipes

```bash
# Set pantry contents
curl -X POST http://localhost:8080/api/v1/pantry \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"items":["chicken","rice","broccoli"]}'

# Find recipes matching your pantry
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/pantry/recipes
```

### Price Tracking

```bash
# Log a grocery price
curl -X POST http://localhost:8080/api/v1/prices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"item":"Organic Milk","store":"Aldi","price_cents":189}'

# Compare prices across stores
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/prices/compare?item=Organic+Milk"

# View spending trends
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/prices/spending
```

### Social Features

```bash
# Invite a friend by email
curl -X POST http://localhost:8080/api/v1/friends/invite \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"friend@example.com"}'

# Share a recipe (generates a shareable link)
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/recipes/5/share

# View the leaderboard
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/leaderboard
```

### Data Export

```bash
# Export all tracking data as CSV
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/export/csv -o export.csv

# Export as JSON
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/export/json -o export.json
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server listen port |
| `DATABASE_PATH` | `bitewise.db` | Path to SQLite database file |
| `JWT_SECRET` | *(required)* | Secret key for signing JWT tokens |
| `JWT_EXPIRY_HOURS` | `72` | Token validity period in hours |

## Mobile App

The Expo-based mobile app provides the full BiteWise experience with:

- Tab-based navigation: Home, Recipes, Meal Plan, Shopping, Track
- Barcode scanner for quick food lookups
- Offline caching for use without connectivity
- Language toggle between English and German

Start the mobile app with `make dev-mobile` or `cd mobile && pnpm start`.

## Troubleshooting

### "JWT_SECRET environment variable is required"

Set the `JWT_SECRET` env var before starting the server:

```bash
export JWT_SECRET="any-strong-random-string"
```

### Database locked errors

SQLite is configured with WAL mode for concurrent reads, but heavy concurrent writes can still cause lock contention. For production use, ensure only one server instance writes to the database file.

### Mobile app can't reach the API

When running on a physical device, `localhost` won't resolve to your dev machine. Update the API URL in `mobile/services/api.ts` to your machine's LAN IP (e.g., `http://192.168.1.42:8080`).
