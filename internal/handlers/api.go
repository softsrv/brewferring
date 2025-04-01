package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/softsrv/brewferring/internal/context"
	"github.com/softsrv/brewferring/internal/database"
	"github.com/softsrv/brewferring/internal/models"
)

type DeviceDataRequest struct {
	Value float64 `json:"value"`
}

func (h *Handlers) authenticateDeviceToken(r *http.Request) (*models.Device, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, nil
	}

	token, err := database.GetDeviceToken(parts[1])
	if err != nil {
		return nil, err
	}

	// Check rate limiting
	rateLimited, err := database.IsDeviceTokenRateLimited(token.ID)
	if err != nil {
		return nil, err
	}
	if rateLimited {
		return nil, nil
	}

	// Update last used time
	err = database.UpdateDeviceTokenLastUsedAt(token.ID)
	if err != nil {
		return nil, err
	}

	var device models.Device
	err = database.DB.First(&device, token.DeviceID).Error
	if err != nil {
		return nil, err
	}

	return &device, nil
}

func DeviceData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	device, ok := context.GetDevice(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req DeviceDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the device token from the request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := database.GetDeviceToken(authHeader[7:]) // Remove "Bearer " prefix
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update last used time
	token.LastUsedAt = time.Now()
	if err := database.DB.Save(token).Error; err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Save device data
	data := &models.DeviceData{
		DeviceID:      device.ID,
		DeviceTokenID: token.ID,
		Value:         req.Value,
		CreatedAt:     time.Now(),
	}

	if err := database.SaveDeviceData(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
