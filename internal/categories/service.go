package categories

import (
	"context"
	"strings"

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

func (s *Service) Add(ctx context.Context, name string) (int, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, myerrors.ErrBadRequest
	}

	id, err := s.storage.Add(ctx, name)
	return id, err
}

func (s *Service) GetAll(ctx context.Context) ([]Category, error) {
	return s.storage.GetAll(ctx)
}

func (s *Service) GetIDByName(ctx context.Context, name string) (int, error) {
	if name == "" {
		return 0, myerrors.ErrBadRequest
	}
	return s.storage.GetIDByName(ctx, name)
}

func (s *Service) Delete(ctx context.Context, name string) error {
	if len(name) <= 0 {
		return myerrors.ErrBadRequest
	}

	return s.storage.Delete(ctx, name)
}
