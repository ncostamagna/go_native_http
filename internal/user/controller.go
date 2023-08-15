package user

import (
	"context"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
	}
}


func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)
		if err != nil {

			return nil, err
		}

		return user, nil

	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		return users, nil
	}
}
