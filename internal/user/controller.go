package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ncostamagna/go_native_http/internal/domain"
)

// Endpoints struct
type (
	Controller func(w http.ResponseWriter, r *http.Request)

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

func MakeEndpoints(ctx context.Context, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAllUser(ctx, s, w)
		case http.MethodPost:
			decoder := json.NewDecoder(r.Body)
			var user domain.User
			if err := decoder.Decode(&user); err != nil {
				MsgResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			PostUser(ctx, s, w, user)
		default:
			InvalidMethod(w)
		}
	}
}

func GetAllUser(ctx context.Context, s Service, w http.ResponseWriter) {
	users, err := s.GetAll(ctx)
	if err != nil {
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	DataResponse(w, http.StatusOK, users)
}

func PostUser(ctx context.Context, s Service, w http.ResponseWriter, data interface{}) {
	req := data.(domain.User)

	if req.FirstName == "" {
		MsgResponse(w, http.StatusBadRequest, "first name is required")
		return
	}

	if req.LastName == "" {
		MsgResponse(w, http.StatusBadRequest, "last name is required")
		return
	}

	if req.Email == "" {
		MsgResponse(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)

	if err != nil {
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	DataResponse(w, http.StatusCreated, user)

}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "message":"method doesn't exist"}`, status)
}

func MsgResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "message":"%s"}`, status, message)
}

func DataResponse(w http.ResponseWriter, status int, data interface{}) {
	value, err := json.Marshal(data)
	if err != nil {
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "data":%s}`, status, string(value))
}

/*
//MakeEndpoints handler endpoints
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
}*/
