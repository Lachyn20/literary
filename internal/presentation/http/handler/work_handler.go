package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/work"
)

type WorkHandler struct {
	createUC *work.CreateWorkUseCase
	getUC    *work.GetWorkUseCase
	listUC   *work.ListWorksUseCase
	updateUC *work.UpdateWorkUseCase
	deleteUC *work.DeleteWorkUseCase
}

func NewWorkHandler(create *work.CreateWorkUseCase, get *work.GetWorkUseCase, list *work.ListWorksUseCase, update *work.UpdateWorkUseCase, del *work.DeleteWorkUseCase) *WorkHandler {
	return &WorkHandler{createUC: create, getUC: get, listUC: list, updateUC: update, deleteUC: del}
}

func (h *WorkHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/works", h.List)
	r.Get("/api/works/{id}", h.Get)
	r.Post("/api/works", h.Create)
	r.Put("/api/works/{id}", h.Update)
	r.Delete("/api/works/{id}", h.Delete)
}

// @Summary List works
// @Description List works with optional search, filters and pagination
// @Param search query string false "search keywords"
// @Param category query string false "category id"
// @Param audience_type query string false "audience type"
// @Param year query int false "publish year"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} handler.JSONResponse
// @Router /api/works [get]
func (h *WorkHandler) List(w http.ResponseWriter, r *http.Request) {
	// parse query params
	q := r.URL.Query()
	var filter repository.WorkFilter
	if s := q.Get("search"); s != "" {
		filter.Search = &s
	}
	if cat := q.Get("category"); cat != "" {
		if id, err := uuid.Parse(cat); err == nil {
			filter.CategoryID = &id
		}
	}
	if at := q.Get("audience_type"); at != "" {
		a := entity.AudienceType(at)
		filter.AudienceType = &a
	}
	if y := q.Get("year"); y != "" {
		if yi, err := strconv.Atoi(y); err == nil {
			filter.PublishYear = &yi
		}
	}
	page := 1
	limit := 20
	if p := q.Get("page"); p != "" {
		if pi, err := strconv.Atoi(p); err == nil && pi > 0 { page = pi }
	}
	if l := q.Get("limit"); l != "" {
		if li, err := strconv.Atoi(l); err == nil && li > 0 { limit = li }
	}
	filter.Page = page
	filter.Limit = limit

	works, total, err := h.listUC.Execute(r.Context(), filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := map[string]interface{}{"data": workResponses(works), "total": total, "page": page, "limit": limit}
	WriteJSON(w, http.StatusOK, resp)
}

// @Summary Get work
// @Description Get single work by id
// @Param id path string true "work id"
// @Success 200 {object} dto.WorkResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/works/{id} [get]
func (h *WorkHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	wk, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, workResponse(wk))
}

// @Summary Create work
// @Description Create a new work
// @Accept json
// @Param payload body dto.WorkCreateRequest true "work payload"
// @Success 201 {object} dto.WorkResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/works [post]
func (h *WorkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.WorkCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := validation.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	catID, _ := uuid.Parse(req.CategoryID)
	now := time.Now()
	work := &entity.Work{
		ID: uuid.New(),
		Title: req.Title,
		CategoryID: catID,
		Content: &req.Content,
		Description: &req.Description,
		AudienceType: entity.AudienceType(req.AudienceType),
		PublishYear: req.PublishYear,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := h.createUC.Execute(r.Context(), work); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, workResponse(work))
}

// @Summary Update work
// @Description Update an existing work by id
// @Param id path string true "work id"
// @Accept json
// @Param payload body dto.WorkCreateRequest true "work payload"
// @Success 200 {object} dto.WorkResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/works/{id} [put]
func (h *WorkHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req dto.WorkCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	catID, _ := uuid.Parse(req.CategoryID)
	now := time.Now()
	work := &entity.Work{
		ID: id,
		Title: req.Title,
		CategoryID: catID,
		Content: &req.Content,
		Description: &req.Description,
		AudienceType: entity.AudienceType(req.AudienceType),
		PublishYear: req.PublishYear,
		UpdatedAt: now,
	}
	if err := h.updateUC.Execute(r.Context(), work); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, workResponse(work))
}

// @Summary Delete work
// @Description Delete an existing work by id
// @Param id path string true "work id"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/works/{id} [delete]
func (h *WorkHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
