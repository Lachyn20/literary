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
	"github.com/hemra-siirow/literary/internal/usecase/film"
)

type FilmHandler struct {
	svc   *film.FilmService
	store repository.FileStorage
}

func NewFilmHandler(svc *film.FilmService, store repository.FileStorage) *FilmHandler {
	return &FilmHandler{svc: svc, store: store}
}

func (h *FilmHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/films", h.List)
	r.Get("/api/films/{id}", h.Get)
	r.Post("/api/films", h.Create)
	r.Put("/api/films/{id}", h.Update)
	r.Delete("/api/films/{id}", h.Delete)
}

// @Summary List films
// @Description List all films with pagination support
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param page query int false "Page number (1-indexed, alternative to offset)"
// @Success 200 {array} dto.FilmResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/films [get]
func (h *FilmHandler) List(w http.ResponseWriter, r *http.Request) {
	paginationParams := pagination.Parse(r)
	items, total, err := h.svc.List(r.Context(), paginationParams.Limit, paginationParams.Offset)
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	
	paginationInfo := pagination.NewInfo(paginationParams.Limit, paginationParams.Offset, total)
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "ok",
		"data":         filmResponses(items),
		"pagination":   paginationInfo,
	})
}

// @Summary Get film
// @Description Get a single film by id
// @Param id path string true "Film ID"
// @Success 200 {object} dto.FilmResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/films/{id} [get]
func (h *FilmHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	f, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	WriteJSON(w, http.StatusOK, filmResponse(f))
}

// @Summary Create film
// @Description Create a new film. Supports JSON or multipart/form-data with optional video upload.
// @Accept json
// @Accept multipart/form-data
// @Param title formData string true "Title"
// @Param film_type formData string true "Film type (film or animation)"
// @Param release_year formData int false "Release year"
// @Param director formData string false "Director"
// @Param based_on_scenario formData bool false "Based on scenario"
// @Param file formData file false "Video file — accepted: .mp4, .mov, .mkv"
// @Success 201 {object} dto.FilmResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/films [post]
func (h *FilmHandler) Create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(128 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var f entity.Film
		f.Title = r.FormValue("title")
		f.FilmType = entity.FilmType(r.FormValue("film_type"))
		if v := r.FormValue("release_year"); v != "" { if i, err := strconv.Atoi(v); err == nil { f.ReleaseYear = &i } }
		if v := r.FormValue("director"); v != "" { f.Director = &v }
		if v := r.FormValue("based_on_scenario"); v != "" { if b, err := strconv.ParseBool(v); err == nil { f.BasedOnScenario = b } }
		if f.ID == uuid.Nil { f.ID = uuid.New() }
		f.CreatedAt = time.Now()
		var savedPath string
		var savedFile bool
		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			if !isAllowedExtension(fh.Filename, []string{".mp4", ".mov", ".mkv"}) {
				WriteError(w, http.StatusBadRequest, "only .mp4, .mov, .mkv video files are allowed")
				return
			}
			savedPath, err = h.store.Save(file, fh.Filename, "video")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			savedFile = true
			f.VideoPath = &savedPath
		}
		if err := h.svc.Create(r.Context(), &f); err != nil {
			if savedFile { _ = h.store.Remove(savedPath) }
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		WriteJSON(w, http.StatusCreated, filmResponse(&f))
		return
	}

	var req dto.FilmCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	f := entity.Film{ID: uuid.New(), Title: req.Title, FilmType: entity.FilmType(req.FilmType), BasedOnScenario: req.BasedOnScenario, Director: req.Director, ReleaseYear: req.ReleaseYear, CreatedAt: time.Now()}
	if err := h.svc.Create(r.Context(), &f); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusCreated, filmResponse(&f))
}

// @Summary Update film
// @Description Update an existing film. Supports JSON or multipart/form-data with optional video upload.
// @Accept json
// @Accept multipart/form-data
// @Param id path string true "Film ID"
// @Param title formData string false "Title"
// @Param film_type formData string false "Film type (film or animation)"
// @Param release_year formData int false "Release year"
// @Param director formData string false "Director"
// @Param based_on_scenario formData bool false "Based on scenario"
// @Param file formData file false "Video file — accepted: .mp4, .mov, .mkv"
// @Success 200 {object} dto.FilmResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/films/{id} [put]
func (h *FilmHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(128 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		// Fetch existing film to preserve film_type if not provided
		existing, err := h.svc.GetByID(r.Context(), id)
		if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
		var f entity.Film
		f.ID = id
		f.Title = r.FormValue("title")
		if f.Title == "" { f.Title = existing.Title }
		f.FilmType = existing.FilmType
		if v := r.FormValue("film_type"); v != "" { f.FilmType = entity.FilmType(v) }
		if v := r.FormValue("release_year"); v != "" { if i, err := strconv.Atoi(v); err == nil { f.ReleaseYear = &i } }
		if v := r.FormValue("director"); v != "" { f.Director = &v }
		if v := r.FormValue("based_on_scenario"); v != "" { if b, err := strconv.ParseBool(v); err == nil { f.BasedOnScenario = b } }
		f.VideoPath = existing.VideoPath
		var savedPath string
		var savedFile bool
		var oldPath string
		if existing.VideoPath != nil { oldPath = *existing.VideoPath }
		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			if !isAllowedExtension(fh.Filename, []string{".mp4", ".mov", ".mkv"}) {
				WriteError(w, http.StatusBadRequest, "only .mp4, .mov, .mkv video files are allowed")
				return
			}
			savedPath, err = h.store.Save(file, fh.Filename, "video")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			savedFile = true
			f.VideoPath = &savedPath
		}
		if err := h.svc.Update(r.Context(), &f); err != nil {
			if savedFile { _ = h.store.Remove(savedPath) }
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if savedFile && oldPath != "" && oldPath != savedPath {
			_ = h.store.Remove(oldPath)
		}
		WriteJSON(w, http.StatusOK, filmResponse(&f))
		return
	}

	var req dto.FilmCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }

	old, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	updated := *old
	if req.Title != "" { updated.Title = req.Title }
	if req.FilmType != "" {
		// validate film_type
		if req.FilmType != "film" && req.FilmType != "animation" {
			WriteError(w, http.StatusBadRequest, "film_type must be 'film' or 'animation'")
			return
		}
		updated.FilmType = entity.FilmType(req.FilmType)
	}
	if req.ReleaseYear != nil { updated.ReleaseYear = req.ReleaseYear }
	if req.Director != nil { updated.Director = req.Director }
	updated.BasedOnScenario = req.BasedOnScenario

	if err := h.svc.Update(r.Context(), &updated); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, filmResponse(&updated))
}

// @Summary Delete film
// @Description Delete an existing film by id
// @Param id path string true "Film ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/films/{id} [delete]
func (h *FilmHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	existing, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	if err := h.svc.Delete(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	if existing.VideoPath != nil { _ = h.store.Remove(*existing.VideoPath) }
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
