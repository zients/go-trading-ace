package repositories

import (
	"database/sql"
	"fmt"
	"time"
	"trading-ace/entities"
)

type IUserRepository interface {
	Create(address string) (*entities.User, error)
	FindByID(id int) (*entities.User, error)
	UpdateAddressById(id int, newAddress string) error
	DeleteByID(id int) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(address string) (*entities.User, error) {
	query := `
		INSERT INTO users (address, created_at, updated_at)
		VALUES ($1, $2, $3) RETURNING id, address, created_at, updated_at
	`

	var user entities.User
	err := u.db.QueryRow(query, address, time.Now().UTC(), time.Now().UTC()).Scan(&user.ID, &user.Address, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindByID(id int) (*entities.User, error) {
	query := "SELECT id, address, created_at, updated_at FROM users WHERE id = $1"
	var user entities.User
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.Address, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) UpdateAddressById(id int, newAddress string) error {
	query := "UPDATE users SET address = $1, updated_at = $2 WHERE id = $3"
	_, err := u.db.Exec(query, newAddress, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) DeleteByID(id int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := u.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
