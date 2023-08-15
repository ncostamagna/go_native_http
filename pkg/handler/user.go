package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		
		
		if pathSize < 3 || pathSize > 4{
			InvalidMethod(w)
			return
		}

		params := make(map[string] string)
		
		if pathSize == 4 && path[2] != ""{
			params["userID"] = path[2]
		}
		ctx = context.WithValue(ctx, "params", params)

		tran := transport.New(w, r, ctx)
		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				tran.Server(
					transport.Endpoint(endpoint.GetAll),
					decodeGetAllUser,
					encodeResponse,
					encodeError)
				return
			case 4:
				tran.Server(
					nil,
					decoGetUser,
					encodeResponse,
					encodeError)
				return 
			}
			
		case http.MethodPost:
			switch pathSize {
			case 3:
				tran.Server(
					transport.Endpoint(endpoint.Create),
					decodeCreateUser,
					encodeResponse,
					encodeError)
					return
			}
			
		}

		InvalidMethod(w)
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
	fmt.Println(params)
	fmt.Println(params["userID"])
	return nil, fmt.Errorf("my error")
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