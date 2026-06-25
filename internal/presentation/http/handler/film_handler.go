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
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/film"
)

type FilmHandler struct {
	createUC *film.CreateFilmUseCase
	getUC    *film.GetFilmUseCase
	listUC   *film.ListFilmsUseCase
	updateUC *film.UpdateFilmUseCase
	deleteUC *film.DeleteFilmUseCase
	store    repository.FileStorage
}

func NewFilmHandler(c *film.CreateFilmUseCase, g *film.GetFilmUseCase, l *film.ListFilmsUseCase, u *film.UpdateFilmUseCase, d *film.DeleteFilmUseCase, s repository.FileStorage) *FilmHandler {
	return &FilmHandler{createUC: c, getUC: g, listUC: l, updateUC: u, deleteUC: d, store: s}
}

func (h *FilmHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/films", h.List)
	r.Get("/api/films/{id}", h.Get)
	r.Post("/api/films", h.Create)
	r.Put("/api/films/{id}", h.Update)
	r.Delete("/api/films/{id}", h.Delete)
}

// @Summary List films
// @Description List all films
// @Success 200 {array} dto.FilmResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/films [get]
func (h *FilmHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.listUC.Execute(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, filmResponses(items))
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
	f, err := h.getUC.Execute(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	WriteJSON(w, http.StatusOK, filmResponse(f))
}

// @Summary Create film
// @Description Create a new film. Supports JSON or multipart/form-data with optional video upload.
// @Accept json
// @Accept multipart/form-data
// @Param title formData string true "Title"
// @Param release_year formData int false "Release year"
// @Param director formData string false "Director"
// @Param based_on_scenario formData bool false "Based on scenario"
// @Param file formData file false "Video file"
// @Success 201 {object} dto.FilmResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/films [post]
func (h *FilmHandler) Create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(128 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var f entity.Film
		f.Title = r.FormValue("title")
		if v := r.FormValue("release_year"); v != "" { if i, err := strconv.Atoi(v); err == nil { f.ReleaseYear = &i } }
		if v := r.FormValue("director"); v != "" { f.Director = &v }
		if v := r.FormValue("based_on_scenario"); v != "" { if b, err := strconv.ParseBool(v); err == nil { f.BasedOnScenario = b } }
		if f.ID == uuid.Nil { f.ID = uuid.New() }
		f.CreatedAt = time.Now()
		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			path, err := h.store.Save(file, fh.Filename, "video")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			f.VideoPath = &path
		}
		if err := h.createUC.Execute(r.Context(), &f); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
		WriteJSON(w, http.StatusCreated, filmResponse(&f))
		return
	}

	var req dto.FilmCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	f := entity.Film{ID: uuid.New(), Title: req.Title, FilmType: entity.FilmType(req.FilmType), BasedOnScenario: req.BasedOnScenario, Director: req.Director, ReleaseYear: req.ReleaseYear, CreatedAt: time.Now()}
	if err := h.createUC.Execute(r.Context(), &f); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusCreated, filmResponse(&f))
}

// @Summary Update film
// @Description Update an existing film. Supports JSON or multipart/form-data with optional video upload.
// @Accept json
// @Accept multipart/form-data
// @Param id path string true "Film ID"
// @Param title formData string false "Title"
// @Param release_year formData int false "Release year"
// @Param director formData string false "Director"
// @Param based_on_scenario formData bool false "Based on scenario"
// @Param file formData file false "Video file"
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
		var f entity.Film
		f.ID = id
		f.Title = r.FormValue("title")
		if v := r.FormValue("release_year"); v != "" { if i, err := strconv.Atoi(v); err == nil { f.ReleaseYear = &i } }
		if v := r.FormValue("director"); v != "" { f.Director = &v }
		if v := r.FormValue("based_on_scenario"); v != "" { if b, err := strconv.ParseBool(v); err == nil { f.BasedOnScenario = b } }
		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			path, err := h.store.Save(file, fh.Filename, "video")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			f.VideoPath = &path
		}
		if err := h.updateUC.Execute(r.Context(), &f); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
		WriteJSON(w, http.StatusOK, filmResponse(&f))
		return
	}

	var req dto.FilmCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	f := entity.Film{ID: id, Title: req.Title, FilmType: entity.FilmType(req.FilmType), BasedOnScenario: req.BasedOnScenario, Director: req.Director, ReleaseYear: req.ReleaseYear}
	if err := h.updateUC.Execute(r.Context(), &f); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, filmResponse(&f))
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
	if err := h.deleteUC.Execute(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
