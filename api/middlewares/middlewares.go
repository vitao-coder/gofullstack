package middlewares

import (
	"errors"
	"net/http"

	"github.com/vitao-coder/gofullstack/api/auth"
	"github.com/vitao-coder/gofullstack/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthenticacao(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValido(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Não autorizado"))
			return
		}
		next(w, r)
	}
}
