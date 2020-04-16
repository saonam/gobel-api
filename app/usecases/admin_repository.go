package usecases

import "github.com/bmf-san/gobel-api/app/domain"

// An AdminRepository is a repository interface for an authentication.
type AdminRepository interface {
	FindByCredential(req RequestCredential) (domain.Admin, error)
	FindIDByToken(token string) (int, error)
	SaveSessionByID(id int) (string, error)
}
