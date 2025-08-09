package repositories

import (
	"context"
	"database/sql"

	"github.com/luckydevil2007/audionotes/entities"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	var u entities.User
	err := repo.db.QueryRow("SELECT id, user_name, password_hash FROM users WHERE user_name = $1", login).Scan(&u.ID, &u.Username, &u.PasswordHash)

	return &u, err
}
