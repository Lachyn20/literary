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
	"github.com/hemra-siirow/literary/internal/usecase/broadcast"
)

type BroadcastHandler struct {
	createUC *broadcast.CreateBroadcastUseCase
	getUC    *broadcast.GetBroadcastUseCase
	listUC   *broadcast.ListBroadcastsUseCase
	updateUC *broadcast.UpdateBroadcastUseCase
	deleteUC *broadcast.DeleteBroadcastUseCase
	store    repository.FileStorage
}

func NewBroadcastHandler(c *broadcast.CreateBroadcastUseCase, g *broadcast.GetBroadcastUseCase, l *broadcast.ListBroadcastsUseCase, u *broadcast.UpdateBroadcastUseCase, d *broadcast.DeleteBroadcastUseCase, s repository.FileStorage) *BroadcastHandler {
	return &BroadcastHandler{createUC: c, getUC: g, listUC: l, updateUC: u, deleteUC: d, store: s}
}

func (h *BroadcastHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/broadcasts", h.List)
	r.Get("/api/broadcasts/{id}", h.Get)
	r.Post("/api/broadcasts", h.Create)
	r.Put("/api/broadcasts/{id}", h.Update)
	r.Delete("/api/broadcasts/{id}", h.Delete)
}

// @Summary List broadcasts
// @Description List all broadcasts
// @Success 200 {array} dto.BroadcastResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/broadcasts [get]
func (h *BroadcastHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.listUC.Execute(r.Context())
	if err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, broadcastResponses(items))
}

// @Summary Get broadcast
// @Description Get a single broadcast by id
// @Param id path string true "Broadcast ID"
// @Success 200 {object} dto.BroadcastResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/broadcasts/{id} [get]
func (h *BroadcastHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	b, err := h.getUC.Execute(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	WriteJSON(w, http.StatusOK, broadcastResponse(b))
}

// @Summary Create broadcast
// @Description Create a new broadcast. Supports JSON or multipart/form-data with optional audio/video upload.
// @Accept json
// @Accept multipart/form-data
// @Param title formData string true "Title"
// @Param broadcast_date formData string false "Broadcast date in RFC3339"
// @Param file formData file false "Audio or video file"
// @Success 201 {object} dto.BroadcastResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/broadcasts [post]
func (h *BroadcastHandler) Create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(64 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var b entity.Broadcast
		b.Title = r.FormValue("title")
		if v := r.FormValue("broadcast_date"); v != "" { if t, err := time.Parse(time.RFC3339, v); err == nil { b.BroadcastDate = t } }
		if b.ID == uuid.Nil { b.ID = uuid.New() }
		b.CreatedAt = time.Now()

		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			// choose type by extension
			typ := "audio"
			if strings.Contains(strings.ToLower(fh.Filename), ".mp4") || strings.Contains(strings.ToLower(fh.Filename), ".mov") || strings.Contains(strings.ToLower(fh.Filename), ".mkv") {
				typ = "video"
			}
			path, err := h.store.Save(file, fh.Filename, typ)
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			b.FilePath = &path
		}
		if err := h.createUC.Execute(r.Context(), &b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
		WriteJSON(w, http.StatusCreated, broadcastResponse(&b))
		return
	}

	var req dto.BroadcastCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	b := entity.Broadcast{ID: uuid.New(), Title: req.Title, BroadcastType: entity.BroadcastType(req.BroadcastType), ChannelName: req.ChannelName, CreatedAt: time.Now()}
	if req.BroadcastDate != nil { b.BroadcastDate = *req.BroadcastDate }
	if err := h.createUC.Execute(r.Context(), &b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusCreated, broadcastResponse(&b))
}

// @Summary Update broadcast
// @Description Update an existing broadcast. Supports JSON or multipart/form-data with optional audio/video upload.
// @Accept json
// @Accept multipart/form-data
// @Param id path string true "Broadcast ID"
// @Param title formData string false "Title"
// @Param broadcast_date formData string false "Broadcast date in RFC3339"
// @Param file formData file false "Audio or video file"
// @Success 200 {object} dto.BroadcastResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/broadcasts/{id} [put]
func (h *BroadcastHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(64 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var b entity.Broadcast
		b.ID = id
		b.Title = r.FormValue("title")
		if v := r.FormValue("broadcast_date"); v != "" { if t, err := time.Parse(time.RFC3339, v); err == nil { b.BroadcastDate = t } }

		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			typ := "audio"
			if strings.Contains(strings.ToLower(fh.Filename), ".mp4") || strings.Contains(strings.ToLower(fh.Filename), ".mov") || strings.Contains(strings.ToLower(fh.Filename), ".mkv") {
				typ = "video"
			}
			path, err := h.store.Save(file, fh.Filename, typ)
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			b.FilePath = &path
		}
		if err := h.updateUC.Execute(r.Context(), &b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
		WriteJSON(w, http.StatusOK, broadcastResponse(&b))
		return
	}

	var req dto.BroadcastCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	b := entity.Broadcast{ID: id, Title: req.Title, BroadcastType: entity.BroadcastType(req.BroadcastType), ChannelName: req.ChannelName}
	if req.BroadcastDate != nil { b.BroadcastDate = *req.BroadcastDate }
	if err := h.updateUC.Execute(r.Context(), &b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, broadcastResponse(&b))
}

// @Summary Delete broadcast
// @Description Delete an existing broadcast by id
// @Param id path string true "Broadcast ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/broadcasts/{id} [delete]
func (h *BroadcastHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }
	if err := h.deleteUC.Execute(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
