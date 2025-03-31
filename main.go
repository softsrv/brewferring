package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/softsrv/brewferring/internal/config"
	"github.com/softsrv/brewferring/internal/handlers"
	"github.com/softsrv/brewferring/internal/middleware"
)

// content holds our static web server content.
//
//go:embed static/*
var content embed.FS

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize handlers
	h := handlers.NewHandlers(cfg)

	// Create router
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/dashboard", middleware.AuthMiddleware(h.Dashboard))
	mux.HandleFunc("/products", middleware.AuthMiddleware(h.Products))
	mux.HandleFunc("/profile", middleware.AuthMiddleware(h.Profile))
	mux.HandleFunc("/orders", middleware.AuthMiddleware(h.Orders))

	// Auth routes
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/callback", h.OAuthCallback)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
