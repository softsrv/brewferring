package middleware

import (
	"context"
	"log"
	"net/http"

	ctx "github.com/softsrv/brewferring/internal/context"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		if cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		log.Printf("got a cookie: %s", cookie.Value)
		// Add the access token to the request context
		reqWithToken := r.WithContext(
			context.WithValue(r.Context(), ctx.AccessTokenKey, cookie.Value),
		)

		next(w, reqWithToken)
	}
}
