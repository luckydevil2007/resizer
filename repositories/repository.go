package repositories

import (
	"context"
	"database/sql"

	"github.com/luckydevil2007/go-lessons/entities"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveImage(ctx context.Context, image *entities.Image) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO images (image_title, image_path, owner_id) VALUES ($1, $2, $3)`,
		image.Title, image.Path, image.Owner)
	return err
}

func (r *Repository) DeleteImage(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE * FROM images WHERE id = $1`,
		id)
	return err
}

func (r *Repository) OpenImage(ctx context.Context, image *entities.Image) error {
	err := r.db.QueryRowContext(ctx,
		`SELECT image_title, image_path, owner_id FROM images WHERE ID = $1`, image.ID).Scan(
		&image.Title, &image.Path, &image.Owner)
	return err
}
