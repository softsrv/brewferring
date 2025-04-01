package middleware

import (
	"net/http"
	"strings"

	"github.com/softsrv/brewferring/internal/context"
	"github.com/softsrv/brewferring/internal/database"
	"github.com/softsrv/brewferring/internal/models"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// Add access token to context
		ctx := context.WithAccessToken(r.Context(), cookie.Value)
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

		token, err := database.GetDeviceToken(parts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check rate limiting
		rateLimited, err := database.IsDeviceTokenRateLimited(token.ID)
		if err != nil || rateLimited {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		var device models.Device
		err = database.DB.First(&device, token.DeviceID).Error
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add device to context
		ctx := context.WithDevice(r.Context(), &device)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
