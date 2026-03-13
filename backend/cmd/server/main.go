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
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "modernc.org/sqlite"

	"github.com/rainyroot/bitewise/backend/internal/config"
	"github.com/rainyroot/bitewise/backend/internal/handler"
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

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	profileSvc := service.NewProfileService(userRepo)
	recipeSvc := service.NewRecipeService(recipeRepo)
	planSvc := service.NewMealPlanService(planRepo, recipeRepo, userRepo)
	trackingSvc := service.NewTrackingService(trackingRepo)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	profileH := handler.NewProfileHandler(profileSvc)
	recipeH := handler.NewRecipeHandler(recipeSvc)
	planH := handler.NewMealPlanHandler(planSvc)
	trackingH := handler.NewTrackingHandler(trackingSvc)

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
				r.Get("/favorites", recipeH.GetFavorites)
				r.Get("/{id}", recipeH.GetByID)
				r.Post("/{id}/favorite", recipeH.AddFavorite)
				r.Delete("/{id}/favorite", recipeH.RemoveFavorite)
			})

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

			// Shopping lists (Phase 2)
			r.Route("/shopping-lists", func(r chi.Router) {
				r.Post("/", placeholder("create shopping list"))
				r.Get("/current", placeholder("get current shopping list"))
				r.Patch("/items/{id}", placeholder("toggle shopping item"))
			})

			// Achievements (Phase 3)
			r.Route("/achievements", func(r chi.Router) {
				r.Get("/", placeholder("list achievements"))
				r.Get("/mine", placeholder("list my achievements"))
			})
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

func placeholder(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.JSON(w, http.StatusNotImplemented, map[string]string{
			"message": fmt.Sprintf("%s: not yet implemented", name),
		})
	}
}

func runMigrations(db *sql.DB) error {
	migration, err := os.ReadFile("migrations/001_initial.sql")
	if err != nil {
		return fmt.Errorf("reading migration file: %w", err)
	}

	if _, err := db.Exec(string(migration)); err != nil {
		return fmt.Errorf("executing migration: %w", err)
	}

	return nil
}
