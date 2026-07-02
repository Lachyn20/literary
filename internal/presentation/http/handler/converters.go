package handler

import (
	"github.com/hemra-siirow/literary/internal/domain/entity"
	"github.com/hemra-siirow/literary/internal/presentation/http/dto"
)

func bookResponse(b *entity.Book) dto.BookResponse {
	return dto.BookResponse{
		ID:                b.ID.String(),
		Title:             b.Title,
		BibliographicInfo: b.BibliographicInfo,
		CoverImagePath:    b.CoverImagePath,
		PDFPath:           b.PDFPath,
		PageCount:         b.PageCount,
		PublishedYear:     b.PublishedYear,
		CreatedAt:         b.CreatedAt,
	}
}

func categoryResponse(c *entity.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:   c.ID.String(),
		Name: c.Name,
		Slug: c.Slug,
	}
}

func categoryResponses(items []*entity.Category) []dto.CategoryResponse {
	resp := make([]dto.CategoryResponse, 0, len(items))
	for _, c := range items {
		resp = append(resp, categoryResponse(c))
	}
	return resp
}

func broadcastResponse(b *entity.Broadcast) dto.BroadcastResponse {
	return dto.BroadcastResponse{
		ID:            b.ID.String(),
		Title:         b.Title,
		BroadcastType: string(b.BroadcastType),
		ChannelName:   b.ChannelName,
		BroadcastDate: b.BroadcastDate,
		FilePath:      b.FilePath,
		CreatedAt:     b.CreatedAt,
	}
}

func filmResponse(f *entity.Film) dto.FilmResponse {
	return dto.FilmResponse{
		ID:              f.ID.String(),
		Title:           f.Title,
		FilmType:        string(f.FilmType),
		BasedOnScenario: f.BasedOnScenario,
		Director:        f.Director,
		ReleaseYear:     f.ReleaseYear,
		VideoPath:       f.VideoPath,
		CreatedAt:       f.CreatedAt,
	}
}

func personalLetterResponse(p *entity.PersonalLetter) dto.PersonalLetterResponse {
	return dto.PersonalLetterResponse{
		ID:            p.ID.String(),
		Title:         p.Title,
		Content:       p.Content,
		LetterDate:    p.LetterDate,
		ScanImagePath: p.ScanImagePath,
		CreatedAt:     p.CreatedAt,
	}
}

func photoArchiveResponse(p *entity.PhotoArchive) dto.PhotoArchiveResponse {
	return dto.PhotoArchiveResponse{
		ID:          p.ID.String(),
		Title:       p.Title,
		ImagePath:   p.ImagePath,
		Description: p.Description,
		TakenDate:   p.TakenDate,
		Category:    string(p.Category),
		CreatedAt:   p.CreatedAt,
	}
}

func theatreResponse(t *entity.TheatreProduction) dto.TheatreResponse {
	return dto.TheatreResponse{
		ID:           t.ID.String(),
		PlayTitle:    t.PlayTitle,
		TheatreName:  t.TheatreName,
		PremiereDate: t.PremiereDate,
		Notes:        t.Notes,
		CreatedAt:    t.CreatedAt,
	}
}

func translatedByAuthorResponse(t *entity.TranslatedByAuthor) dto.TranslatedByAuthorResponse {
	return dto.TranslatedByAuthorResponse{
		ID:                 t.ID.String(),
		OriginalAuthorName: t.OriginalAuthorName,
		OriginalLanguage:   t.OriginalLanguage,
		WorkTitle:          t.WorkTitle,
		Notes:              t.Notes,
	}
}

func translatedIntoLanguageResponse(t *entity.TranslatedIntoLanguage) dto.TranslatedIntoLanguageResponse {
	return dto.TranslatedIntoLanguageResponse{
		ID:             t.ID.String(),
		LanguageName:   t.LanguageName,
		TranslatorName: t.TranslatorName,
		WorkTitle:      t.WorkTitle,
		Notes:          t.Notes,
	}
}

func biographyResponse(b *entity.Biography) dto.BiographyResponse {
	return dto.BiographyResponse{
		ID:        b.ID.String(),
		Content:   b.Content,
		PhotoPath: b.PhotoPath,
		UpdatedAt: b.UpdatedAt,
	}
}

func workResponse(w *entity.Work) dto.WorkResponse {
	return dto.WorkResponse{
		ID:           w.ID.String(),
		Title:        w.Title,
		CategoryID:   w.CategoryID.String(),
		FilePath:     w.FilePath,
		Description:  w.Description,
		AudienceType: string(w.AudienceType),
		PublishYear:  w.PublishYear,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}

func workResponses(items []*entity.Work) []dto.WorkResponse {
	resp := make([]dto.WorkResponse, 0, len(items))
	for _, w := range items {
		resp = append(resp, workResponse(w))
	}
	return resp
}

func bookResponses(books []*entity.Book) []dto.BookResponse {
	resp := make([]dto.BookResponse, 0, len(books))
	for _, b := range books {
		resp = append(resp, bookResponse(b))
	}
	return resp
}

func broadcastResponses(items []*entity.Broadcast) []dto.BroadcastResponse {
	resp := make([]dto.BroadcastResponse, 0, len(items))
	for _, b := range items {
		resp = append(resp, broadcastResponse(b))
	}
	return resp
}

func filmResponses(items []*entity.Film) []dto.FilmResponse {
	resp := make([]dto.FilmResponse, 0, len(items))
	for _, f := range items {
		resp = append(resp, filmResponse(f))
	}
	return resp
}

func personalLetterResponses(items []*entity.PersonalLetter) []dto.PersonalLetterResponse {
	resp := make([]dto.PersonalLetterResponse, 0, len(items))
	for _, p := range items {
		resp = append(resp, personalLetterResponse(p))
	}
	return resp
}

func photoArchiveResponses(items []*entity.PhotoArchive) []dto.PhotoArchiveResponse {
	resp := make([]dto.PhotoArchiveResponse, 0, len(items))
	for _, p := range items {
		resp = append(resp, photoArchiveResponse(p))
	}
	return resp
}

func theatreResponses(items []*entity.TheatreProduction) []dto.TheatreResponse {
	resp := make([]dto.TheatreResponse, 0, len(items))
	for _, t := range items {
		resp = append(resp, theatreResponse(t))
	}
	return resp
}

func translatedByAuthorResponses(items []*entity.TranslatedByAuthor) []dto.TranslatedByAuthorResponse {
	resp := make([]dto.TranslatedByAuthorResponse, 0, len(items))
	for _, t := range items {
		resp = append(resp, translatedByAuthorResponse(t))
	}
	return resp
}

func translatedIntoLanguageResponses(items []*entity.TranslatedIntoLanguage) []dto.TranslatedIntoLanguageResponse {
	resp := make([]dto.TranslatedIntoLanguageResponse, 0, len(items))
	for _, t := range items {
		resp = append(resp, translatedIntoLanguageResponse(t))
	}
	return resp
}

func biographyResponses(items []*entity.Biography) []dto.BiographyResponse {
	resp := make([]dto.BiographyResponse, 0, len(items))
	for _, b := range items {
		resp = append(resp, biographyResponse(b))
	}
	return resp
}
