package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "modernc.org/sqlite"

	"github.com/rainyroot/bitewise/backend/internal/config"
	"github.com/rainyroot/bitewise/backend/internal/handler"
	"github.com/rainyroot/bitewise/backend/internal/nutrition"
	"github.com/rainyroot/bitewise/backend/internal/repository"
	"github.com/rainyroot/bitewise/backend/internal/service"
	"github.com/rainyroot/bitewise/backend/pkg/httputil"
)

func main() {
	migrate := flag.Bool("migrate", false, "run database migrations and exit")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := sql.Open("sqlite", cfg.DatabasePath+"?_pragma=journal_mode(wal)&_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("database open: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database ping: %v", err)
	}

	if *migrate {
		if err := runMigrations(db); err != nil {
			log.Fatalf("migrations: %v", err)
		}
		log.Println("migrations completed successfully")
		return
	}

	// Repositories
	userRepo := repository.NewSQLiteUserRepository(db)
	recipeRepo := repository.NewSQLiteRecipeRepository(db)
	planRepo := repository.NewSQLiteMealPlanRepository(db)
	trackingRepo := repository.NewSQLiteTrackingRepository(db)
	shoppingRepo := repository.NewSQLiteShoppingListRepository(db)
	pantryRepo := repository.NewSQLitePantryRepository(db)
	achievementRepo := repository.NewSQLiteAchievementRepository(db)
	sharingRepo := repository.NewSQLiteSharingRepository(db)
	notifRepo := repository.NewSQLiteNotificationRepository(db)
	statsRepo := repository.NewSQLiteStatsRepository(db)
	diaryRepo := repository.NewSQLiteDiaryRepository(db)
	priceRepo := repository.NewSQLitePriceRepository(db)
	friendRepo := repository.NewSQLiteFriendRepository(db)

	// External providers
	nutritionProvider := nutrition.NewOpenFoodFactsProvider()

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	profileSvc := service.NewProfileService(userRepo)
	recipeSvc := service.NewRecipeService(recipeRepo)
	planSvc := service.NewMealPlanService(planRepo, recipeRepo, userRepo)
	trackingSvc := service.NewTrackingService(trackingRepo)
	shoppingSvc := service.NewShoppingService(shoppingRepo, planRepo)
	pantrySvc := service.NewPantryService(pantryRepo, recipeRepo)
	achievementSvc := service.NewAchievementService(achievementRepo)
	sharingSvc := service.NewSharingService(sharingRepo, recipeRepo)
	notifSvc := service.NewNotificationService(notifRepo)
	statsSvc := service.NewStatsService(statsRepo)
	diarySvc := service.NewDiaryService(diaryRepo)
	priceSvc := service.NewPriceService(priceRepo)
	friendSvc := service.NewFriendService(friendRepo, userRepo)
	exportSvc := service.NewExportService(trackingRepo, userRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	profileH := handler.NewProfileHandler(profileSvc)
	recipeH := handler.NewRecipeHandler(recipeSvc)
	planH := handler.NewMealPlanHandler(planSvc)
	trackingH := handler.NewTrackingHandler(trackingSvc)
	shoppingH := handler.NewShoppingHandler(shoppingSvc)
	pantryH := handler.NewPantryHandler(pantrySvc)
	nutritionH := handler.NewNutritionHandler(nutritionProvider)
	achievementH := handler.NewAchievementHandler(achievementSvc)
	sharingH := handler.NewSharingHandler(sharingSvc)
	notifH := handler.NewNotificationHandler(notifSvc)
	statsH := handler.NewStatsHandler(statsSvc)
	diaryH := handler.NewDiaryHandler(diarySvc)
	priceH := handler.NewPriceHandler(priceSvc)
	friendH := handler.NewFriendHandler(friendSvc)
	exportH := handler.NewExportHandler(exportSvc)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
		})

		// Nutrition lookup (public)
		r.Route("/nutrition", func(r chi.Router) {
			r.Get("/barcode/{code}", nutritionH.LookupBarcode)
			r.Get("/search", nutritionH.SearchFood)
		})

		// Seasonal calendar (public)
		r.Get("/seasonal", nutritionH.GetSeasonal)

		// Shared recipes (public view)
		r.Get("/shared/{code}", sharingH.GetSharedRecipe)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware(authSvc))

			// User profile
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", profileH.GetProfile)
				r.Put("/", profileH.UpdateProfile)
				r.Get("/allergies", profileH.GetAllergies)
				r.Put("/allergies", profileH.SetAllergies)
				r.Get("/preferences", profileH.GetPreferences)
				r.Put("/preferences", profileH.SetPreferences)
			})

			// Recipes
			r.Route("/recipes", func(r chi.Router) {
				r.Get("/", recipeH.Search)
				r.Post("/", recipeH.CreateRecipe)
				r.Get("/favorites", recipeH.GetFavorites)
				r.Get("/mine", recipeH.GetMyRecipes)
				r.Get("/{id}", recipeH.GetByID)
				r.Delete("/{id}", recipeH.DeleteRecipe)
				r.Post("/{id}/favorite", recipeH.AddFavorite)
				r.Delete("/{id}/favorite", recipeH.RemoveFavorite)
				r.Post("/{id}/share", sharingH.ShareRecipe)
			})

			// Shared recipe import
			r.Post("/shared/{code}/save", sharingH.SaveSharedRecipe)

			// Meal plans
			r.Route("/meal-plans", func(r chi.Router) {
				r.Post("/generate", planH.Generate)
				r.Get("/current", planH.GetCurrent)
				r.Get("/{id}", planH.GetByID)
				r.Patch("/{id}/entries/{entryId}", planH.UpdateEntry)
				r.Post("/{id}/activate", planH.Activate)
			})

			// Food tracking
			r.Route("/tracking", func(r chi.Router) {
				r.Post("/food", trackingH.LogFood)
				r.Get("/food", trackingH.GetFoodLogs)
				r.Delete("/food/{id}", trackingH.DeleteFoodLog)
				r.Post("/water", trackingH.LogWater)
				r.Get("/water", trackingH.GetWaterLogs)
				r.Get("/summary", trackingH.GetNutritionSummary)
			})

			// Shopping lists
			r.Route("/shopping-lists", func(r chi.Router) {
				r.Post("/", shoppingH.GenerateFromMealPlan)
				r.Get("/current", shoppingH.GetCurrent)
				r.Patch("/items/{id}", shoppingH.ToggleItem)
			})

			// Pantry / leftover recipes
			r.Route("/pantry", func(r chi.Router) {
				r.Post("/", pantryH.SetItems)
				r.Get("/", pantryH.GetItems)
				r.Get("/recipes", pantryH.FindRecipes)
			})

			// Achievements
			r.Route("/achievements", func(r chi.Router) {
				r.Get("/", achievementH.GetAll)
				r.Get("/mine", achievementH.GetMine)
			})

			// Notifications
			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", notifH.GetSettings)
				r.Put("/", notifH.UpdateSettings)
			})

			// Statistics
			r.Route("/stats", func(r chi.Router) {
				r.Get("/weekly", statsH.GetWeeklyStats)
				r.Get("/monthly", statsH.GetMonthlyStats)
				r.Get("/streaks", statsH.GetStreaks)
			})

			// Diary
			r.Route("/diary", func(r chi.Router) {
				r.Post("/", diaryH.CreateOrUpdate)
				r.Get("/", diaryH.GetByDate)
				r.Get("/monthly", diaryH.GetMonthly)
				r.Delete("/{id}", diaryH.Delete)
			})

			// Prices
			r.Route("/prices", func(r chi.Router) {
				r.Post("/", priceH.LogPrice)
				r.Get("/", priceH.GetLogs)
				r.Get("/trends", priceH.GetTrend)
				r.Get("/compare", priceH.CompareStores)
				r.Get("/spending", priceH.GetSpending)
			})

			// Friends
			r.Route("/friends", func(r chi.Router) {
				r.Post("/invite", friendH.InviteFriend)
				r.Get("/invites", friendH.GetPendingInvites)
				r.Post("/invites/{id}", friendH.RespondToInvite)
				r.Get("/", friendH.GetFriends)
				r.Delete("/{id}", friendH.RemoveFriend)
			})

			// Leaderboard
			r.Get("/leaderboard", friendH.GetLeaderboard)

			// Export & Account
			r.Route("/export", func(r chi.Router) {
				r.Get("/csv", exportH.ExportCSV)
				r.Get("/json", exportH.ExportJSON)
			})
			r.Delete("/account", exportH.DeleteAccount)
		})
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("BiteWise API server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-done
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}

	log.Println("server stopped")
}

func runMigrations(db *sql.DB) error {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return fmt.Errorf("finding migration files: %w", err)
	}

	sort.Strings(files)

	for _, file := range files {
		migration, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading %s: %w", file, err)
		}

		if _, err := db.Exec(string(migration)); err != nil {
			return fmt.Errorf("executing %s: %w", file, err)
		}
		log.Printf("applied migration: %s", file)
	}

	return nil
}
