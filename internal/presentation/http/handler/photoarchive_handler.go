package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/photoarchive"
)

type PhotoArchiveHandler struct {
	svc   *photoarchive.PhotoArchiveService
	store repository.FileStorage
}

func NewPhotoArchiveHandler(svc *photoarchive.PhotoArchiveService, store repository.FileStorage) *PhotoArchiveHandler {
	return &PhotoArchiveHandler{svc: svc, store: store}
}

func (h *PhotoArchiveHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/photo-archive", h.List)
	r.Get("/api/photo-archive/{id}", h.Get)
	r.Post("/api/photo-archive", h.Create)
	r.Put("/api/photo-archive/{id}", h.Update)
	r.Delete("/api/photo-archive/{id}", h.Delete)
}

// @Summary List photo archive
// @Description List all photo archive entries
// @Success 200 {array} dto.PhotoArchiveResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/photo-archive [get]
func (h *PhotoArchiveHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, photoArchiveResponses(items))
}

// @Summary Get photo archive
// @Description Get a single photo archive entry by id
// @Param id path string true "Photo archive ID"
// @Success 200 {object} dto.PhotoArchiveResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/photo-archive/{id} [get]
func (h *PhotoArchiveHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	WriteJSON(w, http.StatusOK, photoArchiveResponse(p))
}

// @Summary Create photo archive entry
// @Description Create a new photo archive entry. Supports JSON or multipart/form-data with optional image upload.
// @Accept json
// @Accept multipart/form-data
// @Param title formData string true "Title"
// @Param description formData string false "Description"
// @Param category formData string false "Category"
// @Param image formData file false "Image file"
// @Success 201 {object} dto.PhotoArchiveResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/photo-archive [post]
func (h *PhotoArchiveHandler) Create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(32 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var p entity.PhotoArchive
		p.Title = r.FormValue("title")
		if v := r.FormValue("description"); v != "" { p.Description = &v }
		if v := r.FormValue("category"); v != "" {
			switch entity.PhotoCategory(v) {
			case entity.PhotoCategoryArchive, entity.PhotoCategoryPersonal:
				p.Category = entity.PhotoCategory(v)
			default:
				WriteError(w, http.StatusBadRequest, "category must be archive or personal")
				return
			}
		} else {
			p.Category = entity.PhotoCategoryArchive
		}
		if p.ID == uuid.Nil { p.ID = uuid.New() }
		p.CreatedAt = time.Now()
		img, ih, err := r.FormFile("image")
		if err == nil {
			defer img.Close()
			lower := strings.ToLower(ih.Filename)
			if !strings.HasSuffix(lower, ".jpg") && !strings.HasSuffix(lower, ".jpeg") && !strings.HasSuffix(lower, ".png") && !strings.HasSuffix(lower, ".gif") && !strings.HasSuffix(lower, ".webp") {
				WriteError(w, http.StatusBadRequest, "unsupported image file type")
				return
			}
			path, err := h.store.Save(img, ih.Filename, "image")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			p.ImagePath = path
		}
		if err := h.svc.Create(r.Context(), &p); err != nil {
			if p.ImagePath != "" { _ = h.store.Remove(p.ImagePath) }
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		WriteJSON(w, http.StatusCreated, photoArchiveResponse(&p))
		return
	}

	var req dto.PhotoArchiveCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	p := entity.PhotoArchive{ID: uuid.New(), Title: req.Title, ImagePath: "", Description: req.Description, CreatedAt: time.Now()}
	if req.TakenDate != nil { p.TakenDate = req.TakenDate }
	if req.Category != nil && *req.Category != "" {
		switch entity.PhotoCategory(*req.Category) {
		case entity.PhotoCategoryArchive, entity.PhotoCategoryPersonal:
			p.Category = entity.PhotoCategory(*req.Category)
		default:
			WriteError(w, http.StatusBadRequest, "category must be archive or personal")
			return
		}
	} else {
		p.Category = entity.PhotoCategoryArchive
	}
	if err := h.svc.Create(r.Context(), &p); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusCreated, photoArchiveResponse(&p))
}

// @Summary Update photo archive entry
// @Description Update an existing photo archive entry. Supports JSON or multipart/form-data with optional image upload.
// @Accept json
// @Accept multipart/form-data
// @Param id path string true "Photo archive ID"
// @Param title formData string false "Title"
// @Param description formData string false "Description"
// @Param image formData file false "Image file"
// @Success 200 {object} dto.PhotoArchiveResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/photo-archive/{id} [put]
func (h *PhotoArchiveHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(32 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		old, err := h.svc.GetByID(r.Context(), id)
		if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
		var p entity.PhotoArchive
		p.ID = id
		p.Title = r.FormValue("title")
		if v := r.FormValue("description"); v != "" { p.Description = &v }
		if v := r.FormValue("category"); v != "" {
			switch entity.PhotoCategory(v) {
			case entity.PhotoCategoryArchive, entity.PhotoCategoryPersonal:
				p.Category = entity.PhotoCategory(v)
			default:
				WriteError(w, http.StatusBadRequest, "category must be archive or personal")
				return
			}
		} else {
			p.Category = entity.PhotoCategoryArchive
		}
		p.ImagePath = old.ImagePath
		var oldPath string
		if old.ImagePath != "" { oldPath = old.ImagePath }
		img, ih, err := r.FormFile("image")
		if err == nil {
			defer img.Close()
			lower := strings.ToLower(ih.Filename)
			if !strings.HasSuffix(lower, ".jpg") && !strings.HasSuffix(lower, ".jpeg") && !strings.HasSuffix(lower, ".png") && !strings.HasSuffix(lower, ".gif") && !strings.HasSuffix(lower, ".webp") {
				WriteError(w, http.StatusBadRequest, "unsupported image file type")
				return
			}
			path, err := h.store.Save(img, ih.Filename, "image")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			p.ImagePath = path
		}
		if err := h.svc.Update(r.Context(), &p); err != nil {
			if p.ImagePath != "" && p.ImagePath != oldPath { _ = h.store.Remove(p.ImagePath) }
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if p.ImagePath != "" && oldPath != "" && p.ImagePath != oldPath { _ = h.store.Remove(oldPath) }
		WriteJSON(w, http.StatusOK, photoArchiveResponse(&p))
		return
	}

	var req dto.PhotoArchiveCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	p := entity.PhotoArchive{ID: id, Title: req.Title, ImagePath: "", Description: req.Description}
	if req.TakenDate != nil { p.TakenDate = req.TakenDate }
	if req.Category != nil && *req.Category != "" {
		switch entity.PhotoCategory(*req.Category) {
		case entity.PhotoCategoryArchive, entity.PhotoCategoryPersonal:
			p.Category = entity.PhotoCategory(*req.Category)
		default:
			WriteError(w, http.StatusBadRequest, "category must be archive or personal")
			return
		}
	} else {
		p.Category = entity.PhotoCategoryArchive
	}
	if err := h.svc.Update(r.Context(), &p); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, photoArchiveResponse(&p))
}

// @Summary Delete photo archive entry
// @Description Delete an existing photo archive entry by id
// @Param id path string true "Photo archive ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/photo-archive/{id} [delete]
func (h *PhotoArchiveHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	old, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	if err := h.svc.Delete(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	if old.ImagePath != "" { _ = h.store.Remove(old.ImagePath) }
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
