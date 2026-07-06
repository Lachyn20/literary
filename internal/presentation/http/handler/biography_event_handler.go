package handler

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/hemra-siirow/literary/internal/domain/entity"
    "github.com/hemra-siirow/literary/internal/presentation/http/dto"
    "github.com/hemra-siirow/literary/internal/presentation/http/validation"
    "github.com/hemra-siirow/literary/internal/usecase/biography"
)

type BiographyEventHandler struct {
    svc *biography.BiographyEventService
    bioSvc *biography.BiographyService
}

func NewBiographyEventHandler(svc *biography.BiographyEventService, bioSvc *biography.BiographyService) *BiographyEventHandler {
    return &BiographyEventHandler{svc: svc, bioSvc: bioSvc}
}

func (h *BiographyEventHandler) RegisterRoutes(r chi.Router) {
    r.Get("/api/biography/events", h.List)
    r.Post("/api/biography/events", h.Create)
    r.Put("/api/biography/events/{id}", h.Update)
    r.Delete("/api/biography/events/{id}", h.Delete)
}

// @Summary List biography events
// @Description Get biography timeline events in the specified language. Falls back to Turkmen if the requested language field is empty.
// @Tags biography
// @Produce json
// @Param lang query string false "Language code: tk (default), ru, en"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/biography/events [get]
func (h *BiographyEventHandler) List(w http.ResponseWriter, r *http.Request) {
    lang := r.URL.Query().Get("lang")
    if lang != "tk" && lang != "ru" && lang != "en" {
        lang = "tk"
    }

    bio, err := h.bioSvc.GetLatest(r.Context())
    if err != nil {
        WriteError(w, http.StatusNotFound, "biography not found")
        return
    }

    events, err := h.svc.List(r.Context(), bio.ID)
    if err != nil {
        WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

    if events == nil {
        events = []*entity.BiographyEvent{}
    }

    localized := make([]map[string]interface{}, 0, len(events))
    for _, e := range events {
        localized = append(localized, localizeEvent(e, lang))
    }

    WriteJSON(w, http.StatusOK, map[string]interface{}{
        "status": "ok",
        "data":   localized,
        "total":  len(localized),
    })
}

// @Summary Create biography event
// @Accept json
// @Param request body dto.BiographyEventCreateRequest true "Biography event"
// @Success 201 {object} dto.BiographyEventResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Security BearerAuth
// @Router /api/biography/events [post]
func (h *BiographyEventHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.BiographyEventCreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
    if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }

    event := &entity.BiographyEvent{
        Year: req.Year,
        TitleTk: ptr(req.TitleTk),
        TitleRu: ptr(req.TitleRu),
        TitleEn: ptr(req.TitleEn),
        DescriptionTk: ptr(req.DescriptionTk),
        DescriptionRu: ptr(req.DescriptionRu),
        DescriptionEn: ptr(req.DescriptionEn),
        SortOrder: req.SortOrder,
        CreatedAt: time.Now(),
    }

    if err := h.svc.Create(r.Context(), event); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
    WriteJSON(w, http.StatusCreated, biographyEventToDTOPtr(event))
}

// @Summary Update biography event
// @Accept json
// @Param id path string true "Event ID"
// @Param request body dto.BiographyEventUpdateRequest true "Biography event partial"
// @Success 200 {object} dto.BiographyEventResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Security BearerAuth
// @Router /api/biography/events/{id} [put]
func (h *BiographyEventHandler) Update(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }

    var req dto.BiographyEventUpdateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
    if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }

    ev, err := h.svc.GetByID(r.Context(), id)
    if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }

    if req.Year != nil { ev.Year = *req.Year }
    if req.TitleTk != "" { ev.TitleTk = ptr(req.TitleTk) }
    if req.TitleRu != "" { ev.TitleRu = ptr(req.TitleRu) }
    if req.TitleEn != "" { ev.TitleEn = ptr(req.TitleEn) }
    if req.DescriptionTk != "" { ev.DescriptionTk = ptr(req.DescriptionTk) }
    if req.DescriptionRu != "" { ev.DescriptionRu = ptr(req.DescriptionRu) }
    if req.DescriptionEn != "" { ev.DescriptionEn = ptr(req.DescriptionEn) }
    if req.SortOrder != nil { ev.SortOrder = *req.SortOrder }

    if err := h.svc.Update(r.Context(), ev); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
    WriteJSON(w, http.StatusOK, biographyEventToDTO(ev))
}

// @Summary Delete biography event
// @Param id path string true "Event ID"
// @Success 204
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Security BearerAuth
// @Router /api/biography/events/{id} [delete]
func (h *BiographyEventHandler) Delete(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
    if err := h.svc.Delete(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
    w.WriteHeader(http.StatusNoContent)
}

func ptr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}

// coalesceStr returns primary if non-nil and non-empty, otherwise fallback
func coalesceStr(primary *string, fallback *string) *string {
    if primary != nil && *primary != "" {
        return primary
    }
    return fallback
}

// localizeEvent returns event with title and description in requested lang.
// Falls back to Turkmen if requested lang field is empty.
func localizeEvent(e *entity.BiographyEvent, lang string) map[string]interface{} {
    item := map[string]interface{}{
        "id":           e.ID,
        "biography_id": e.BiographyID,
        "year":         e.Year,
        "sort_order":   e.SortOrder,
        "created_at":   e.CreatedAt,
    }
    switch lang {
    case "ru":
        item["title"] = coalesceStr(e.TitleRu, e.TitleTk)
        item["description"] = coalesceStr(e.DescriptionRu, e.DescriptionTk)
    case "en":
        item["title"] = coalesceStr(e.TitleEn, e.TitleTk)
        item["description"] = coalesceStr(e.DescriptionEn, e.DescriptionTk)
    default:
        item["title"] = e.TitleTk
        item["description"] = e.DescriptionTk
    }
    return item
}

func biographyEventToDTO(e *entity.BiographyEvent) dto.BiographyEventResponse {
    return dto.BiographyEventResponse{
        ID: e.ID.String(),
        BiographyID: e.BiographyID.String(),
        Year: e.Year,
        TitleTk: e.TitleTk,
        TitleRu: e.TitleRu,
        TitleEn: e.TitleEn,
        DescriptionTk: e.DescriptionTk,
        DescriptionRu: e.DescriptionRu,
        DescriptionEn: e.DescriptionEn,
        SortOrder: e.SortOrder,
        CreatedAt: e.CreatedAt,
    }
}

func biographyEventToDTOPtr(e *entity.BiographyEvent) *dto.BiographyEventResponse {
    d := biographyEventToDTO(e)
    return &d
}
