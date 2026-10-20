package repositories

import (
	"context"
	"time"
	"todo_with_auth/internals/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user models.User) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := ` INSERT INTO users(email, password)
			   VALUES ($1, $2)
			   RETURNING id, email, created_at, updated_at
	`

	var createdUser models.User
	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&createdUser.Id,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := ` SELECT id, email, password, created_at, updated_at
			   FROM users
			   WHERE email = $1
	`

	var user models.User
	err := pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
