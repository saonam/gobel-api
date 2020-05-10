package usecases

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bmf-san/gobel-api/app/domain"
	"golang.org/x/crypto/bcrypt"
)

// A AuthInteractor is an interactor for authentication.
type AuthInteractor struct {
	AdminRepository AdminRepository
	JSONResponse    JSONResponse
	Logger          Logger
}

// HandleSignIn sign in with a credential.
func (ai *AuthInteractor) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	ai.Logger.LogAccess(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ai.Logger.LogError(err)
		ai.JSONResponse.Error500(w)
		return
	}

	var req RequestCredential
	err = json.Unmarshal(body, &req)
	if err != nil {
		ai.Logger.LogError(err)
		ai.JSONResponse.Error500(w)
		return
	}

	var admin domain.Admin
	admin, err = ai.AdminRepository.FindByCredential(req)
	if err != nil {
		ai.Logger.LogError(err)
		ai.JSONResponse.Error500(w)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		ai.Logger.LogError(err)
		ai.JSONResponse.Error403(w)
		return
	}

	token, err := ai.AdminRepository.SaveLoginSessionByID(admin.ID)
	if err != nil {
		ai.Logger.LogError(err)
		ai.JSONResponse.Error500(w)
		return
	}

	ai.JSONResponse.Success200(w, []byte(token))
	return
}

// HandleSignOut signs out.
func (ai *AuthInteractor) HandleSignOut(w http.ResponseWriter, r *http.Request) {
	ai.Logger.LogAccess(r)

	token := r.Header.Get("Authorization")

	if err := ai.AdminRepository.RemoveLoginSession(token); err != nil {
		ai.Logger.LogError(err)
		ai.JSONResponse.Error500(w)
		return
	}

	ai.JSONResponse.Success200(w, []byte(token))
	return
}
