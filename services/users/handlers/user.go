package handlers

import (
	"ecommerce/services/users/models"
	"ecommerce/services/users/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in create user", zap.Error(err))
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.sendJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := h.extractID(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUser(id)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Int("id", id), zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user == nil {
		h.sendError(w, http.StatusNotFound, "User not found")
		return
	}

	h.sendJSON(w, http.StatusOK, user)
}

func (h *UserHandler) extractID(r *http.Request) (int, error) {
	path := r.URL.Path
	// Expect path like "/api/v1/users/123"
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return 0, fmt.Errorf("invalid path")
	}

	idStr := parts[len(parts)-1] // Last part should be ID
	return strconv.Atoi(idStr)
}

func (h *UserHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON", zap.Error(err))
	}
}

func (h *UserHandler) sendError(w http.ResponseWriter, status int, message string) {
	response := models.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}
	h.sendJSON(w, status, response)
}
