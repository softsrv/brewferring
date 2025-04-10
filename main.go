package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/softsrv/brewferring/internal/config"
	"github.com/softsrv/brewferring/internal/database"
	"github.com/softsrv/brewferring/internal/handlers"
	"github.com/softsrv/brewferring/internal/middleware"
)

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
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public routes
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/callback", h.OAuthCallback)

	// Protected routes
	mux.HandleFunc("/logout", h.Logout)
	mux.Handle("GET /dashboard", middleware.Auth(http.HandlerFunc(h.Dashboard)))
	mux.Handle("GET /products", middleware.Auth(http.HandlerFunc(h.Products)))
	mux.Handle("GET /orders", middleware.Auth(http.HandlerFunc(h.Orders)))
	mux.Handle("GET /profile", middleware.Auth(http.HandlerFunc(h.Profile)))
	mux.Handle("GET /devices", middleware.Auth(http.HandlerFunc(h.Devices)))
	mux.Handle("GET /schedulers", middleware.Auth(http.HandlerFunc(h.Schedulers)))
	mux.Handle("POST /devices", middleware.Auth(http.HandlerFunc(h.CreateDevice)))
	mux.Handle("DELETE /devices/{id}", middleware.Auth(http.HandlerFunc(h.DeleteDevice)))
	mux.Handle("POST /schedulers", middleware.Auth(http.HandlerFunc(h.CreateScheduler)))
	mux.Handle("DELETE /schedulers/{id}", middleware.Auth(http.HandlerFunc(h.DeleteScheduler)))

	// API routes
	mux.Handle("/api/devices/data", middleware.DeviceAuth(http.HandlerFunc(h.CreateDeviceData)))

	// Start server
	log.Println("Starting server on", cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), mux); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
