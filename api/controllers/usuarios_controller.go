package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vitao-coder/gofullstack/api/auth"
	"github.com/vitao-coder/gofullstack/api/models"
	"github.com/vitao-coder/gofullstack/api/responses"
	"github.com/vitao-coder/gofullstack/api/utils/formaterror"
)

//CriarUsuario : Rota para criação de usuários
func (servidor *Servidor) CriarUsuario(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	usuario := models.Usuario{}
	err = json.Unmarshal(body, &usuario)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	usuario.Preparar()
	err = usuario.Validar("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	usuarioCriado, err := usuario.SalvarUsuario(servidor.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, usuarioCriado.ID))
	responses.JSON(w, http.StatusCreated, usuarioCriado)
}

//GetUsuarios : Rota para buscar todos usuários
func (servidor *Servidor) GetUsuarios(w http.ResponseWriter, r *http.Request) {

	usuario := models.Usuario{}

	usuarios, err := usuario.BuscarTodosUsuarios(servidor.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, usuarios)
}

//GetUsuario : Rota para buscar um usuário por ID
func (servidor *Servidor) GetUsuario(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	usuario := models.Usuario{}
	usuarioRetornado, err := usuario.BuscarUsuarioID(servidor.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, usuarioRetornado)
}

//AlterarUsuario : Rota para alterar um usuário
func (servidor *Servidor) AlterarUsuario(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
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
	tokenID, err := auth.ExtrairIDDoToken(r)
	if err != nil || tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	usuario.Preparar()
	err = usuario.Validar("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	usuarioAlterado, err := usuario.AlterarUsuario(servidor.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, usuarioAlterado)
}

//DeletarUsuario : Rota para alterar um usuário
func (servidor *Servidor) DeletarUsuario(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	usuario := models.Usuario{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtrairIDDoToken(r)
	if err != nil || tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = usuario.ExcluirUsuario(servidor.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}
