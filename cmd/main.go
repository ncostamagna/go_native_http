package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/ncostamagna/go_native_http/internal/user"
	"github.com/ncostamagna/go_native_http/pkg/bootstrap"
	"github.com/ncostamagna/go_native_http/pkg/handler"
)

func main() {

	_ = godotenv.Load()
	server := http.NewServeMux()

	db, err := bootstrap.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	logger := bootstrap.NewLogger()

	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()
	handler.NewUserHTTPServer(ctx, server, user.MakeEndpoints(service))

	port := os.Getenv("PORT")
	fmt.Println("Server started at port ", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}
