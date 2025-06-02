package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/softsrv/brewferring/internal/config"
	"github.com/softsrv/brewferring/internal/database"
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
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create handlers
	h := handlers.NewHandlers(cfg)

	// Create server
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.FileServer(http.FS(content)))

	// Public routes
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/callback", h.OAuthCallback)
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/about", h.About)

	// Protected routes
	mux.Handle("GET /products", middleware.Auth(http.HandlerFunc(h.Products)))
	mux.Handle("GET /orders", middleware.Auth(http.HandlerFunc(h.Orders)))
	mux.Handle("GET /profile", middleware.Auth(http.HandlerFunc(h.Profile)))

	// buffers
	mux.Handle("GET /buffers", middleware.Auth(http.HandlerFunc(h.Buffers)))
	mux.Handle("POST /buffers", middleware.Auth(http.HandlerFunc(h.CreateBuffer)))
	mux.Handle("DELETE /buffers/{id}", middleware.Auth(http.HandlerFunc(h.DeleteBuffer)))

	// API routes
	mux.Handle("POST /api/devices/data", middleware.DeviceAuth(http.HandlerFunc(h.CreateDeviceData)))

	// Start server
	log.Println("Starting server on", cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), mux); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
