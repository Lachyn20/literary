package repository

import (
	"io"
)

// FileStorage defines storage operations for uploaded files.
type FileStorage interface {
	// Save saves the provided file stream under the given filename and type (e.g. "book","video","image").
	// Returns the relative path where the file is stored.
	Save(r io.Reader, filename string, typ string) (string, error)
	// Remove deletes a stored file by relative path.
	Remove(path string) error
}
