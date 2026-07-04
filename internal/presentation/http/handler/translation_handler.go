package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
	"github.com/hemra-siirow/literary/internal/presentation/http/validation"
	"github.com/hemra-siirow/literary/internal/usecase/translation"
)

type TranslationHandler struct {
	svc *translation.TranslationService
}

func NewTranslationHandler(svc *translation.TranslationService) *TranslationHandler {
	return &TranslationHandler{svc: svc}
}

func (h *TranslationHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/translations/by-author", h.ListByAuthor)
	r.Get("/api/translations/by-author/{id}", h.GetByAuthor)
	r.Post("/api/translations/by-author", h.CreateByAuthor)
	r.Delete("/api/translations/by-author/{id}", h.DeleteByAuthor)

	r.Get("/api/translations/into-language", h.ListInto)
	r.Get("/api/translations/into-language/{id}", h.GetInto)
	r.Post("/api/translations/into-language", h.CreateInto)
	r.Delete("/api/translations/into-language/{id}", h.DeleteInto)
}

// @Summary List translations by author
// @Description List all translations by author
// @Success 200 {object} dto.TranslatedByAuthorListResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/translations/by-author [get]
func (h *TranslationHandler) ListByAuthor(w http.ResponseWriter, r *http.Request) {
	items, total, err := h.svc.ListByAuthor(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []*entity.TranslatedByAuthor{}
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "data": translatedByAuthorResponses(items), "total": total})
}

// @Summary Get translation by author
// @Description Get a single translation by author by id
// @Param id path string true "Translation ID"
// @Success 200 {object} dto.TranslatedByAuthorResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/translations/by-author/{id} [get]
func (h *TranslationHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	i, err := h.svc.GetByAuthor(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, translatedByAuthorResponse(i))
}

// @Summary Create translation by author
// @Description Create a new translation by author
// @Accept json
// @Param request body dto.TranslatedByAuthorCreateRequest true "Translation by author data"
// @Success 201 {object} dto.TranslatedByAuthorResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/translations/by-author [post]
func (h *TranslationHandler) CreateByAuthor(w http.ResponseWriter, r *http.Request) {
	var req dto.TranslatedByAuthorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := validation.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	p := entity.TranslatedByAuthor{ID: uuid.New(), OriginalAuthorName: req.OriginalAuthorName, OriginalLanguage: req.OriginalLanguage, WorkTitle: req.WorkTitle, Notes: req.Notes}
	if err := h.svc.CreateByAuthor(r.Context(), &p); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, translatedByAuthorResponse(&p))
}

// @Summary Delete translation by author
// @Description Delete an existing translation by author by id
// @Param id path string true "Translation ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/translations/by-author/{id} [delete]
func (h *TranslationHandler) DeleteByAuthor(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteByAuthor(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}

// @Summary List translations into language
// @Description List all translations into language
// @Success 200 {object} dto.TranslatedIntoLanguageListResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/translations/into-language [get]
func (h *TranslationHandler) ListInto(w http.ResponseWriter, r *http.Request) {
	items, total, err := h.svc.ListIntoLanguage(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []*entity.TranslatedIntoLanguage{}
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"status": "ok", "data": translatedIntoLanguageResponses(items), "total": total})
}

// @Summary Get translation into language
// @Description Get a single translation into language by id
// @Param id path string true "Translation ID"
// @Success 200 {object} dto.TranslatedIntoLanguageResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/translations/into-language/{id} [get]
func (h *TranslationHandler) GetInto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	i, err := h.svc.GetIntoLanguage(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, translatedIntoLanguageResponse(i))
}

// @Summary Create translation into language
// @Description Create a new translation into language
// @Accept json
// @Param request body dto.TranslatedIntoLanguageCreateRequest true "Translation into language data"
// @Success 201 {object} dto.TranslatedIntoLanguageResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/translations/into-language [post]
func (h *TranslationHandler) CreateInto(w http.ResponseWriter, r *http.Request) {
	var req dto.TranslatedIntoLanguageCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := validation.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	p := entity.TranslatedIntoLanguage{ID: uuid.New(), LanguageName: req.LanguageName, TranslatorName: req.TranslatorName, WorkTitle: req.WorkTitle, Notes: req.Notes}
	if err := h.svc.CreateIntoLanguage(r.Context(), &p); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, translatedIntoLanguageResponse(&p))
}

// @Summary Delete translation into language
// @Description Delete an existing translation into language by id
// @Param id path string true "Translation ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/translations/into-language/{id} [delete]
func (h *TranslationHandler) DeleteInto(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteIntoLanguage(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
