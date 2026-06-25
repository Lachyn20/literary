package dto

import "time"

type FilmCreateRequest struct {
    Title           string  `json:"title" validate:"required,min=1,max=255"`
    FilmType        string  `json:"film_type" validate:"required,oneof=film animation"`
    Director        *string `json:"director"`
    ReleaseYear     *int    `json:"release_year"`
    BasedOnScenario bool    `json:"based_on_scenario"`
    // file handled as multipart
}

type FilmResponse struct {
    ID              string     `json:"id"`
    Title           string     `json:"title"`
    FilmType        string     `json:"film_type"`
    BasedOnScenario bool       `json:"based_on_scenario"`
    Director        *string    `json:"director,omitempty"`
    ReleaseYear     *int       `json:"release_year,omitempty"`
    VideoPath       *string    `json:"video_path,omitempty"`
    CreatedAt       time.Time  `json:"created_at"`
}
