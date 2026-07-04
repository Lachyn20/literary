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
// @Description Create the initial author biography (supports JSON or multipart for image upload)
// @Tags biography
// @Accept json, multipart/form-data
// @Produce json
// @Param content formData string false "Biography content (used in multipart)"
// @Param photo formData file false "Biography image file (used in multipart)"
// @Param request body dto.BiographyCreateRequest false "Biography content (used in JSON)"
// @Success 201 {object} dto.BiographyResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Security BearerAuth
// @Router /api/biography [post]
func (h *BiographyHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Ýagdaý: Eger multipart/form-data (suratly) ugradylan bolsa
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid multipart form")
			return
		}

		content := r.FormValue("content")
		if content == "" {
			WriteError(w, http.StatusBadRequest, "content is required")
			return
		}

		b := &entity.Biography{
			ID:        uuid.New(),
			Content:   content,
			UpdatedAt: time.Now(),
		}

		var photoPath string
		var savedPhoto bool

		// Suraty okamak we h.store arkaly ýazdyrmak (BookHandler-däki ýaly)
		photo, ph, err := r.FormFile("photo")
		if err == nil {
			defer photo.Close()
			photoPath, err = h.store.Save(photo, ph.Filename, "biography")
			if err != nil {
				WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			savedPhoto = true
			b.PhotoPath = &photoPath
		}

		// UseCase çagyrmak
		if err := h.svc.Create(r.Context(), b); err != nil {
			if savedPhoto {
				_ = h.store.Remove(photoPath) // Ýalňyşlyk boldy, ýüklenen suraty öçürýäris
			}
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		WriteJSON(w, http.StatusCreated, biographyResponse(b))
		return
	}

	// 2. Ýagdaý: Eger diňe JSON ugradylan bolsa (suratsyz)
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
		ID:        uuid.New(),
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}

	if err := h.svc.Create(r.Context(), b); err != nil {
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
	b, err := h.svc.GetLatest(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	b.Content = req.Content
	if err := h.svc.Update(r.Context(), b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, biographyResponse(b))
}
