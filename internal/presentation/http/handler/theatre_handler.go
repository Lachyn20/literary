package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/pagination"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/theatre"
)

type TheatreHandler struct {
	svc *theatre.TheatreService
}

func NewTheatreHandler(svc *theatre.TheatreService) *TheatreHandler {
	return &TheatreHandler{svc: svc}
}

func (h *TheatreHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/theatre-productions", h.List)
	r.Get("/api/theatre-productions/{id}", h.Get)
	r.Post("/api/theatre-productions", h.Create)
	r.Put("/api/theatre-productions/{id}", h.Update)
	r.Delete("/api/theatre-productions/{id}", h.Delete)
}

// @Summary List theatre productions
// @Description List all theatre productions with pagination support
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param page query int false "Page number (1-indexed, alternative to offset)"
// @Success 200 {array} dto.TheatreResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/theatre-productions [get]
func (h *TheatreHandler) List(w http.ResponseWriter, r *http.Request) {
	paginationParams := pagination.Parse(r)
	items, total, err := h.svc.List(r.Context(), paginationParams.Limit, paginationParams.Offset)
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	
	paginationInfo := pagination.NewInfo(paginationParams.Limit, paginationParams.Offset, total)
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "ok",
		"data":         theatreResponses(items),
		"pagination":   paginationInfo,
	})
}

// @Summary Get theatre production
// @Description Get a single theatre production by id
// @Param id path string true "Theatre production ID"
// @Success 200 {object} dto.TheatreResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/theatre-productions/{id} [get]
func (h *TheatreHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	t, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	WriteJSON(w, http.StatusOK, theatreResponse(t))
}

// @Summary Create theatre production
// @Description Create a new theatre production
// @Accept json
// @Param request body dto.TheatreCreateRequest true "Theatre production data"
// @Success 201 {object} dto.TheatreResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/theatre-productions [post]
func (h *TheatreHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.TheatreCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	t := entity.TheatreProduction{ID: uuid.New(), PlayTitle: req.PlayTitle, TheatreName: req.TheatreName, CreatedAt: time.Now()}
	if req.PremiereDate != nil {
		parsedTime, _ := time.Parse("2006-01-02", *req.PremiereDate)
		t.PremiereDate = parsedTime
	}
	t.Notes = req.Notes
	if err := h.svc.Create(r.Context(), &t); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusCreated, theatreResponse(&t))
}

// @Summary Update theatre production
// @Description Update an existing theatre production
// @Accept json
// @Param id path string true "Theatre production ID"
// @Param request body dto.TheatreCreateRequest true "Theatre production data"
// @Success 200 {object} dto.TheatreResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/theatre-productions/{id} [put]
func (h *TheatreHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	var req dto.TheatreCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	t := entity.TheatreProduction{ID: id, PlayTitle: req.PlayTitle, TheatreName: req.TheatreName}
	if req.PremiereDate != nil {
		parsedTime, _ := time.Parse("2006-01-02", *req.PremiereDate)
		t.PremiereDate = parsedTime
	}
	t.Notes = req.Notes
	if err := h.svc.Update(r.Context(), &t); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, theatreResponse(&t))
}

// @Summary Delete theatre production
// @Description Delete an existing theatre production by id
// @Param id path string true "Theatre production ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/theatre-productions/{id} [delete]
func (h *TheatreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	if err := h.svc.Delete(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
