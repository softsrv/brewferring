package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/softsrv/brewferring/internal/config"
	ctx "github.com/softsrv/brewferring/internal/context"
	"github.com/softsrv/brewferring/internal/database"
	"github.com/softsrv/brewferring/internal/models"
	"github.com/softsrv/brewferring/internal/templates"
	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal-sdk-go/option"
)

type Handlers struct {
	oauthConfig struct {
		Issuer        string   `json:"issuer"`
		AuthEndpoint  string   `json:"authorization_endpoint"`
		TokenEndpoint string   `json:"token_endpoint"`
		JWKSEndpoint  string   `json:"jwks_uri"`
		ResponseTypes []string `json:"response_types_supported"`
		ClientID      string
		ClientSecret  string
		RedirectURI   string
	}
	config *config.Config
}

func NewHandlers(cfg *config.Config) *Handlers {
	h := &Handlers{}

	// Initialize OAuth config
	h.oauthConfig.Issuer = "https://auth.terminal.shop"
	h.oauthConfig.AuthEndpoint = "https://auth.terminal.shop/authorize"
	h.oauthConfig.TokenEndpoint = "https://auth.terminal.shop/token"
	h.oauthConfig.JWKSEndpoint = "https://auth.terminal.shop/.well-known/jwks.json"
	h.oauthConfig.ResponseTypes = []string{"code", "token"}
	h.oauthConfig.ClientID = cfg.OAuth.ClientID
	h.oauthConfig.ClientSecret = cfg.OAuth.ClientSecret
	h.oauthConfig.RedirectURI = cfg.OAuth.RedirectURI

	return h
}

// getClient returns a new Terminal.shop client initialized with the access token from context
func (h *Handlers) getClient(r *http.Request) *terminal.Client {
	accessToken := ""
	if token := r.Context().Value(ctx.AccessTokenKey); token != nil {
		accessToken = token.(string)
	}
	return terminal.NewClient(option.WithBearerToken(accessToken))
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	component := templates.Home()
	component.Render(r.Context(), w)
}

func (h *Handlers) Dashboard(w http.ResponseWriter, r *http.Request) {
	component := templates.Dashboard()
	component.Render(r.Context(), w)
}

func (h *Handlers) Products(w http.ResponseWriter, r *http.Request) {
	client := h.getClient(r)
	products, err := client.Product.List(r.Context())
	if err != nil {
		log.Printf("products error: %s", err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	var templateProducts []templates.Product
	for _, p := range products.Data {
		templateProducts = append(templateProducts, templates.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       float64(p.Variants[0].Price) / 100, // Convert cents to dollars
		})
	}

	templates.Products(templateProducts).Render(r.Context(), w)
}

func (h *Handlers) Profile(w http.ResponseWriter, r *http.Request) {
	client := h.getClient(r)
	profile, err := client.Profile.Me(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}

	templateProfile := templates.Profile{
		ID:        profile.Data.User.ID,
		Email:     profile.Data.User.Email,
		Name:      profile.Data.User.Name,
		CreatedAt: time.Now().Format(time.RFC3339), // Since we don't have CreatedAt from the API
	}

	templates.ProfileView(templateProfile).Render(r.Context(), w)
}

func (h *Handlers) Orders(w http.ResponseWriter, r *http.Request) {
	client := h.getClient(r)
	orders, err := client.Order.List(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	var templateOrders []templates.Order
	for _, o := range orders.Data {
		var items []templates.OrderItem
		for _, i := range o.Items {
			items = append(items, templates.OrderItem{
				ProductName: i.Description,
				Quantity:    int(i.Quantity),
				Price:       float64(i.Amount),
			})
		}

		templateOrders = append(templateOrders, templates.Order{
			ID:        o.ID,
			Status:    o.Tracking.URL,
			Total:     float64(o.Amount.Subtotal),
			Items:     items,
			CreatedAt: time.Now().Format(time.RFC3339), // Since we don't have CreatedAt from the API
		})
	}

	component := templates.Orders(templateOrders)
	component.Render(r.Context(), w)
}

func (h *Handlers) Devices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, ok := ctx.GetAccessToken(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Get user from access token
	user, err := h.getUserFromToken(accessToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	devices, err := database.GetUserDevices(user.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	component := templates.Devices(devices)
	component.Render(r.Context(), w)
}

func (h *Handlers) CreateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, ok := ctx.GetAccessToken(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Get user from access token
	user, err := h.getUserFromToken(accessToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := database.GenerateDeviceToken()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	device := &models.Device{
		Name:   req.Name,
		UserID: user.ID,
	}

	if err := database.CreateDevice(device); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	deviceToken := &models.DeviceToken{
		Token:    token,
		DeviceID: device.ID,
	}

	if err := database.DB.Create(deviceToken).Error; err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handlers) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, ok := ctx.GetAccessToken(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Get user from access token
	user, err := h.getUserFromToken(accessToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	deviceID, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	// Check if device belongs to user
	var device models.Device
	if err := database.DB.First(&device, deviceID).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	if device.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := database.DeleteDevice(device.ID); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) Schedulers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, ok := ctx.GetAccessToken(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Get user from access token
	user, err := h.getUserFromToken(accessToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	schedulers, err := database.GetUserSchedulers(user.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	devices, err := database.GetUserDevices(user.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	component := templates.Schedulers(schedulers, devices)
	component.Render(r.Context(), w)
}

func (h *Handlers) CreateScheduler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, ok := ctx.GetAccessToken(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Get user from access token
	user, err := h.getUserFromToken(accessToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name      string  `json:"name"`
		Type      string  `json:"type"`
		DeviceID  uint    `json:"device_id"`
		Threshold float64 `json:"threshold,omitempty"`
		Date      string  `json:"date,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Type != "device" && req.Type != "date" {
		http.Error(w, "Invalid scheduler type", http.StatusBadRequest)
		return
	}

	if req.Type == "device" && req.Threshold == 0 {
		http.Error(w, "Threshold is required for device-based scheduling", http.StatusBadRequest)
		return
	}

	if req.Type == "date" && req.Date == "" {
		http.Error(w, "Date is required for date-based scheduling", http.StatusBadRequest)
		return
	}

	// Check if device belongs to user
	var device models.Device
	if err := database.DB.First(&device, req.DeviceID).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	if device.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	scheduler := &models.Scheduler{
		Name:      req.Name,
		UserID:    user.ID,
		DeviceID:  req.DeviceID,
		Type:      req.Type,
		Threshold: req.Threshold,
	}

	if req.Type == "date" {
		date, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		scheduler.Date = date
	}

	if err := database.CreateScheduler(scheduler); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) DeleteScheduler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken, ok := ctx.GetAccessToken(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Get user from access token
	user, err := h.getUserFromToken(accessToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	schedulerID, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid scheduler ID", http.StatusBadRequest)
		return
	}

	// Check if scheduler belongs to user
	var scheduler models.Scheduler
	if err := database.DB.First(&scheduler, schedulerID).Error; err != nil {
		http.Error(w, "Scheduler not found", http.StatusNotFound)
		return
	}

	if scheduler.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := database.DeleteScheduler(scheduler.ID); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	// Redirect to Terminal.shop OAuth
	redirectURL := h.oauthConfig.AuthEndpoint + "?" +
		"client_id=" + h.oauthConfig.ClientID +
		"&redirect_uri=" + h.oauthConfig.RedirectURI +
		"&response_type=code" +
		"&scope=openid profile email"

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the access token cookie
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handlers) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
		return
	}

	// Prepare the token exchange request
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", h.oauthConfig.RedirectURI)
	values.Set("client_id", h.oauthConfig.ClientID)
	values.Set("client_secret", h.oauthConfig.ClientSecret)

	// Make the token exchange request
	resp, err := http.PostForm(h.oauthConfig.TokenEndpoint, values)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read token response: %v", err)
		http.Error(w, "Failed to read token response", http.StatusInternalServerError)
		return
	}

	// Parse the response
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		log.Printf("Failed to parse token response: %v", err)
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	// Get user profile from Terminal.shop
	client := terminal.NewClient(option.WithBearerToken(tokenResp.AccessToken))
	profile, err := client.Profile.Me(r.Context())
	if err != nil {
		log.Printf("Failed to fetch profile: %v", err)
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}

	// Create or update user in database
	user := &models.User{
		TerminalID: profile.Data.User.ID,
		Email:      profile.Data.User.Email,
		Name:       profile.Data.User.Name,
	}

	existingUser, err := database.GetUserByTerminalID(user.TerminalID)
	if err != nil {
		if err := database.CreateUser(user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else {
		user.ID = existingUser.ID
		if err := database.DB.Save(user).Error; err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	// Set the access token cookie
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokenResp.AccessToken,
		Path:     "/",
		MaxAge:   86400, // 24 hours in seconds
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)

	// Redirect to products page
	http.Redirect(w, r, "/products", http.StatusTemporaryRedirect)
}

// getUserFromToken extracts the user information from the access token and returns the corresponding user from the database
func (h *Handlers) getUserFromToken(accessToken string) (*models.User, error) {
	// Create a new Terminal client with the access token
	client := terminal.NewClient(option.WithBearerToken(accessToken))

	// Get the user profile from Terminal.shop
	profile, err := client.Profile.Me(context.Background())
	if err != nil {
		return nil, err
	}

	// Get or create the user in our database
	user, err := database.GetUserByTerminalID(profile.Data.User.ID)
	if err != nil {
		// If user doesn't exist, create them
		user = &models.User{
			TerminalID: profile.Data.User.ID,
			Email:      profile.Data.User.Email,
			Name:       profile.Data.User.Name,
		}
		if err := database.CreateUser(user); err != nil {
			return nil, err
		}
	} else {
		// Update user information if needed
		user.Email = profile.Data.User.Email
		user.Name = profile.Data.User.Name
		if err := database.DB.Save(user).Error; err != nil {
			return nil, err
		}
	}

	return user, nil
}

// GetDeviceData returns the device data for a specific device
func (h *Handlers) GetDeviceData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceID, err := strconv.ParseUint(r.URL.Query().Get("device_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	// Get the device to check ownership
	var device models.Device
	if err := database.DB.Preload("User").First(&device, deviceID).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	// Check if user is authorized to access this device
	accessToken, ok := ctx.GetAccessToken(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.getUserFromToken(accessToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if device.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get device data
	data, err := database.GetDeviceDataByDeviceID(uint(deviceID))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// CreateDeviceData creates new device data from an authenticated device
func (h *Handlers) CreateDeviceData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get device token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "No authorization header", http.StatusUnauthorized)
		return
	}

	// Extract token from "Bearer <token>"
	if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}
	token := authHeader[7:]

	// Get device token from database
	deviceToken, err := database.GetDeviceToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check rate limiting
	isLimited, err := database.IsDeviceTokenRateLimited(deviceToken.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if isLimited {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Parse request body
	var req struct {
		Value float64 `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create device data
	data := &models.DeviceData{
		DeviceID:      deviceToken.DeviceID,
		DeviceTokenID: deviceToken.ID,
		Value:         req.Value,
		CreatedAt:     time.Now(),
	}

	if err := database.CreateDeviceData(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update last used timestamp
	if err := database.UpdateDeviceTokenLastUsedAt(deviceToken.ID); err != nil {
		log.Printf("Failed to update device token last used at: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
}
