package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/samaita/quick-go/internal/auth"
	"github.com/samaita/quick-go/internal/config"
	"github.com/samaita/quick-go/internal/db"
	appMiddleware "github.com/samaita/quick-go/internal/middleware"
)

func main() {
	cfg := config.Load()

	// Database
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("server: db open: %v", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatalf("server: db migrate: %v", err)
	}

	// Redis session store
	sessions := auth.NewSessionStore(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := sessions.Ping(context.Background()); err != nil {
		log.Printf("server: WARNING: redis not reachable (%v) — session revocation disabled", err)
	}

	// Auth handler
	authHandler := auth.NewHandler(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
		cfg.JWTSecret,
		cfg.AESKey,
		cfg.AppURL,
		cfg.CookieDomain,
		cfg.IsProd,
		sessions,
		database,
	)

	// Router
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Auth routes (public)
	r.Get("/auth/google", authHandler.GoogleLogin)
	r.Get("/auth/google/callback", authHandler.GoogleCallback)
	r.Post("/auth/logout", authHandler.Logout)
	r.Get("/auth/logout", authHandler.Logout) // allow GET for nav link

	// /auth/me — protected, used by Alpine.js to check login state
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.JWT(cfg.JWTSecret, sessions))
		r.Get("/auth/me", authHandler.Me)
	})

	// Protected API routes — generated CRUD goes here
	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.JWT(cfg.JWTSecret, sessions))
		// generated.RegisterRoutes(r, database)
		// ↑ Uncomment after running: quickgen --schema example.sql
	})

	// Static files: serve Hugo-built public/ from filesystem
	// Run `hugo` in frontend/hugo-site/ before starting the server.
	staticDir := cfg.StaticDir
	if staticDir == "" {
		staticDir = defaultStaticDir()
	}
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("server: WARNING: static dir %q not found — run `hugo` in frontend/hugo-site/ first", staticDir)
	}
	r.Handle("/*", http.FileServer(http.Dir(staticDir)))

	addr := ":" + cfg.Port
	log.Printf("server: listening on %s (static: %s)", addr, staticDir)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}

// defaultStaticDir resolves frontend/hugo-site/public relative to the binary's location.
func defaultStaticDir() string {
	exe, err := os.Executable()
	if err != nil {
		return filepath.Join("frontend", "hugo-site", "public")
	}
	// Walk up from the binary to the repo root (up to 5 levels)
	dir := filepath.Dir(exe)
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, "frontend", "hugo-site", "public")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	return filepath.Join("frontend", "hugo-site", "public")
}
