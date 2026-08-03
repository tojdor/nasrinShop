package categories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/tojdor/nasrinShop/internal/myErrors"
)

type Category struct {
	ID   int `json:"id"`
	Name string `json:"name"`
}

type Storage struct {
	pool *pgxpool.Pool
}

type Storer interface {
	Add(ctx context.Context, name string) (int, error)
	GetAll(ctx context.Context) ([]Category, error)
	GetIDByName(ctx context.Context, name string) (int, error)
	Delete(ctx context.Context, name string) error
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) Add(ctx context.Context, name string) (int, error) {
	var id int
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO categories (name) VALUES ($1) RETURNING id",
		name,
	).Scan(&id)
	if err!=nil{
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505"{
			return 0, myerrors.ErrConflict
		}
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetIDByName(ctx context.Context, name string) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx,
		"SELECT id FROM categories WHERE name = $1",
		name).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, myerrors.ErrNotFound
		}
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]Category, error) {
	categories := make([]Category, 0)

	rows, err := s.pool.Query(ctx, "SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *Storage) Delete(ctx context.Context, name string) error {
	result, err := s.pool.Exec(ctx, "DELETE FROM categories WHERE name = $1", name)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
