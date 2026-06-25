package dto

type CategoryCreateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
