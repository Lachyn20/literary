package entity

import "github.com/google/uuid"

type AudienceType string

type BroadcastType string

type FileType string

type FilmType string

type PhotoCategory string

const (
	AudienceAdult    AudienceType = "adult"
	AudienceChildren AudienceType = "children"

	BroadcastTV    BroadcastType = "tv"
	BroadcastRadio BroadcastType = "radio"

	FileVideo FileType = "video"
	FileAudio FileType = "audio"

	FilmLiveAction FilmType = "film"
	FilmAnimation  FilmType = "animation"

	PhotoCategoryArchive PhotoCategory = "archive"
	PhotoCategoryPersonal PhotoCategory = "personal"
)

var (
	_ = uuid.UUID{}
)
