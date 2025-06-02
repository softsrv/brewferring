package context

import (
	"context"

	"github.com/softsrv/brewferring/internal/models"
	"github.com/softsrv/brewferring/internal/provider"
	"github.com/terminaldotshop/terminal-sdk-go"
)

type contextKey string

const (
	// AccessTokenKey is the key used to store the access token in the request context
	AccessTokenKey contextKey = "access_token"
	UserKey        contextKey = "user_key"
	ClientKey      contextKey = "client_key"
	BufferKey      contextKey = "buffer"
	ProviderKey    contextKey = "provider"
)

func WithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AccessTokenKey, token)
}

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func WithBuffer(ctx context.Context, buffer *models.Buffer) context.Context {
	return context.WithValue(ctx, BufferKey, buffer)
}

func WithTerminalClient(ctx context.Context, client *terminal.Client) context.Context {
	return context.WithValue(ctx, ClientKey, client)
}
func WithProvider(ctx context.Context, provider provider.Provider) context.Context {
	return context.WithValue(ctx, ProviderKey, provider)
}

func HasAccessTokenValue(ctx context.Context) bool {
	_, ok := ctx.Value(AccessTokenKey).(string)
	return ok
}

func GetAccessToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(AccessTokenKey).(string)
	return token, ok
}

func GetTerminalClient(ctx context.Context) (*terminal.Client, bool) {
	client, ok := ctx.Value(ClientKey).(*terminal.Client)
	return client, ok
}

func GetProvider(ctx context.Context) (provider.Provider, bool) {
	provider, ok := ctx.Value(ProviderKey).(provider.Provider)
	return provider, ok
}

func GetUser(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserKey).(*models.User)
	return user, ok
}
func GetBuffer(ctx context.Context) (*models.Buffer, bool) {
	buffer, ok := ctx.Value(BufferKey).(*models.Buffer)
	return buffer, ok
}
