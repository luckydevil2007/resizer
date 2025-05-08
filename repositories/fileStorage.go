package repositories

import (
	"context"
	"database/sql"
	"os"

	"github.com/luckydevil2007/go-lessons/entities"
)

type FileStorage struct {
	db       *sql.DB
	location string
}

func NewFileStorage(db *sql.DB, location string) *FileStorage {
	return &FileStorage{db: db, location: location}
}

func (r *FileStorage) Save(ctx context.Context, image *entities.Image) error {
	fullPath := r.location + "/" + image.Title
	err := os.WriteFile(fullPath, image.Data, 0644)
	if err == nil {
		image.Path = fullPath
	}
	return err
}

func (r *FileStorage) Delete(ctx context.Context, image *entities.Image) error {
	return os.Remove(image.Path)
}

func (r *FileStorage) Open(ctx context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}
