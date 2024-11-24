package repositories

import (
	"database/sql"
	"time"
	"trading-ace/entities"
)

type IUserRepository interface {
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(address string) (*entities.User, error) {
	// 插入資料
	query := `
		INSERT INTO users (address, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING id, address, created_at, updated_at
	`

	var user entities.User
	err := u.db.QueryRow(query, address, time.Now(), time.Now()).Scan(&user.ID, &user.Address, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
