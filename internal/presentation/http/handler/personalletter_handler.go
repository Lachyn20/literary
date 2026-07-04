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
	"github.com/hemra-siirow/literary/internal/usecase/personalletter"
)

type PersonalLetterHandler struct {
	svc   *personalletter.PersonalLetterService
	store repository.FileStorage
}

func NewPersonalLetterHandler(svc *personalletter.PersonalLetterService, store repository.FileStorage) *PersonalLetterHandler {
	return &PersonalLetterHandler{svc: svc, store: store}
}

func (h *PersonalLetterHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/personal-letters", h.List)
	r.Get("/api/personal-letters/{id}", h.Get)
	r.Post("/api/personal-letters", h.Create)
	r.Put("/api/personal-letters/{id}", h.Update)
	r.Delete("/api/personal-letters/{id}", h.Delete)
}

// @Summary List personal letters
// @Description List all personal letters
// @Success 200 {array} dto.PersonalLetterResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/personal-letters [get]
func (h *PersonalLetterHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, personalLetterResponses(items))
}

// @Summary Get personal letter
// @Description Get a single personal letter by id
// @Param id path string true "Personal letter ID"
// @Success 200 {object} dto.PersonalLetterResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/personal-letters/{id} [get]
func (h *PersonalLetterHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	WriteJSON(w, http.StatusOK, personalLetterResponse(p))
}

// @Summary Create personal letter
// @Description Create a new personal letter. Supports JSON or multipart/form-data with optional scan upload.
// @Accept json
// @Accept multipart/form-data
// @Param title formData string true "Title"
// @Param content formData string true "Content"
// @Param letter_date formData string false "Letter date in RFC3339"
// @Param scan formData file false "Scan image or PDF"
// @Success 201 {object} dto.PersonalLetterResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/personal-letters [post]
func (h *PersonalLetterHandler) Create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(16 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var p entity.PersonalLetter
		p.Title = r.FormValue("title")
		p.Content = r.FormValue("content")
		if v := r.FormValue("letter_date"); v != "" { if t, err := time.Parse(time.RFC3339, v); err == nil { p.LetterDate = t } }
		if p.ID == uuid.Nil { p.ID = uuid.New() }
		p.CreatedAt = time.Now()
		scan, sh, err := r.FormFile("scan")
		if err == nil {
			defer scan.Close()
			path, err := h.store.Save(scan, sh.Filename, "scan")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			p.ScanImagePath = &path
		}
		if err := h.svc.Create(r.Context(), &p); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
		WriteJSON(w, http.StatusCreated, personalLetterResponse(&p))
		return
	}

	var req dto.PersonalLetterCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	p := entity.PersonalLetter{ID: uuid.New(), Title: req.Title, Content: req.Content, CreatedAt: time.Now()}
	if req.LetterDate != nil { p.LetterDate = *req.LetterDate }
	if err := h.svc.Create(r.Context(), &p); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusCreated, personalLetterResponse(&p))
}

// @Summary Update personal letter
// @Description Update an existing personal letter. Supports JSON or multipart/form-data with optional scan upload.
// @Accept json
// @Accept multipart/form-data
// @Param id path string true "Personal letter ID"
// @Param title formData string false "Title"
// @Param content formData string false "Content"
// @Param letter_date formData string false "Letter date in RFC3339"
// @Param scan formData file false "Scan image or PDF"
// @Success 200 {object} dto.PersonalLetterResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/personal-letters/{id} [put]
func (h *PersonalLetterHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(16 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var p entity.PersonalLetter
		p.ID = id
		p.Title = r.FormValue("title")
		p.Content = r.FormValue("content")
		if v := r.FormValue("letter_date"); v != "" { if t, err := time.Parse(time.RFC3339, v); err == nil { p.LetterDate = t } }
		scan, sh, err := r.FormFile("scan")
		if err == nil {
			defer scan.Close()
			path, err := h.store.Save(scan, sh.Filename, "scan")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			p.ScanImagePath = &path
		}
		if err := h.svc.Update(r.Context(), &p); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
		WriteJSON(w, http.StatusOK, personalLetterResponse(&p))
		return
	}

	var req dto.PersonalLetterCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	p := entity.PersonalLetter{ID: id, Title: req.Title, Content: req.Content}
	if req.LetterDate != nil { p.LetterDate = *req.LetterDate }
	if err := h.svc.Update(r.Context(), &p); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, personalLetterResponse(&p))
}

// @Summary Delete personal letter
// @Description Delete an existing personal letter by id
// @Param id path string true "Personal letter ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/personal-letters/{id} [delete]
func (h *PersonalLetterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	if err := h.svc.Delete(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
