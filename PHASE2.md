# BiteWise — Phase 2 Implementation Guide

## What's Done (Phase 1)
The complete Phase 1 MVP is implemented and compiles cleanly:

### Backend (Go, chi router, SQLite)
- **Auth**: JWT-based register/login, auth middleware, password hashing (bcrypt)
- **Profile**: CRUD for user data, allergies, and food preferences
- **Recipes**: Repository with search/filter (LIKE, category, allergen exclusion), favorites
- **Scraper**: Chefkoch.de DataProvider interface + basic API integration
- **Meal Plans**: Generation (7 days × 3 meals, allergen-aware), CRUD, entry updates, activation
- **Tracking**: Food logs, water logs, daily nutrition summaries (SUM aggregation)
- **All routes wired** in `cmd/server/main.go` with auth middleware on protected routes
- **Shopping lists & Achievements**: Routes exist as placeholders

### Mobile (Expo/React Native + TypeScript)
- **Tab navigation**: Home, Plan, Track, Recipes, Profile
- **Auth screens**: Login, Register (wired to API client)
- **Home**: Calorie/macro overview, water tracking with quick-add buttons, today's meals
- **Plan**: Day selector, meal cards with swap/lock icons, generate button
- **Track**: Meal type selector, food log display, quick-add foods, manual entry, barcode button
- **Recipes**: Search, category filter, recipe cards with favorites
- **Recipe Detail**: Nutrition, ingredients, instructions, add-to-plan button
- **Profile**: Stats, allergies, settings menu, logout
- **API Client**: Full typed client in `services/api.ts` with JWT token management

### What's NOT wired yet
- Mobile screens use placeholder data (TODO comments mark where API calls go)
- `pnpm install` hasn't been run yet (just package.json exists)
- No `.env` file created yet for JWT_SECRET

## Phase 2 Tasks

### 1. Barcode Scanner (OpenFoodFacts)
**Backend:**
- Create `backend/internal/nutrition/openfoodfacts.go`
- Interface: `NutritionProvider` with `LookupBarcode(ctx, code string) (*FoodItem, error)` and `SearchFood(ctx, query string) ([]FoodItem, error)`
- API: `https://world.openfoodfacts.org/api/v2/product/{barcode}.json`
- Parse: product_name, nutriments (energy-kcal_100g, proteins_100g, carbohydrates_100g, fat_100g, fiber_100g)
- Create handler: `GET /api/v1/nutrition/barcode/{code}` and `GET /api/v1/nutrition/search?q=`
- Wire into main.go

**Mobile:**
- Add `expo-barcode-scanner` to package.json
- Create `app/scanner.tsx` — camera view with barcode overlay
- On scan → call `/api/v1/nutrition/barcode/{code}` → show result → allow logging directly
- Add scanner button on Track screen that navigates to scanner

### 2. Shopping List
**Backend:**
- Create `backend/internal/repository/shopping_repository.go` and `service/shopping_service.go`
- `GenerateFromMealPlan(ctx, userID, mealPlanID)` — iterate all entries, collect ingredients, merge duplicates (same name+unit → sum quantities), categorize (Obst, Gemüse, Milchprodukte, Fleisch, etc.)
- `ToggleItem(ctx, itemID, userID)` — toggle is_checked
- Replace placeholder routes in main.go

**Mobile:**
- Create shopping list screen (could be a new tab or accessible from Plan screen)
- Group items by category, show checkboxes
- "Generate from current plan" button

### 3. Saisonkalender (Seasonal Calendar)
**Backend:**
- Create `backend/internal/nutrition/seasonal.go` with a static map of German seasonal produce by month
- `GET /api/v1/seasonal` — returns current month's seasonal items
- Modify meal plan generation to prefer recipes with seasonal ingredients (bonus scoring)

### 4. Reste-Verwertung (Leftover Recipes)
**Backend:**
- `POST /api/v1/pantry` — save user's available ingredients
- `GET /api/v1/pantry/recipes` — search recipes where most ingredients match the pantry
- Simple: score recipes by how many ingredients the user already has, return top matches

### 5. Wire Mobile Screens to API
Replace all TODO/placeholder data in mobile screens with actual API calls:
- Home: fetch today's summary from `/api/v1/tracking/summary?date=today`
- Plan: fetch current plan from `/api/v1/meal-plans/current`
- Track: log food via `/api/v1/tracking/food`, water via `/api/v1/tracking/water`
- Recipes: search via `/api/v1/recipes?q=...`
- Profile: fetch/update from `/api/v1/profile`

## Architecture Notes
- **Pattern**: Repository (interface + SQLite impl) → Service (business logic) → Handler (thin HTTP layer)
- **Auth**: JWT token in Authorization Bearer header, `handler.UserIDFromContext(ctx)` to get user ID
- **Data**: All JSON columns (allergens, categories, ingredients) stored as TEXT, marshaled/unmarshaled in repository layer
- **Module path**: `github.com/rainyroot/bitewise/backend`
- **Mobile path alias**: `@/` maps to project root (configured in tsconfig.json)

## Running
```bash
# Backend
export JWT_SECRET="your-secret-here"
cd backend && go run ./cmd/server/ -migrate  # first time
cd backend && go run ./cmd/server/           # start server on :8080

# Mobile
cd mobile && pnpm install && pnpm start
```

## Key Files
| Component | Files |
|-----------|-------|
| Domain types | `backend/internal/domain/*.go` |
| Repositories | `backend/internal/repository/*.go` |
| Services | `backend/internal/service/*.go` |
| Handlers | `backend/internal/handler/*.go` |
| Auth middleware | `backend/internal/handler/middleware.go` |
| Scraper | `backend/internal/scraper/chefkoch.go` |
| DB schema | `backend/migrations/001_initial.sql` |
| API client | `mobile/services/api.ts` |
| Types | `mobile/types/index.ts` |
| Auth hook | `mobile/hooks/useAuth.tsx` |
