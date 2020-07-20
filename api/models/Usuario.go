package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Usuario : Entidade Usuario do sistema
type Usuario struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Apelido      string    `gorm:"size:255;not null;unique" json:"apelido"`
	Email        string    `gorm:"size:100;not null;unique" json:"email"`
	Senha        string    `gorm:"size:100;not null;" json:"senha"`
	CriadoEm     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"criadoEm"`
	AtualizadoEm time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"atualizadoEm"`
}

//Encriptar : Encriptação da senha em hash
func Encriptar(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

//VerificarSenha : Verificação da senha com o hash
func VerificarSenha(senhaCriptografada, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaCriptografada), []byte(senha))
}

//AntesDeSalvar : Executa ações importantes antes de salvar
func (u *Usuario) AntesDeSalvar() error {
	senhaCriptografada, err := Encriptar(u.Senha)
	if err != nil {
		return err
	}
	u.Senha = string(senhaCriptografada)
	return nil
}

//Preparar : Prepara o usuário para ser salvo
func (u *Usuario) Preparar() {
	u.ID = 0
	u.Apelido = html.EscapeString(strings.TrimSpace(u.Apelido))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CriadoEm = time.Now()
	u.AtualizadoEm = time.Now()
}

//Validar : Valida todos campos de usuário antes de salvar de acordo com a ação
func (u *Usuario) Validar(acao string) error {
	switch strings.ToLower(acao) {
	case "update":
		if u.Apelido == "" {
			return errors.New("Apelido obrigatório")
		}
		if u.Senha == "" {
			return errors.New("Senha obrigatória")
		}
		if u.Email == "" {
			return errors.New("Email obrigatório")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email inválido")
		}

		return nil
	case "login":
		if u.Senha == "" {
			return errors.New("Senha obrigatória")
		}
		if u.Email == "" {
			return errors.New("Email obrigatório")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email inválido")
		}
		return nil

	default:
		if u.Apelido == "" {
			return errors.New("Apelido obrigatório")
		}
		if u.Senha == "" {
			return errors.New("Senha obrigatória")
		}
		if u.Email == "" {
			return errors.New("Email obrigatório")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email inválido")
		}
		return nil
	}
}

//SalvarUsuario : Salva usuário no banco de dados
func (u *Usuario) SalvarUsuario(db *gorm.DB) (*Usuario, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	return u, nil
}

//BuscarTodosUsuarios : Busca todos usuário da base de dados
func (u *Usuario) BuscarTodosUsuarios(db *gorm.DB) (*[]Usuario, error) {
	var err error
	usuarios := []Usuario{}
	err = db.Debug().Model(&Usuario{}).Limit(100).Find(&usuarios).Error
	if err != nil {
		return &[]Usuario{}, err
	}
	return &usuarios, err
}

//BuscarUsuarioID : Busca usuário da base de dados pelo seu ID
func (u *Usuario) BuscarUsuarioID(db *gorm.DB, uid uint32) (*Usuario, error) {
	var err error
	err = db.Debug().Model(Usuario{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Usuario{}, errors.New("Usuário não encontrado")
	}
	return u, err
}

//AlterarUsuario : Altera usuário na base de dados
func (u *Usuario) AlterarUsuario(db *gorm.DB, uid uint32) (*Usuario, error) {

	// To hash the password
	err := u.AntesDeSalvar()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&Usuario{}).UpdateColumns(
		map[string]interface{}{
			"senha":        u.Senha,
			"apelido":      u.Apelido,
			"email":        u.Email,
			"atualizadoEm": time.Now(),
		},
	)
	if db.Error != nil {
		return &Usuario{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Usuario{}, err
	}
	return u, nil
}

//ExcluirUsuario : Exclui usuário da base de dados
func (u *Usuario) ExcluirUsuario(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Usuario{}).Where("id = ?", uid).Take(&Usuario{}).Delete(&Usuario{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
