package repositories

import (
	"context"
	"database/sql"
	"os"

	"github.com/luckydevil2007/audionotes/entities"
)

type FileStorage struct {
	db       *sql.DB
	location string
}

func NewFileStorage(db *sql.DB, location string) *FileStorage {
	return &FileStorage{db: db, location: location}
}

func (r *FileStorage) Save(ctx context.Context, note *entities.Note) error {
	fullPath := r.location + "/" + note.Title
	err := os.WriteFile(fullPath, note.Data, 0644)
	if err == nil {
		note.Path = fullPath
	}
	return err
}

func (r *FileStorage) Delete(ctx context.Context, note *entities.Note) error {
	err := os.Remove(note.Path)
	if err != nil {
		if note.Prev != nil {
			note.Prev.Next = note.Next
		}
	}
	return err
}

func (r *FileStorage) Open(ctx context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}
