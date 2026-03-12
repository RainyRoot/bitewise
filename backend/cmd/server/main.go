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
		// Auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", placeholder("register"))
			r.Post("/login", placeholder("login"))
		})

		// Protected routes (will add auth middleware later)
		r.Group(func(r chi.Router) {
			// User profile
			r.Route("/profile", func(r chi.Router) {
				r.Get("/", placeholder("get profile"))
				r.Put("/", placeholder("update profile"))
				r.Get("/allergies", placeholder("get allergies"))
				r.Put("/allergies", placeholder("set allergies"))
				r.Get("/preferences", placeholder("get preferences"))
				r.Put("/preferences", placeholder("set preferences"))
			})

			// Recipes
			r.Route("/recipes", func(r chi.Router) {
				r.Get("/", placeholder("search recipes"))
				r.Get("/{id}", placeholder("get recipe"))
				r.Post("/{id}/favorite", placeholder("favorite recipe"))
				r.Delete("/{id}/favorite", placeholder("unfavorite recipe"))
				r.Get("/favorites", placeholder("list favorites"))
			})

			// Meal plans
			r.Route("/meal-plans", func(r chi.Router) {
				r.Post("/generate", placeholder("generate meal plan"))
				r.Get("/current", placeholder("get current meal plan"))
				r.Get("/{id}", placeholder("get meal plan"))
				r.Patch("/{id}/entries/{entryId}", placeholder("update meal plan entry"))
				r.Post("/{id}/activate", placeholder("activate meal plan"))
			})

			// Food tracking
			r.Route("/tracking", func(r chi.Router) {
				r.Post("/food", placeholder("log food"))
				r.Get("/food", placeholder("get food logs"))
				r.Delete("/food/{id}", placeholder("delete food log"))
				r.Post("/water", placeholder("log water"))
				r.Get("/water", placeholder("get water logs"))
				r.Get("/summary", placeholder("get nutrition summary"))
			})

			// Shopping lists
			r.Route("/shopping-lists", func(r chi.Router) {
				r.Post("/", placeholder("create shopping list"))
				r.Get("/current", placeholder("get current shopping list"))
				r.Patch("/items/{id}", placeholder("toggle shopping item"))
			})

			// Achievements
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
