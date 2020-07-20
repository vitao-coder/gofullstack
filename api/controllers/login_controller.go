package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vitao-coder/gofullstack/api/auth"

	"github.com/vitao-coder/gofullstack/api/models"
	"github.com/vitao-coder/gofullstack/api/responses"
	"github.com/vitao-coder/gofullstack/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

//Login : Rota para efetuar login
func (servidor *Servidor) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	usuario := models.Usuario{}
	err = json.Unmarshal(body, &usuario)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	usuario.Preparar()
	err = usuario.Validar("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := servidor.Autenticar(usuario.Email, usuario.Senha)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

//Autenticar : Método para autenticar o usuário
func (servidor *Servidor) Autenticar(email, senha string) (string, error) {

	var err error

	usuario := models.Usuario{}

	err = servidor.DB.Debug().Model(models.Usuario{}).Where("email = ?", email).Take(&usuario).Error
	if err != nil {
		return "", err
	}
	err = models.VerificarSenha(usuario.Senha, senha)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CriarToken(usuario.ID)
}
