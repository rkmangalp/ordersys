package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rkmangalp/ordersys/internal/database"
	"github.com/rkmangalp/ordersys/internal/product"
)

func main() {
	// -------------------------------------------------------------------------
	// 1. DATABASE CONNECTION
	// -------------------------------------------------------------------------
	db, err := database.NewConnection(database.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected")

	// -------------------------------------------------------------------------
	// 2. DEPENDENCY INJECTION (Repository -> Service -> Handler)
	// -------------------------------------------------------------------------
	productRepo := product.NewMySQLRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// -------------------------------------------------------------------------
	// 3. ROUTER SETUP (chi)
	// -------------------------------------------------------------------------
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	// Health endpoint (no API version)
	r.Get("/health", healthHandler)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		productHandler.RegisterRoutes(r)
	})

	// -------------------------------------------------------------------------
	// 4. HTTP SERVER CONFIG
	// -------------------------------------------------------------------------
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// -------------------------------------------------------------------------
	// 5. GRACEFUL SHUTDOWN
	// -------------------------------------------------------------------------
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// healthHandler handles GET /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Simple health response
	_, _ = w.Write([]byte(`{"status":"ok","time":"` + time.Now().Format(time.RFC3339) + `"}`))
}