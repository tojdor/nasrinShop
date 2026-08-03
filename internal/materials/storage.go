package materials

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/tojdor/nasrinShop/internal/myErrors"
)

type Material struct {
	ID         int `json:"id"`
	CategoryID int `json:"category_id"`
	Price      int `json:"price"`
	ImgURL     string `json:"image_url"`
}

type Storage struct {
	pool *pgxpool.Pool
}

type Storer interface {
	Add(
		ctx context.Context,
		categoryID int,
		price int,
		imgURL string,
	) (int, error)
	GetByCategoryID(
		ctx context.Context,
		categoryID int,
	) ([]Material, error)
	Delete(
		ctx context.Context,
		id int,
	) error
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) Add(
	ctx context.Context,
	categoryID int,
	price int,
	imgURL string,
) (int, error) {

	var id int
	err := s.pool.QueryRow(ctx,
		"INSERT INTO materials (category_id, price, image_url) VALUES ($1, $2, $3) RETURNING id",
		categoryID,
		price,
		imgURL,
	).Scan(&id)

	return id, err
}

func (s *Storage) GetByCategoryID(ctx context.Context, categoryID int) ([]Material, error) {
	materials := make([]Material, 0)

	rows, err := s.pool.Query(
		ctx,
		"SELECT id, category_id, price, image_url FROM materials WHERE category_id = $1",
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var material Material
		err := rows.Scan(
			&material.ID,
			&material.CategoryID,
			&material.Price,
			&material.ImgURL,
		)
		if err != nil {
			return nil, err
		}

		materials = append(materials, material)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return materials, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	result, err := s.pool.Exec(ctx, "DELETE FROM materials WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
