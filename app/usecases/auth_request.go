package usecases

// A RequestCredential represents the singular of jwtauth for jwtauth.
type RequestCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
