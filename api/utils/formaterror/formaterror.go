package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "apelido") {
		return errors.New("Apelido já está em uso")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email já está em uso")
	}

	if strings.Contains(err, "titulo") {
		return errors.New("Titulo já está em uso")
	}
	if strings.Contains(err, "senha") {
		return errors.New("Senha inválida")
	}
	return errors.New("Detalhes incorretos")
}
