package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/softsrv/brewferring/internal/context"
	"github.com/softsrv/brewferring/internal/database"
	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal-sdk-go/option"
)

// getClient returns a new Terminal.shop client initialized with the access token from context
func getClient(token string) *terminal.Client {
	return terminal.NewClient(option.WithBearerToken(token))
}

func GetAccessTokenFromHeader(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetAccessTokenFromHeader(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add access token to context
		ctx := context.WithAccessToken(r.Context(), token)
		// Add terminal client to context
		client := getClient(token)
		ctx = context.WithTerminalClient(ctx, client)
		// Add user to context
		terminalUser, err := client.Profile.Me(ctx)
		if err != nil {
			log.Println("Unable to get terminal user")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		user, err := database.GetUserByTerminalID(terminalUser.Data.User.ID)
		if err != nil {
			log.Println("Unable to find user in DB for current session")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		ctx = context.WithUser(ctx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DeviceAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		device, err := database.GetDeviceByToken(parts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check rate limiting
		rateLimited, err := database.CheckIsTokenLimited(device.TokenLastUsedAt)
		if err != nil || rateLimited {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Add device to context
		ctx := context.WithDevice(r.Context(), &device)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
