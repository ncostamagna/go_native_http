package user

import (
	"context"
	"database/sql"
	"log"

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
		db  *sql.DB
		log *log.Logger
	}
)

// NewRepo is a repositories handler
func NewRepo(db *sql.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {

	sql := "INSERT INTO users(first_name, last_name, email) VALUES(?,?,?)"
	res, err := r.db.Exec(sql,user.FirstName, user.LastName, user.Email)
if err != nil {
	r.log.Println(err.Error())
	return  err
}
id, err := res.LastInsertId()
if err != nil {
	r.log.Println(err.Error())
	return  err
}
user.ID = uint64(id)
	r.log.Println("user created with id: ", id)
	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {

	var users []domain.User
	sql := "select id, first_name, last_name, email from users"
	r.log.Println(sql)
	rows, err := r.db.Query(sql)
	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
			r.log.Println(err.Error())
			return nil, err
		}
		users = append(users, u)
	}
	r.log.Println("user get all: ", len(users))
	return users, nil
}

func (r *repo) Get(ctx context.Context, userID uint64) (*domain.User, error) {
/*
	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
		return v.ID == userID
	})

	if index < 0 {
		return nil, ErrNotFound{userID}
	}
	*/
	return nil, nil
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