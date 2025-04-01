package context

import (
	"context"

	"github.com/softsrv/brewferring/internal/models"
)

type contextKey string

const (
	// AccessTokenKey is the key used to store the access token in the request context
	AccessTokenKey contextKey = "access_token"
	deviceKey      contextKey = "device"
)

func WithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, AccessTokenKey, token)
}

func GetAccessToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(AccessTokenKey).(string)
	return token, ok
}

func WithDevice(ctx context.Context, device *models.Device) context.Context {
	return context.WithValue(ctx, deviceKey, device)
}

func GetDevice(ctx context.Context) (*models.Device, bool) {
	device, ok := ctx.Value(deviceKey).(*models.Device)
	return device, ok
}
