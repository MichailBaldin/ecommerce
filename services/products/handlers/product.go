package handlers

import (
	"ecommerce/services/products/models"
	"ecommerce/services/products/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ProductHandler struct {
	service service.ProductService
	logger  *zap.Logger
}

func NewProductHandler(service service.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid JSON in create product", zap.Error(err))
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	product, err := h.service.CreateProduct(&req)
	if err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	h.sendJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := h.extractID(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		h.logger.Error("Failed to get product", zap.Int("id", id), zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "Failed to get product")
		return
	}

	if product == nil {
		h.sendError(w, http.StatusNotFound, "Product not found")
		return
	}

	h.sendJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) extractID(r *http.Request) (int, error) {
	path := r.URL.Path
	// Expect path like "/api/v1/products/123"
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return 0, fmt.Errorf("invalid path")
	}

	idStr := parts[len(parts)-1] // Last part should be ID
	return strconv.Atoi(idStr)
}

func (h *ProductHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON", zap.Error(err))
	}
}

func (h *ProductHandler) sendError(w http.ResponseWriter, status int, message string) {
	response := models.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}
	h.sendJSON(w, status, response)
}
