package materials

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/tojdor/nasrinShop/internal/myErrors"
)

type Material struct {
	ID         int
	CategoryID int
	Price      string
	ImgURL     string
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
	GetByCategorieID(
		ctx context.Context,
		categoryID int,
	) ([]Material, error)
	Delete(
		ctx context.Context,
		id int,
	) error
}

func (s *Storage) Add(
	ctx context.Context,
	categorieID int,
	price int,
	imgURL string,
) (int, error) {

	var id int
	err := s.pool.QueryRow(ctx,
		"INSERT INTO materials(categorie_id, price, image_url) VALUES $1, $2, $3 RETURNING id",
		categorieID,
		price,
		imgURL,
	).Scan(&id)

	return id, err
}

func (s *Storage) GetByCategorieID(ctx context.Context, categorieID int) ([]Material, error) {
	materials := make([]Material, 0)

	rows, err := s.pool.Query(
		ctx,
		"SELECT id, category_id, price, image_url FROM materials WHERE categorie_id = $1",
		categorieID,
	)
	if err != nil {
		return nil, err
	}

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
