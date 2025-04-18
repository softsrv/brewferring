package context

import (
	"context"

	"github.com/softsrv/brewferring/internal/models"
	"github.com/terminaldotshop/terminal-sdk-go"
)

type contextKey string

const (
	// AccessTokenKey is the key used to store the access token in the request context
	AccessTokenKey contextKey = "access_token"
	UserKey        contextKey = "user_key"
	ClientKey      contextKey = "client_key"
	DeviceKey      contextKey = "device"
)

func WithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AccessTokenKey, token)
}

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func WithTerminalClient(ctx context.Context, client *terminal.Client) context.Context {
	return context.WithValue(ctx, ClientKey, client)
}

func GetAccessToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(AccessTokenKey).(string)
	return token, ok
}

func HasAccessTokenValue(ctx context.Context) bool {
	_, ok := ctx.Value(AccessTokenKey).(string)
	return ok
}

func WithDevice(ctx context.Context, device *models.Device) context.Context {
	return context.WithValue(ctx, DeviceKey, device)
}

func GetDevice(ctx context.Context) (*models.Device, bool) {
	device, ok := ctx.Value(DeviceKey).(*models.Device)
	return device, ok
}
