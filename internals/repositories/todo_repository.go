package repositories

import (
	"context"
	"fmt"
	"time"
	"todo_with_auth/internals/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userId string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `INSERT INTO todos (title, completed, user_id)
			   VALUES ($1, $2, $3)
			   RETURNING id, title, completed, created_at, updated_at, user_id
	`

	var todo models.Todo

	err := pool.QueryRow(ctx, query, title, completed, userId).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetTodos(pool *pgxpool.Pool, userId string) ([]models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, completed, created_at, updated_at, user_id
			 FROM todos
			 WHERE user_id = $1
			 ORDER BY created_at DESC 
	`

	var todos []models.Todo = []models.Todo{}

	rows, err := pool.Query(ctx, query, userId)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo

		err := rows.Scan(
			&todo.Id,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserId,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}
	return todos, nil
}

func GetTodoById(pool *pgxpool.Pool, id int, userId string) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	var todo models.Todo

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := ` SELECT id, title, completed, created_at, updated_at, user_id
		       FROM todos
			   WHERE id = $1 AND user_id = $2
	`

	err := pool.QueryRow(ctx, query, id, userId).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool) (*models.Todo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todo models.Todo

	query := `UPDATE todos SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $3
			  RETURNING id, title, completed, created_at, updated_at, user_id
	`

	err := pool.QueryRow(ctx, query, title, completed, id).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int, userId string) error {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := ` DELETE FROM todos
			   WHERE id = $1 AND user_id = $2
	`

	commandTag, err := pool.Exec(ctx, query, id, userId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Todo with the given id not found")
	}
	return nil
}
