package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mobile-crud-backend/model"
	"mobile-crud-backend/service"
)

// CategoryHandler handles HTTP requests for product categories.
type CategoryHandler struct {
	service service.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler instance.
func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// RegisterRoutes registers the handlers to the provided ServeMux.
func (h *CategoryHandler) RegisterRoutes(mux *http.ServeMux, auth *AuthMiddleware) {
	mux.Handle("POST /categories", auth.Handler(http.HandlerFunc(h.CreateCategory)))
	mux.Handle("GET /categories", auth.Handler(http.HandlerFunc(h.GetAllCategories)))
	mux.Handle("PUT /categories/{id}", auth.Handler(http.HandlerFunc(h.UpdateCategory)))
	mux.Handle("DELETE /categories/{id}", auth.Handler(http.HandlerFunc(h.DeleteCategory)))
}

// CreateCategory handles category creation.
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.CreateCategory(r.Context(), &category); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, category)
}

// GetAllCategories handles listing all categories.
func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch categories: "+err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, categories)
}

// UpdateCategory handles category updating.
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	var category model.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.UpdateCategory(r.Context(), id, &category); err != nil {
		if err.Error() == "category not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondWithSuccess(w, http.StatusOK, category)
}

// DeleteCategory handles category deleting.
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID format")
		return
	}

	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		if err.Error() == "category not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithSuccess(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
