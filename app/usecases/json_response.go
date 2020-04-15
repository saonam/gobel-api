package usecases

import "net/http"

// A JSONResponse is a presenter interface for response.
type JSONResponse interface {
	Success200(w http.ResponseWriter, res []byte)
	Success201(w http.ResponseWriter, res []byte)
	Error403(w http.ResponseWriter)
	Error404(w http.ResponseWriter)
	Error500(w http.ResponseWriter)
}
