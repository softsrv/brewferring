package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gorm.io/datatypes"

	"github.com/softsrv/brewferring/internal/components"
	"github.com/softsrv/brewferring/internal/config"
	ctx "github.com/softsrv/brewferring/internal/context"
	"github.com/softsrv/brewferring/internal/database"
	"github.com/softsrv/brewferring/internal/middleware"
	"github.com/softsrv/brewferring/internal/models"
	"github.com/softsrv/brewferring/internal/templates"
	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal-sdk-go/option"
)

type Handlers struct {
	config *config.Config
}

func NewHandlers(cfg *config.Config) *Handlers {
	h := &Handlers{}
	h.config = cfg

	return h
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	token, _ := middleware.GetAccessTokenFromHeader(r)

	// Add access token to context so that home page navbar can show
	// if user is logged in
	ct1 := r.Context()
	if len(token) > 0 {
		ct1 = ctx.WithAccessToken(ct1, token)
	}
	component := templates.Home()
	component.Render(ct1, w)
}

func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	token, _ := middleware.GetAccessTokenFromHeader(r)

	// Add access token to context so that home page navbar can show
	// if user is logged in
	ct1 := r.Context()
	if len(token) > 0 {
		ct1 = ctx.WithAccessToken(ct1, token)
	}
	component := templates.About()
	component.Render(ct1, w)
}

func (h *Handlers) Products(w http.ResponseWriter, r *http.Request) {

	provider, ok := ctx.GetProvider(r.Context())
	if !ok {
		log.Printf("terminal provider not found")
		http.Error(w, "Failed to find terminal provider", http.StatusInternalServerError)
	}

	products, err := provider.ListProducts(r.Context())
	if err != nil {
		log.Printf("products error: %s", err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	templates.Products(products).Render(r.Context(), w)
}

func (h *Handlers) Profile(w http.ResponseWriter, r *http.Request) {
	// profile page contains mostly external data, including:
	// addresses, credit cards, orders

	provider, ok := ctx.GetProvider(r.Context())
	if !ok {
		log.Printf("terminal provider not found")
		http.Error(w, "Failed to find terminal provider", http.StatusInternalServerError)
	}
	profile, err := provider.GetProfile(r.Context())
	if err != nil {
		log.Printf("failed to fetch profile: %s", err)
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}
	orders, err := provider.ListOrders(r.Context())
	if err != nil {
		log.Printf("failed to fetch orders: %s", err)
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	cards, err := provider.ListCards(r.Context())
	if err != nil {
		log.Printf("failed to fetch cards: %s", err)
		http.Error(w, "Failed to fetch cards", http.StatusInternalServerError)
		return
	}

	addresses, err := provider.ListAddresses(r.Context())
	if err != nil {
		log.Printf("failed to fetch addresses: %s", err)
		http.Error(w, "Failed to fetch addresses", http.StatusInternalServerError)
		return
	}

	templates.ProfileView(profile, orders, cards, addresses).Render(r.Context(), w)
}

func (h *Handlers) Orders(w http.ResponseWriter, r *http.Request) {

	provider, ok := ctx.GetProvider(r.Context())
	if !ok {
		log.Printf("terminal provider not found")
		http.Error(w, "Failed to find terminal provider", http.StatusInternalServerError)
	}

	orders, err := provider.ListOrders(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	component := templates.Orders(orders)
	component.Render(r.Context(), w)
}

func (h *Handlers) Buffers(w http.ResponseWriter, r *http.Request) {

	user, ok := ctx.GetUser(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	buffers, err := database.GetBuffersByUserID(user.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	component := templates.Buffers(buffers)
	component.Render(r.Context(), w)
}

func (h *Handlers) CreateBuffer(w http.ResponseWriter, r *http.Request) {
	user, ok := ctx.GetUser(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("unable to parse form data: %s", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	threshold := r.FormValue("threshold")
	date := r.FormValue("date")
	sType := r.FormValue("type")

	if len(name) <= 0 || len(name) > 50 {
		log.Println("invalid device name specified")
		http.Error(w, "Invalid device name. Must be 1-50 characters", http.StatusBadRequest)
		return
	}

	token, err := database.GenerateDeviceToken()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	buffer := &models.Buffer{
		Name:   name,
		UserID: user.ID,
		Token:  token,
	}

	if sType == "date" {
		sDate, err := time.Parse(time.DateOnly, date)
		if err != nil {
			log.Printf("failed to create new buffer due to invalid date: %s", err)
			http.Error(w, "Invalid date provided", http.StatusBadRequest)
		}
		buffer.OrderDate = datatypes.Date(sDate)
	} else {
		// Check if device belongs to user

		thresh, err := strconv.ParseFloat(threshold, 64)
		if err != nil {
			http.Error(w, "Invalid buffer threshold provided", http.StatusBadRequest)
			return
		}
		buffer.Threshold = thresh
	}

	if err := database.CreateBuffer(buffer); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	component := components.CreateBufferResponseComponent(buffer)
	component.Render(r.Context(), w)
}

func (h *Handlers) DeleteBuffer(w http.ResponseWriter, r *http.Request) {
	user, ok := ctx.GetUser(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	bufferID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Invalid buffer ID", http.StatusBadRequest)
		return
	}

	// Check if buffer belongs to user
	var buffer models.Buffer
	if err := database.DB.First(&buffer, bufferID).Error; err != nil {
		http.Error(w, "Buffer not found", http.StatusNotFound)
		return
	}

	if buffer.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := database.DeleteBuffer(buffer.ID); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	// Redirect to Terminal.shop OAuth
	redirectURL := h.config.OAuth.Provider.AuthEndpoint + "?" +
		"client_id=" + h.config.OAuth.ClientID +
		"&redirect_uri=" + h.config.OAuth.RedirectURI +
		"&response_type=code" +
		"&scope=allyourbasearebelongtous"

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the access token cookie
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
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
	values.Set("redirect_uri", h.config.OAuth.RedirectURI)
	values.Set("client_id", h.config.OAuth.ClientID)
	values.Set("client_secret", h.config.OAuth.ClientSecret)

	// Make the token exchange request
	resp, err := http.PostForm(h.config.OAuth.Provider.TokenEndpoint, values)
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
	log.Printf("the raw response body: %s", string(body))

	// Parse the response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	log.Printf("the token response: %+v", tokenResp)
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

// GetDeviceData returns the device data for a specific device
func (h *Handlers) GetDeviceData(w http.ResponseWriter, r *http.Request) {
	user, ok := ctx.GetUser(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bufferID, err := strconv.ParseUint(r.URL.Query().Get("buffer_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	// Get the buffer to check ownership
	var buffer models.Buffer
	if err := database.DB.Preload("User").First(&buffer, bufferID).Error; err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	if buffer.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get device data
	data, err := database.GetDeviceDataByBufferID(uint(bufferID))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// CreateDeviceData creates new device data from an authenticated device
func (h *Handlers) CreateDeviceData(w http.ResponseWriter, r *http.Request) {
	buffer, ok := ctx.GetBuffer(r.Context())
	if !ok {
		http.Error(w, "No buffer found", http.StatusNotFound)
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
		BufferID: buffer.ID,
		Value:    req.Value,
	}

	if err := database.CreateDeviceData(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update last used timestamp
	if err := database.UpdateBufferTokenLastUsedAt(buffer); err != nil {
		log.Printf("Failed to update device token last used at: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
}
