package transport

import (
	"context"
	"net/http"
	"strings"
)

type Trasport interface {
	Server(
		endpoint Endpoint,
		decode func(ctx context.Context, r *http.Request) (interface{}, error),
		encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
		encodeError func(ctx context.Context, err error, w http.ResponseWriter))
}

type Endpoint func(ctx context.Context, request interface{}) (interface{}, error)

type transport struct {
	w http.ResponseWriter
	r *http.Request
	ctx context.Context
}

func New(w http.ResponseWriter, r *http.Request, ctx context.Context) Trasport {
	return &transport{
		w: w,
		r: r,
		ctx: ctx,
	}
}

func (t *transport) Server(endpoint Endpoint,
	decode func(ctx context.Context, r *http.Request) (interface{}, error),
	encode func(ctx context.Context, w http.ResponseWriter, resp interface{}) error,
	encodeError func(ctx context.Context, err error, w http.ResponseWriter)){


		data, err := decode(t.ctx, t.r)
		if err != nil {
			encodeError(t.ctx, err, t.w)
			return
		}

		res, err := endpoint(t.ctx, data)
		if err != nil {
			encodeError(t.ctx, err, t.w)
			return
		}

		if err := encode(t.ctx, t.w, res); err != nil {
			encodeError(t.ctx, err, t.w)
			return
		}

}

func Clean(url string) ([]string, int) {
	
		if url[0] != '/' {
			url = "/" + url
		}

		if url[len(url)-1] != '/' {
			url =  url + "/"
		}

		parts := strings.Split(url, "/")

		return parts, len(parts)
}