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
	"github.com/hemra-siirow/literary/internal/usecase/book"
)

type BookHandler struct {
	svc   *book.BookService
	store repository.FileStorage
}

func NewBookHandler(svc *book.BookService, store repository.FileStorage) *BookHandler {
	return &BookHandler{svc: svc, store: store}
}

func (h *BookHandler) RegisterRoutes(r chi.Router) {
	r.Get("/api/books", h.List)
	r.Get("/api/books/{id}", h.Get)
	r.Post("/api/books", h.Create)
	r.Put("/api/books/{id}", h.Update)
	r.Delete("/api/books/{id}", h.Delete)
}

// @Summary List books
// @Description List all books with pagination support (limit, offset/page)
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param page query int false "Page number (1-indexed, alternative to offset)"
// @Success 200 {object} dto.BookListResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/books [get]
func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	paginationParams := pagination.Parse(r)
	
	books, total, err := h.svc.List(r.Context(), paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if books == nil {
		books = []*entity.Book{}
	}
	
	paginationInfo := pagination.NewInfo(paginationParams.Limit, paginationParams.Offset, total)
	
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "ok",
		"data":         bookResponses(books),
		"pagination":   paginationInfo,
	})
}

// @Summary Get book
// @Description Get a single book by id
// @Param id path string true "Book ID"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 404 {object} handler.JSONResponse
// @Router /api/books/{id} [get]
func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	b, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, bookResponse(b))
}

// @Summary Create book
// @Description Create a new book. Supports JSON or multipart/form-data with optional cover image and PDF upload.
// @Accept json
// @Accept multipart/form-data
// @Param title formData string true "Title"
// @Param bibliographic_info formData string false "Bibliographic info"
// @Param page_count formData int false "Page count"
// @Param published_year formData int false "Published year"
// @Param cover formData file false "Cover image — accepted: .jpg, .jpeg, .png, .webp"
// @Param pdf formData file false "Book PDF — accepted: .pdf"
// @Success 201 {object} dto.BookResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/books [post]
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	// support multipart upload or JSON
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid multipart")
			return
		}
		var b entity.Book
		b.Title = r.FormValue("title")
		if v := r.FormValue("bibliographic_info"); v != "" {
			b.BibliographicInfo = &v
		}
		if v := r.FormValue("page_count"); v != "" {
			if i, err := strconv.Atoi(v); err == nil {
				b.PageCount = &i
			}
		}
		if v := r.FormValue("published_year"); v != "" {
			if i, err := strconv.Atoi(v); err == nil {
				b.PublishedYear = &i
			}
		}
		if b.ID == uuid.Nil {
			b.ID = uuid.New()
		}
		b.CreatedAt = time.Now()

		var coverPath, pdfPath string
		var savedCover, savedPDF bool

		// handle cover
		cover, ch, err := r.FormFile("cover")
		if err == nil {
			defer cover.Close()
			if !isAllowedExtension(ch.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}) {
				WriteError(w, http.StatusBadRequest, "only .jpg, .jpeg, .png, .webp files are allowed for cover")
				return
			}
			coverPath, err := h.store.Save(cover, ch.Filename, "image")
			if err != nil {
				WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			savedCover = true
			b.CoverImagePath = &coverPath
		}
		// handle pdf
		pdf, ph, err := r.FormFile("pdf")
		if err == nil {
			defer pdf.Close()
			if !isAllowedExtension(ph.Filename, []string{".pdf"}) {
				WriteError(w, http.StatusBadRequest, "only .pdf files are allowed for book")
				return
			}
			pdfPath, err := h.store.Save(pdf, ph.Filename, "book")
			if err != nil {
				if savedCover {
					_ = h.store.Remove(coverPath)
				}
				WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			savedPDF = true
			b.PDFPath = &pdfPath
		}

		if err := h.svc.Create(r.Context(), &b); err != nil {
			if savedCover {
				_ = h.store.Remove(coverPath)
			}
			if savedPDF {
				_ = h.store.Remove(pdfPath)
			}
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		WriteJSON(w, http.StatusCreated, bookResponse(&b))
		return
	}

	var req dto.BookCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := validation.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	b := entity.Book{ID: uuid.New(), Title: req.Title, BibliographicInfo: req.BibliographicInfo, PageCount: req.PageCount, PublishedYear: req.PublishedYear, CreatedAt: time.Now()}
	if err := h.svc.Create(r.Context(), &b); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, bookResponse(&b))
}

// @Summary Update book
// @Description Update an existing book. Supports JSON or multipart/form-data with optional cover image and PDF upload.
// @Accept json
// @Accept multipart/form-data
// @Param id path string true "Book ID"
// @Param title formData string false "Title"
// @Param bibliographic_info formData string false "Bibliographic info"
// @Param page_count formData int false "Page count"
// @Param published_year formData int false "Published year"
// @Param cover formData file false "Cover image — accepted: .jpg, .jpeg, .png, .webp"
// @Param pdf formData file false "Book PDF — accepted: .pdf"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/books/{id} [put]
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	// support multipart or json
	if ct := r.Header.Get("Content-Type"); ct != "" && strings.HasPrefix(ct, "multipart/") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid multipart")
			return
		}

		old, err := h.svc.GetByID(r.Context(), id)
		if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
		updated := *old

		if v := r.FormValue("title"); v != "" { updated.Title = v }
		if v := r.FormValue("bibliographic_info"); v != "" { updated.BibliographicInfo = &v }
		if v := r.FormValue("page_count"); v != "" {
			if i, err := strconv.Atoi(v); err == nil { updated.PageCount = &i }
		}
		if v := r.FormValue("published_year"); v != "" {
			if i, err := strconv.Atoi(v); err == nil { updated.PublishedYear = &i }
		}

		prevCover := old.CoverImagePath
		prevPDF := old.PDFPath
		var savedCoverPath, savedPDFPath string
		savedCover := false
		savedPDF := false

		cover, ch, err := r.FormFile("cover")
		if err == nil {
			defer cover.Close()
			if !isAllowedExtension(ch.Filename, []string{".jpg", ".jpeg", ".png", ".webp"}) {
				WriteError(w, http.StatusBadRequest, "only .jpg, .jpeg, .png, .webp files are allowed for cover")
				return
			}
			coverpath, err := h.store.Save(cover, ch.Filename, "image")
			if err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }
			savedCover = true
			savedCoverPath = coverpath
			updated.CoverImagePath = &coverpath
		}
		pdf, ph, err := r.FormFile("pdf")
		if err == nil {
			defer pdf.Close()
			if !isAllowedExtension(ph.Filename, []string{".pdf"}) {
				WriteError(w, http.StatusBadRequest, "only .pdf files are allowed for book")
				return
			}
			path, err := h.store.Save(pdf, ph.Filename, "book")
			if err != nil {
				if savedCover { _ = h.store.Remove(savedCoverPath) }
				WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			savedPDF = true
			savedPDFPath = path
			updated.PDFPath = &path
		}

		if err := h.svc.Update(r.Context(), &updated); err != nil {
			if savedCover { _ = h.store.Remove(savedCoverPath) }
			if savedPDF { _ = h.store.Remove(savedPDFPath) }
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// cleanup old files after successful update
		if savedCover && prevCover != nil && (updated.CoverImagePath == nil || *prevCover != *updated.CoverImagePath) {
			_ = h.store.Remove(*prevCover)
		}
		if savedPDF && prevPDF != nil && (updated.PDFPath == nil || *prevPDF != *updated.PDFPath) {
			_ = h.store.Remove(*prevPDF)
		}

		WriteJSON(w, http.StatusOK, bookResponse(&updated))
		return
	}

	var req dto.BookCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w, http.StatusBadRequest, "invalid payload"); return }
	if err := validation.Struct(req); err != nil { WriteError(w, http.StatusBadRequest, err.Error()); return }

	old, err := h.svc.GetByID(r.Context(), id)
	if err != nil { WriteError(w, http.StatusNotFound, err.Error()); return }
	updated := *old
	if req.Title != "" { updated.Title = req.Title }
	if req.BibliographicInfo != nil { updated.BibliographicInfo = req.BibliographicInfo }
	if req.PageCount != nil { updated.PageCount = req.PageCount }
	if req.PublishedYear != nil { updated.PublishedYear = req.PublishedYear }

	if err := h.svc.Update(r.Context(), &updated); err != nil { WriteError(w, http.StatusInternalServerError, err.Error()); return }
	WriteJSON(w, http.StatusOK, bookResponse(&updated))
}

// @Summary Delete book
// @Description Delete an existing book by id
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.JSONResponse
// @Failure 500 {object} handler.JSONResponse
// @Router /api/books/{id} [delete]
func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"deleted": id.String()})
}
