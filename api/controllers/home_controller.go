package controllers

import (
	"net/http"

	"github.com/vitao-coder/gofullstack/api/responses"
)

//Home : Home controller
func (servidor *Servidor) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Usuários e Posts - GOLang API")
}
