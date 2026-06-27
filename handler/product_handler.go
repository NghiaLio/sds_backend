package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mobile-crud-backend/model"
	"mobile-crud-backend/service"
)

// ProductHandler handles HTTP requests for products.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// RegisterRoutes registers the handlers to the provided ServeMux.
func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux, auth *AuthMiddleware) {
	mux.Handle("POST /products", auth.Handler(http.HandlerFunc(h.CreateProduct)))
	mux.Handle("GET /products", auth.Handler(http.HandlerFunc(h.GetAllProducts)))
	mux.Handle("PUT /products/{id}", auth.Handler(http.HandlerFunc(h.UpdateProduct)))
	mux.Handle("PATCH /products/{id}", auth.Handler(http.HandlerFunc(h.PatchProduct)))
	mux.Handle("DELETE /products/{id}", auth.Handler(http.HandlerFunc(h.DeleteProduct)))
}

// CreateProduct handles product creation.
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.CreateProduct(r.Context(), &product); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithSuccess(w, http.StatusCreated, product)
}

// GetAllProducts handles listing products with optional filters and pagination.
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category_id")
	var categoryID int64
	if categoryIDStr != "" {
		categoryID, _ = strconv.ParseInt(categoryIDStr, 10, 64)
	}

	keyword := r.URL.Query().Get("keyword")

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	products, _, err := h.service.GetAllProducts(r.Context(), categoryID, keyword, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products: "+err.Error())
		return
	}

	respondWithSuccess(w, http.StatusOK, products)
}

// UpdateProduct handles product updates.
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.UpdateProduct(r.Context(), id, &product); err != nil {
		if err.Error() == "product not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondWithSuccess(w, http.StatusOK, product)
}

// PatchProduct handles product partial updating.
func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	existing, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(existing); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.UpdateProduct(r.Context(), id, existing); err != nil {
		if err.Error() == "product not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respondWithSuccess(w, http.StatusOK, existing)
}

// DeleteProduct handles deleting a product.
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if err.Error() == "product not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithSuccess(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
