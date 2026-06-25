package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/category"
)

type CategoryHandler struct {
	createUC *category.CreateCategoryUseCase
	getUC    *category.GetCategoryUseCase
	listUC   *category.ListCategoriesUseCase
	updateUC *category.UpdateCategoryUseCase
	deleteUC *category.DeleteCategoryUseCase
}

func NewCategoryHandler(
	create *category.CreateCategoryUseCase,
	get *category.GetCategoryUseCase,
	list *category.ListCategoriesUseCase,
	update *category.UpdateCategoryUseCase,
	del *category.DeleteCategoryUseCase,
) *CategoryHandler {
	return &CategoryHandler{createUC: create, getUC: get, listUC: list, updateUC: update, deleteUC: del}
}

func (h *CategoryHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/categories", h.List)
	r.Get("/api/categories/{id}", h.Get)
	r.Post("/api/categories", h.Create)
	r.Put("/api/categories/{id}", h.Update)
	r.Delete("/api/categories/{id}", h.Delete)
}

// @Summary List categories
// @Description List all categories
// @Success 200 {array} dto.CategoryResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/categories [get]
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.listUC.Execute(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, categoryResponses(categories))
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

	category, err := h.getUC.Execute(r.Context(), id)
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
	if err := h.createUC.Execute(r.Context(), category); err != nil {
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
	if err := validation.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	category := &entity.Category{
		ID:   id,
		Name: req.Name,
	}
	if err := h.updateUC.Execute(r.Context(), category); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, categoryResponse(category))
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
	if err := h.deleteUC.Execute(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
