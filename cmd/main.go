package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ncostamagna/go_native_http/internal/domain"
	"github.com/ncostamagna/go_native_http/internal/user"
)

func main() {

	server := http.NewServeMux()

	db := user.DB{
		Users: []domain.User{{
			ID:        1,
			FirstName: "Nahuel",
			LastName:  "Costamagna",
			Email:     "nahuel@domain.com",
		}, {
			ID:        2,
			FirstName: "Eren",
			LastName:  "Jaeger",
			Email:     "eren@domain.com",
		}, {
			ID:        3,
			FirstName: "Poca",
			LastName:  "Costamagna",
			Email:     "poca@domain.com",
		}},
		MaxUserID: 3,
	}
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()
	server.HandleFunc("/user", user.MakeEndpoints(ctx, service))

	//handler.NewUserHTTPServer(ctx, server, userEndpoint)

	fmt.Println("Server started at port 8080")

	log.Fatal(http.ListenAndServe(":8080", server))
}
