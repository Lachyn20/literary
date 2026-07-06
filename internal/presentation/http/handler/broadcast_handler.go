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
	svc   *broadcast.BroadcastService
	store repository.FileStorage
}

func NewBroadcastHandler(svc *broadcast.BroadcastService, store repository.FileStorage) *BroadcastHandler {
	return &BroadcastHandler{svc: svc, store: store}
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
	items, err := h.svc.List(r.Context())
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
	b, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	WriteJSON(w, http.StatusOK, broadcastResponse(b))
}

// @Summary Create broadcast
// @Description Create a new broadcast. Supports JSON or multipart/form-data with optional audio/video upload.
// @Accept json
// @Accept multipart/form-data
// @Param title formData string true "Title"
// @Param broadcast_type formData string true "Broadcast type: tv or radio"
// @Param channel_name formData string false "Channel name"
// @Param broadcast_date formData string false "Broadcast date in RFC3339"
// @Param file formData file false "Media file — tv: .mp4 .mov .mkv | radio: .mp3 .wav .aac"
// @Success 201 {object} dto.BroadcastResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/broadcasts [post]
func (h *BroadcastHandler) Create(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(64 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }
		var b entity.Broadcast
		b.Title = r.FormValue("title")

		broadcastType := r.FormValue("broadcast_type")
		if broadcastType != "tv" && broadcastType != "radio" {
			WriteError(w, http.StatusBadRequest, "broadcast_type must be 'tv' or 'radio'")
			return
		}
		b.BroadcastType = entity.BroadcastType(broadcastType)

		b.ChannelName = r.FormValue("channel_name")

		if v := r.FormValue("broadcast_date"); v != "" { if t, err := time.Parse(time.RFC3339, v); err == nil { b.BroadcastDate = t } }
		if b.ID == uuid.Nil { b.ID = uuid.New() }
		b.CreatedAt = time.Now()

		var savedPath string
		var savedFile bool

		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			lower := strings.ToLower(fh.Filename)
			isVideo := strings.HasSuffix(lower, ".mp4") || strings.HasSuffix(lower, ".mov") || strings.HasSuffix(lower, ".mkv")
			if broadcastType == "tv" {
				if !isAllowedExtension(fh.Filename, []string{".mp4", ".mov", ".mkv"}) {
					WriteError(w, http.StatusBadRequest, "broadcast_type 'tv' requires video file: .mp4, .mov, .mkv")
					return
				}
			} else if broadcastType == "radio" {
				if !isAllowedExtension(fh.Filename, []string{".mp3", ".wav", ".aac"}) {
					WriteError(w, http.StatusBadRequest, "broadcast_type 'radio' requires audio file: .mp3, .wav, .aac")
					return
				}
			}
			// choose type by extension
			typ := "audio"
			if isVideo {
				typ = "video"
			}
			savedPath, err := h.store.Save(file, fh.Filename, typ)
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
		
			savedFile = true
			b.FilePath = &savedPath
			b.FileType = entity.FileType(typ)

			fileType := "audio"
			if typ == "video" {
				fileType = "video"
			}
			b.FileType = entity.FileType(fileType)
		}
			
            if err := h.svc.Create(r.Context(), &b); err != nil {
			 if savedFile {
			    	_ = h.store.Remove(savedPath)
			    }
			    WriteError(w, http.StatusInternalServerError, err.Error())
			    return
	    	}
 
		WriteJSON(w, http.StatusCreated, broadcastResponse(&b))
		return
	}

	var req dto.BroadcastCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
	b := entity.Broadcast{ID: uuid.New(), Title: req.Title, BroadcastType: entity.BroadcastType(req.BroadcastType), ChannelName: req.ChannelName, CreatedAt: time.Now()}
	if req.BroadcastDate != nil { b.BroadcastDate = *req.BroadcastDate }
	if err := h.svc.Create(r.Context(), &b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusCreated, broadcastResponse(&b))
}

// @Summary Update broadcast
// @Description Update an existing broadcast. Supports JSON or multipart/form-data with optional audio/video upload.
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Broadcast ID"
// @Param title formData string false "Title"
// @Param broadcast_type formData string false "Broadcast type: tv or radio"
// @Param channel_name formData string false "Channel name"
// @Param broadcast_date formData string false "Broadcast date in RFC3339"
// @Param file formData file false "Media file — tv: .mp4, .mov, .mkv | radio: .mp3, .wav, .aac"
// @Success 200 {object} dto.BroadcastResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/broadcasts/{id} [put]
func (h *BroadcastHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil { WriteError(w, http.StatusBadRequest, "invalid id"); return }

	old , err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	b := *old

	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(64 << 20); err != nil { WriteError(w, http.StatusBadRequest, "invalid multipart"); return }

		if v := r.FormValue("title"); v != "" { b.Title = v }
		if v := r.FormValue("broadcast_type"); v != "" {
			if v != "tv" && v != "radio" {
				WriteError(w, http.StatusBadRequest, "broadcast_type must be 'tv' or 'radio'")
				return
			}
			b.BroadcastType = entity.BroadcastType(v)
		}
		if v := r.FormValue("channel_name"); v != "" { b.ChannelName = v }
		if v := r.FormValue("broadcast_date"); v != "" { if t, err := time.Parse(time.RFC3339, v); err == nil { b.BroadcastDate = t } }
		
		var savedFile bool
		var newFilePath string

		file, fh, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			lower := strings.ToLower(fh.Filename)
			isVideo := strings.Contains(lower, ".mp4") || strings.Contains(lower, ".mov") || strings.Contains(lower, ".mkv")
			isAudio := strings.Contains(lower, ".mp3") || strings.Contains(lower, ".wav") || strings.Contains(lower, ".aac")

			if !isVideo && !isAudio {
				WriteError(w, http.StatusBadRequest, "unsupported file type")
				return
			}

			currentType := string(b.BroadcastType)
        if currentType == "tv" {
			if !isVideo {
				WriteError(w, http.StatusBadRequest, "broadcast_type 'tv' requires video file: .mp4, .mov, .mkv")
				return
			}
		} else if currentType == "radio" {
			if !isAudio {
				WriteError(w, http.StatusBadRequest, "broadcast_type 'radio' requires audio file: .mp3, .wav, .aac")
				return
			}
			}

			typ := "audio"
			if isVideo {
				typ = "video"
			}
			newFilePath, err = h.store.Save(file, fh.Filename, typ)
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			savedFile = true
			b.FilePath = &newFilePath
			b.FileType = entity.FileType(typ)
		}

        if err := h.svc.Update(r.Context(), &b); err != nil {
			if savedFile {
				_ = h.store.Remove(newFilePath) // täze ýazylan faýly yzyna poz
			}
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Täze faýl ýazyldy we DB täzelendi -> köne faýly poz
		if savedFile && old.FilePath != nil && *old.FilePath != newFilePath {
			_ = h.store.Remove(*old.FilePath)
		}

		WriteJSON(w, http.StatusOK, broadcastResponse(&b))
		return
	}

	// === JSON bölegi - şol bir partial update ýörelgesi ===
	var req dto.BroadcastCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }

	if req.Title != "" { b.Title = req.Title }
	if req.BroadcastType != "" {
		if req.BroadcastType != "tv" && req.BroadcastType != "radio" {
			WriteError(w, http.StatusBadRequest, "broadcast_type must be 'tv' or 'radio'")
			return
		}
		b.BroadcastType = entity.BroadcastType(req.BroadcastType)
	}
	if req.ChannelName != "" {
		b.ChannelName = req.ChannelName
	}
	if req.BroadcastDate != nil {
		b.BroadcastDate = *req.BroadcastDate
	}

	if err := h.svc.Update(r.Context(), &b); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
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
	if err := h.svc.Delete(r.Context(), id); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
