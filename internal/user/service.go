package user

import (
	"context"
	"log"

	"github.com/ncostamagna/go_native_http/internal/domain"
)

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error) {

	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {

	user, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}
