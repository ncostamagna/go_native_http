package user

import (
	"context"
	"log"
	"slices"

	"github.com/ncostamagna/go_native_http/internal/domain"
)

type DB struct {
	Users     []domain.User
	MaxUserID uint64
}

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, userID uint64) (*domain.User, error)
		Update(ctx context.Context, userID uint64, firstName, lastName, email *string) error
	}

	repo struct {
		db  DB
		log *log.Logger
	}
)

// NewRepo is a repositories handler
func NewRepo(db DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {

	r.db.MaxUserID++
	user.ID = r.db.MaxUserID
	r.db.Users = append(r.db.Users, *user)
	r.log.Println("repository create")
	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	r.log.Println("repository get all")
	return r.db.Users, nil
}

func (r *repo) Get(ctx context.Context, userID uint64) (*domain.User, error) {

	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
		return v.ID == userID
	})

	if index < 0 {
		return nil, ErrNotFound{userID}
	}
	return &r.db.Users[index], nil
}

func (r *repo) Update(ctx context.Context, userID uint64, firstName, lastName, email *string) error {
	user, err := r.Get(ctx, userID)
	if err != nil {
		return err
	}

	if firstName != nil {
		user.FirstName = *firstName
	}

	if lastName != nil {
		user.LastName = *lastName
	}

	if email != nil {
		user.Email = *email
	}

	return nil
}