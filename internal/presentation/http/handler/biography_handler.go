package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/biography"
)

type BiographyHandler struct {
	getUC    *biography.GetBiographyUseCase
	createUC *biography.CreateBiographyUseCase
	updateUC *biography.UpdateBiographyUseCase
}

func NewBiographyHandler(g *biography.GetBiographyUseCase, c *biography.CreateBiographyUseCase, u *biography.UpdateBiographyUseCase) *BiographyHandler {
	return &BiographyHandler{getUC: g, createUC: c, updateUC: u}
}

func (h *BiographyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/biography", h.GetLatest)
	r.Post("/api/biography", h.Create)
	r.Put("/api/biography", h.Update)
}

// @Summary Get biography
// @Description Get the author biography
// @Success 200 {object} dto.BiographyResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/biography [get]
func (h *BiographyHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	b, err := h.getUC.Execute(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, biographyResponse(b))
}


// @Summary Create biography
// @Description Create the initial author biography (only used once, when no biography exists yet)
// @Tags biography
// @Accept json
// @Produce json
// @Param request body dto.BiographyCreateRequest true "Biography content"
// @Success 201 {object} dto.BiographyResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Security BearerAuth
// @Router /api/biography [post]
func (h *BiographyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.BiographyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := validation.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	b := &entity.Biography{
		ID:   uuid.New(),
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}

	if err := h.createUC.Execute(r.Context(), b); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, biographyResponse(b))
}

// @Summary Update biography
// @Description Update the author biography
// @Accept json
// @Param request body dto.BiographyUpdateRequest true "Biography content"
// @Success 200 {object} dto.BiographyResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/biography [put]
func (h *BiographyHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.BiographyUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	b, err := h.getUC.Execute(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	b.Content = req.Content
	if err := h.updateUC.Execute(r.Context(), b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, biographyResponse(b))
}
