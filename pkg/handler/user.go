package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ncostamagna/go_native_http/internal/user"
	"github.com/ncostamagna/go_native_http/pkg/response"
	"github.com/ncostamagna/go_native_http/pkg/transport"
)

func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoint user.Endpoints) {

	router.HandleFunc("/users/", UserServer(ctx, endpoint))

}

func UserServer(ctx context.Context, endpoint user.Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		log.Println(r.Method, ": ", url)
		path, pathSize := transport.Clean(url)

		params := make(map[string]string)

		if pathSize == 4 && path[2] != "" {
			params["userID"] = path[2]
		}

		params["token"] = r.Header.Get("Authorization")

		tran := transport.New(w, r, context.WithValue(ctx, "params", params))

		var end transport.Endpoint
		var deco func(ctx context.Context, r *http.Request) (interface{}, error)

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				end = transport.Endpoint(endpoint.GetAll)
				deco = decodeGetAllUser
			case 4:
				end = transport.Endpoint(endpoint.Get)
				deco = decoGetUser
			}

		case http.MethodPost:
			switch pathSize {
			case 3:
				end = transport.Endpoint(endpoint.Create)
				deco = decodeCreateUser
			}
		case http.MethodPatch:
			switch pathSize {
			case 4:
				end = transport.Endpoint(endpoint.Update)
				deco = decoUpdateUser
			}
		}

		if end != nil && deco != nil {
			tran.Server(
				end,
				deco,
				encodeResponse,
				encodeError)
		} else {
			InvalidMethod(w)
		}

	}

}

func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string]string)
	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decoGetUser(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string]string)

	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return user.GetReq{
		UserID: userID,
	}, nil
}

func decoUpdateUser(ctx context.Context, r *http.Request) (interface{}, error) {

	var req user.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request format: '%v'", err.Error())
	}

	params := ctx.Value("params").(map[string]string)

	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, err
	}

	req.UserID = userID

	return req, nil
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string]string)

	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	return nil, nil
}

func tokenVerify(token string) error {
	if os.Getenv("TOKEN") != token {
		return errors.New("invalid token")
	}
	return nil
}
func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {

	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)

}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)

	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "message":"method doesn't exist"}`, status)
}
