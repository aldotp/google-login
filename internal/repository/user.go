package repository

import (
	"database/sql"
	"errors"

	"github.com/aldotp/golang-login-with-google/internal/model"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	FindUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, email, role, provider FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.Role, &user.Provider); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (id, name, email, role, password, provider, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.Role, user.Password, user.Provider, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
