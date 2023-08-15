package bootstrap

import (
	"log"
	"os"

	"github.com/ncostamagna/go_native_http/internal/domain"
	"github.com/ncostamagna/go_native_http/internal/user"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func NewDB() user.DB {
	return user.DB{
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
}