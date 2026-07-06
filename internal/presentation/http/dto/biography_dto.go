package dto

import "time"

type BiographyResponse struct {
    ID        string     `json:"id"`
    PhotoPath *string    `json:"photo_path,omitempty"`
    UpdatedAt time.Time  `json:"updated_at"`
}
type BiographyUpdateRequest struct {
    // photo handled as multipart only
}

type BiographyCreateRequest struct {
    // photo handled as multipart only
}

type BiographyEventResponse struct {
    ID            string    `json:"id"`
    BiographyID   string    `json:"biography_id"`
    Year          int       `json:"year"`
    TitleTk       *string   `json:"title_tk,omitempty"`
    TitleRu       *string   `json:"title_ru,omitempty"`
    TitleEn       *string   `json:"title_en,omitempty"`
    DescriptionTk *string   `json:"description_tk,omitempty"`
    DescriptionRu *string   `json:"description_ru,omitempty"`
    DescriptionEn *string   `json:"description_en,omitempty"`
    SortOrder     int       `json:"sort_order"`
    CreatedAt     time.Time `json:"created_at"`
}

type BiographyEventCreateRequest struct {
    Year          int    `json:"year" validate:"required,min=1800,max=2100"`
    TitleTk       string `json:"title_tk" validate:"required"`
    TitleRu       string `json:"title_ru"`
    TitleEn       string `json:"title_en"`
    DescriptionTk string `json:"description_tk"`
    DescriptionRu string `json:"description_ru"`
    DescriptionEn string `json:"description_en"`
    SortOrder     int    `json:"sort_order"`
}

type BiographyEventUpdateRequest struct {
    Year          *int   `json:"year"`
    TitleTk       string `json:"title_tk"`
    TitleRu       string `json:"title_ru"`
    TitleEn       string `json:"title_en"`
    DescriptionTk string `json:"description_tk"`
    DescriptionRu string `json:"description_ru"`
    DescriptionEn string `json:"description_en"`
    SortOrder     *int   `json:"sort_order"`
}