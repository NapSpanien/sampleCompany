package middlewares

import (
	"net/http"

	"gitlab.com/napspan/SampleCompany/api/models"
)

// Middleware Middleware
type Middleware func(http.HandlerFunc) http.HandlerFunc

// ValidateContentType ValidateContentType
func ValidateContentType() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			contentType := r.Header.Get("content-type")
			if contentType != "application/json" {
				mR := models.MyResponse{}
				mR.Code = 1
				mR.Msg = "Invalid Content-Type"
				models.GenerateResponse(w, mR, http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// ValidateAuthorizationApp ValidateAuthorizationApp
func ValidateAuthorization() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			authorization := r.Header.Get("Authorization")
			if authorization == "" {
				mR := models.MyResponse{}
				mR.Code = 1
				mR.Msg = "Authorization is missing!"
				models.GenerateResponse(w, mR, http.StatusBadRequest)
				return
			}
			if authorization != "auth test" {
				mR := models.MyResponse{}
				mR.Code = 1
				mR.Msg = "Authorization is invalid!"
				models.GenerateResponse(w, mR, http.StatusBadRequest)
				return
			}
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
