package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ncostamagna/go_native_http/internal/user"
	"github.com/ncostamagna/go_native_http/pkg/bootstrap"
	"github.com/ncostamagna/go_native_http/pkg/handler"
)

func main() {

	server := http.NewServeMux()

	db := bootstrap.NewDB()
	logger := bootstrap.NewLogger()

	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()
	handler.NewUserHTTPServer(ctx, server, user.MakeEndpoints(service))

	fmt.Println("Server started at port 8080")

	log.Fatal(http.ListenAndServe(":8080", server))
}