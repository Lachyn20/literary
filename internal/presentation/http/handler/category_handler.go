package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/pagination"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/category"
)

type CategoryHandler struct {
	svc *category.CategoryService
}

func NewCategoryHandler(svc *category.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/categories", h.List)
	r.Get("/api/categories/{id}", h.Get)
	r.Post("/api/categories", h.Create)
	r.Put("/api/categories/{id}", h.Update)
	r.Delete("/api/categories/{id}", h.Delete)
}

// @Summary List categories
// @Description List all categories with pagination support
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param page query int false "Page number (1-indexed, alternative to offset)"
// @Success 200 {array} dto.CategoryResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/categories [get]
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	paginationParams := pagination.Parse(r)
	categories, total, err := h.svc.List(r.Context(), paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	paginationInfo := pagination.NewInfo(paginationParams.Limit, paginationParams.Offset, total)
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "ok",
		"data":         categoryResponses(categories),
		"pagination":   paginationInfo,
	})
}

// @Summary Get category
// @Description Get a single category by id
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	category, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, categoryResponse(category))
}

// @Summary Create category
// @Description Create a new category
// @Accept json
// @Param payload body dto.CategoryCreateRequest true "category payload"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := validation.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	category := &entity.Category{
		ID:   uuid.New(),
		Name: req.Name,
	}
	if err := h.svc.Create(r.Context(), category); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, categoryResponse(category))
}

// @Summary Update category
// @Description Update an existing category by id
// @Accept json
// @Param id path string true "Category ID"
// @Param payload body dto.CategoryCreateRequest true "category payload"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.CategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	// allow partial update: fetch existing and overlay
	old, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if req.Name != "" {
		old.Name = req.Name
	}
	if err := h.svc.Update(r.Context(), old); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, categoryResponse(old))
}

// @Summary Delete category
// @Description Delete an existing category by id
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
