package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

var (
	ErrInvalidExtension = errors.New("invalid file extension")
	ErrInvalidMime      = errors.New("invalid mime type")
	ErrTooLarge         = errors.New("file too large")
)

// size limits in bytes
var sizeLimits = map[string]int64{
	"book": 100 * 1024 * 1024, // 100MB
	"video": 2 * 1024 * 1024 * 1024, // 2GB
	"image": 20 * 1024 * 1024, // 20MB
	"audio": 500 * 1024 * 1024, // 500MB
	"scan": 25 * 1024 * 1024, // 25MB for scans
}

var allowedExt = map[string][]string{
	"book": {".pdf"},
	"video": {".mp4", ".mov", ".mkv"},
	"image": {".jpg", ".jpeg", ".png", ".webp"},
	"audio": {".mp3", ".wav", ".m4a"},
	"scan": {".jpg", ".jpeg", ".png", ".pdf"},
}

var allowedMimePrefixes = map[string][]string{
	"book": {"application/pdf"},
	"video": {"video/"},
	"image": {"image/"},
	"audio": {"audio/"},
	"scan": {"image/", "application/pdf"},
}

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) Remove(path string) error {
	full := filepath.Join(s.basePath, path)
	return os.Remove(full)
}

func (s *LocalStorage) Save(r io.Reader, filename string, typ string) (string, error) {
	limit, ok := sizeLimits[typ]
	if !ok {
		return "", fmt.Errorf("unknown type: %s", typ)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "", ErrInvalidExtension
	}

	allowed, ok := allowedExt[typ]
	if ok {
		valid := false
		for _, e := range allowed {
			if e == ext {
				valid = true
				break
			}
		}
		if !valid {
			return "", ErrInvalidExtension
		}
	}

	// read first 512 bytes to detect mime type
	buf := make([]byte, 512)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", err
	}
	head := buf[:n]
	mimeType := http.DetectContentType(head)

	// validate mime prefix
	mimes := allowedMimePrefixes[typ]
	validMime := false
	for _, p := range mimes {
		if strings.HasPrefix(mimeType, p) {
			validMime = true
			break
		}
	}
	if !validMime {
		return "", ErrInvalidMime
	}

	// prepare reader that includes the head and the remaining stream
	body := io.MultiReader(bytes.NewReader(head), r)

	// enforce size limit with LimitReader
	limited := io.LimitReader(body, limit+1)

	// generate filename
	id := uuid.New()
	safeName := fmt.Sprintf("%s%s", id.String(), ext)
	t := time.Now()
	year := fmt.Sprintf("%d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	dir := filepath.Join(typ, year, month)
	fullDir := filepath.Join(s.basePath, dir)
	if err := os.MkdirAll(fullDir, 0o755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(fullDir, safeName)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	written, err := io.Copy(f, limited)
	if err != nil {
		return "", err
	}
	if written > limit {
		// remove partial file
		f.Close()
		os.Remove(fullPath)
		return "", ErrTooLarge
	}

	rel := filepath.Join(dir, safeName)
	// normalize to forward slashes for DB storage
	rel = filepath.ToSlash(rel)
	return rel, nil
}

// ensure LocalStorage implements repository.FileStorage
var _ repository.FileStorage = (*LocalStorage)(nil)
