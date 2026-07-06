package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/domain/repository"
	"github.com/hemra-siirow/literary/internal/usecase/biography"
)

type BiographyHandler struct {
	svc   *biography.BiographyService
	store repository.FileStorage
}

func NewBiographyHandler(svc *biography.BiographyService, store repository.FileStorage) *BiographyHandler {
	return &BiographyHandler{svc: svc, store: store}
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
	b, err := h.svc.GetLatest(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, biographyResponse(b))
}


// @Summary Create biography
// @Description Create the initial author biography (multipart, optional photo)
// @Tags biography
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file false "Author photo (jpg, jpeg, png, webp)"
// @Success 201 {object} dto.BiographyResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Security BearerAuth
// @Router /api/biography [post]
func (h *BiographyHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Only accept multipart/form-data with optional photo
	ct := r.Header.Get("Content-Type")
	if ct == "" || !strings.HasPrefix(ct, "multipart/") {
		WriteError(w, http.StatusBadRequest, "only multipart/form-data is accepted")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	b := &entity.Biography{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
	}

	photo, ph, err := r.FormFile("photo")
	if err == nil {
		defer photo.Close()
		if !isAllowedExtension(ph.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}) {
			WriteError(w, http.StatusBadRequest, "photo must be .jpg, .jpeg, .png, or .webp")
			return
		}
		savedPath, err := h.store.Save(photo, ph.Filename, "biography")
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "failed to save photo")
			return
		}
		b.PhotoPath = &savedPath
	}

	if err := h.svc.Create(r.Context(), b); err != nil {
		if b.PhotoPath != nil { _ = h.store.Remove(*b.PhotoPath) }
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
	// Only support multipart photo update
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid multipart form")
			return
		}
		b, err := h.svc.GetLatest(r.Context())
		if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }

		photo, ph, err := r.FormFile("photo")
		if err != nil {
			WriteError(w, http.StatusBadRequest, "photo is required")
			return
		}
		defer photo.Close()
		if !isAllowedExtension(ph.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}) {
			WriteError(w, http.StatusBadRequest, "photo must be .jpg, .jpeg, .png, or .webp")
			return
		}
		photoPath, err := h.store.Save(photo, ph.Filename, "biography")
		if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }

		oldPath := b.PhotoPath
		b.PhotoPath = &photoPath
		b.UpdatedAt = time.Now()
		if err := h.svc.Update(r.Context(), b); err != nil {
			_ = h.store.Remove(photoPath)
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if oldPath != nil && *oldPath != photoPath {
			_ = h.store.Remove(*oldPath)
		}
		WriteJSON(w, http.StatusOK, biographyResponse(b))
		return
	}
	WriteError(w, http.StatusBadRequest, "only multipart photo update supported")
}
