package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ncostamagna/go_native_http/internal/user"
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

		params := make(map[string] string)
		
		if pathSize == 4 && path[2] != ""{
			params["userID"] = path[2]
		}

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
		}else{
			InvalidMethod(w)
		}

	}

}

func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {

	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request format: '%v'", err.Error())
	}

	return req, nil
}

func decoGetUser(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string] string)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, err
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

	params := ctx.Value("params").(map[string] string)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		return nil, err
	}

	req.UserID = userID
	
	return req, nil
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	status := http.StatusOK
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "data":%s}`, status, data)

	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {

	status := http.StatusInternalServerError
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "message":"%s"}`, status, err.Error())

}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status":%d, "message":"method doesn't exist"}`, status)
}