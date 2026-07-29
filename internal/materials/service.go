package materials

import (
	"context"

	myerrors "github.com/tojdor/nasrinShop/internal/myErrors"
)

type Service struct {
	storage Storer
}

func NewService(storage Storer) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Add(
	ctx context.Context,
	categorieID int,
	price int,
	imgURL string,
) (int, error) {

	if categorieID <= 0 || price <= 0 || imgURL == "" {
		return 0, myerrors.ErrBadRequest
	}

	id, err := s.storage.Add(ctx, categorieID, price, imgURL)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetByCategorieID(ctx context.Context, categorieID int) ([]Material, error) {
	if categorieID <= 0 {
		return nil, myerrors.ErrBadRequest
	}

	materials, err := s.storage.GetByCategorieID(ctx, categorieID)
	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return myerrors.ErrBadRequest
	}

	return s.storage.Delete(ctx, id)
}
