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
	SchedulerKey   contextKey = "scheduler"
	ProviderKey    contextKey = "provider"
)

func WithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AccessTokenKey, token)
}

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func WithScheduler(ctx context.Context, scheduler *models.Scheduler) context.Context {
	return context.WithValue(ctx, SchedulerKey, scheduler)
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
func GetScheduler(ctx context.Context) (*models.Scheduler, bool) {
	scheduler, ok := ctx.Value(SchedulerKey).(*models.Scheduler)
	return scheduler, ok
}
