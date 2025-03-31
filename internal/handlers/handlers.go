package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/softsrv/brewferring/internal/config"
	ctx "github.com/softsrv/brewferring/internal/context"
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
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
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
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
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
