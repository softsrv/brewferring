package context

type contextKey string

const (
	// AccessTokenKey is the key used to store the access token in the request context
	AccessTokenKey contextKey = "access_token"
)
