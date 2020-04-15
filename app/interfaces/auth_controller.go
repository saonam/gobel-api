package interfaces

import (
	"database/sql"
	"net/http"

	"github.com/bmf-san/gobel-api/app/usecases"
	"github.com/go-redis/redis/v7"
)

// An AuthController is a controller for an authentication.
type AuthController struct {
	AuthInteractor usecases.AuthInteractor
}

// NewAuthController creates an AuthController.
func NewAuthController(connMySQL *sql.DB, connRedis *redis.Client, logger usecases.Logger) *AuthController {
	return &AuthController{
		AuthInteractor: usecases.AuthInteractor{
			AdminRepository: &AdminRepository{
				ConnMySQL: connMySQL,
				ConnRedis: connRedis,
			},
			JSONResponse: &JSONResponse{},
			Logger:       logger,
		},
	}
}

// SignIn sign in with a credential.
func (ac *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	ac.AuthInteractor.HandleSignIn(w, r)
	return
}
