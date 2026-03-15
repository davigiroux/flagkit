package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/davigiroux/flagkit/api/internal/cache"
	"github.com/davigiroux/flagkit/api/internal/config"
	"github.com/davigiroux/flagkit/api/internal/db"
	"github.com/davigiroux/flagkit/api/internal/handler"
	"github.com/davigiroux/flagkit/api/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Run migrations
	runMigrations(cfg.DatabaseURL)

	// Connect to database
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// Connect to Redis
	flagCache, err := cache.New(cfg.RedisURL, 30*time.Second)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}

	queries := db.NewQueries(pool)

	// Bootstrap API key if none exist
	bootstrapAPIKey(ctx, queries)

	// Handlers
	flagHandler := handler.NewFlagHandler(queries, flagCache)
	evalHandler := handler.NewEvalHandler(queries, flagCache)
	auditHandler := handler.NewAuditHandler(queries)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(cfg.CORSOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", handler.Health)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(queries))

		r.Route("/flags", func(r chi.Router) {
			r.Get("/", flagHandler.List)
			r.Post("/", flagHandler.Create)
			r.Get("/{key}", flagHandler.Get)
			r.Patch("/{key}", flagHandler.Update)
			r.Delete("/{key}", flagHandler.Delete)
			r.Post("/{key}/toggle", flagHandler.Toggle)
		})

		r.Get("/evaluate/{key}", evalHandler.Evaluate)
		r.Get("/audit", auditHandler.List)
	})

	addr := ":" + cfg.Port
	log.Printf("FlagKit API listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func runMigrations(databaseURL string) {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatalf("migration init: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration: %v", err)
	}
	log.Println("migrations applied")
}

func bootstrapAPIKey(ctx context.Context, queries *db.Queries) {
	count, err := queries.CountAPIKeys(ctx)
	if err != nil {
		log.Printf("warning: could not check api keys: %v", err)
		return
	}
	if count > 0 {
		return
	}

	token := generateToken()
	hash := middleware.HashToken(token)
	if err := queries.CreateAPIKey(ctx, hash); err != nil {
		log.Fatalf("bootstrap api key: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("  BOOTSTRAP API KEY (save this!):")
	fmt.Printf("  %s\n", token)
	fmt.Println("========================================")
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "fk_" + hex.EncodeToString(b)
}
