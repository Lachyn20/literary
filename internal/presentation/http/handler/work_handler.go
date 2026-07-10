package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/pagination"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/work"
)

type WorkHandler struct {
	svc   *work.WorkService
	store repository.FileStorage
}

func NewWorkHandler(svc *work.WorkService, store repository.FileStorage) *WorkHandler {
	return &WorkHandler{svc: svc, store: store}
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
// @Success 200 {object} dto.WorkListResponse
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
	paginationParams := pagination.Parse(r)

	works, total, err := h.svc.List(r.Context(), filter, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if works == nil {
		works = []*entity.Work{}
	}

	paginationInfo := pagination.NewInfo(paginationParams.Limit, paginationParams.Offset, total)
	resp := map[string]interface{}{"status": "ok", "data": workResponses(works), "pagination": paginationInfo}
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
	wk, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, workResponse(wk))
}

// @Summary Create work
// @Description Create a new work
// @Accept json
// @Accept multipart/form-data
// @Param payload body dto.WorkCreateRequest true "work payload"
// @Param title formData string true "Title"
// @Param category_id formData string true "Category ID"
// @Param audience_type formData string true "Audience type: adult or children"
// @Param description formData string false "Description"
// @Param publish_year formData int false "Publish year"
// @Param file formData file false "Work file (.pdf or .txt)"
// @Success 201 {object} dto.WorkResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/works [post]
func (h *WorkHandler) Create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid multipart")
			return
		}

		title := r.FormValue("title")
		if title == "" {
			WriteError(w, http.StatusBadRequest, "title is required")
			return
		}

		categoryID := r.FormValue("category_id")
		if categoryID == "" {
			WriteError(w, http.StatusBadRequest, "category_id is required")
			return
		}
		catID, err := uuid.Parse(categoryID)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "category_id must be a valid UUID")
			return
		}

		audienceType := r.FormValue("audience_type")
		if audienceType != "adult" && audienceType != "children" {
			WriteError(w, http.StatusBadRequest, "audience_type must be 'adult' or 'children'")
			return
		}

		var description *string
		if d := r.FormValue("description"); d != "" {
			description = &d
		}

		var publishYear *int
		if y := r.FormValue("publish_year"); y != "" {
			yi, err := strconv.Atoi(y)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "publish_year must be a valid integer")
				return
			}
			publishYear = &yi
		}

		now := time.Now()
		work := &entity.Work{
			ID:           uuid.New(),
			Title:        title,
			CategoryID:   catID,
			Description:  description,
			AudienceType: entity.AudienceType(audienceType),
			PublishYear:  publishYear,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		var savedPath string
		savedFile := false
		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			if !isAllowedExtension(fh.Filename, []string{".pdf", ".txt"}) {
				WriteError(w, http.StatusBadRequest, "only .pdf and .txt files are allowed")
				return
			}

			savedPath, err = h.store.Save(file, fh.Filename, "work")
			if err != nil {
				WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			savedFile = true
			work.FilePath = &savedPath
		}

		if err := h.svc.Create(r.Context(), work); err != nil {
			if savedFile {
				_ = h.store.Remove(savedPath)
			}
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		WriteJSON(w, http.StatusCreated, workResponse(work))
		return
	}

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
		ID:           uuid.New(),
		Title:        req.Title,
		CategoryID:   catID,
		AudienceType: entity.AudienceType(req.AudienceType),
		PublishYear:  req.PublishYear,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if req.Description != "" {
		work.Description = &req.Description
	}
	if err := h.svc.Create(r.Context(), work); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, workResponse(work))
}

// @Summary Update work
// @Description Update an existing work by id
// @Param id path string true "work id"
// @Accept json
// @Accept multipart/form-data
// @Param title formData string false "Title"
// @Param category_id formData string false "Category ID"
// @Param audience_type formData string false "Audience type: adult or children"
// @Param description formData string false "Description"
// @Param publish_year formData int false "Publish year"
// @Param file formData file false "Work file (.pdf or .txt)"
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

	oldWork, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid multipart")
			return
		}

		if v := r.FormValue("title"); v != "" {
			oldWork.Title = v
		}
		if v := r.FormValue("category_id"); v != "" {
			catID, err := uuid.Parse(v)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "category_id must be a valid UUID")
				return
			}
			oldWork.CategoryID = catID
		}
		if v := r.FormValue("audience_type"); v != "" {
			if v != "adult" && v != "children" {
				WriteError(w, http.StatusBadRequest, "audience_type must be 'adult' or 'children'")
				return
			}
			oldWork.AudienceType = entity.AudienceType(v)
		}
		if v := r.FormValue("description"); v != "" {
			oldWork.Description = &v
		}
		if v := r.FormValue("publish_year"); v != "" {
			yi, err := strconv.Atoi(v)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "publish_year must be a valid integer")
				return
			}
			oldWork.PublishYear = &yi
		}

		oldWork.UpdatedAt = time.Now()
		prevFilePath := oldWork.FilePath
		var savedPath string
		savedFile := false

		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			if !isAllowedExtension(fh.Filename, []string{".pdf", ".txt"}) {
				WriteError(w, http.StatusBadRequest, "only .pdf and .txt files are allowed")
				return
			}

			savedPath, err = h.store.Save(file, fh.Filename, "work")
			if err != nil {
				WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			savedFile = true
			oldWork.FilePath = &savedPath
		}

		if err := h.svc.Update(r.Context(), oldWork); err != nil {
			if savedFile {
				_ = h.store.Remove(savedPath)
			}
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if savedFile && prevFilePath != nil && (oldWork.FilePath == nil || *prevFilePath != *oldWork.FilePath) {
			_ = h.store.Remove(*prevFilePath)
		}

		WriteJSON(w, http.StatusOK, workResponse(oldWork))
		return
	}

	var req dto.WorkCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if req.Title != "" {
		oldWork.Title = req.Title
	}
	if req.CategoryID != "" {
		catID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "category_id must be a valid UUID")
			return
		}
		oldWork.CategoryID = catID
	}
	if req.AudienceType != "" {
		if req.AudienceType != "adult" && req.AudienceType != "children" {
			WriteError(w, http.StatusBadRequest, "audience_type must be 'adult' or 'children'")
			return
		}
		oldWork.AudienceType = entity.AudienceType(req.AudienceType)
	}
	if req.Description != "" {
		oldWork.Description = &req.Description
	}
	if req.PublishYear != nil {
		oldWork.PublishYear = req.PublishYear
	}
	oldWork.UpdatedAt = time.Now()

	if err := h.svc.Update(r.Context(), oldWork); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, workResponse(oldWork))
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
	if err := h.svc.Delete(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
