package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//Post : Entidade Post do sistema
type Post struct {
	ID           uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Titulo       string    `gorm:"size:255;not null;unique" json:"titulo"`
	Conteudo     string    `gorm:"size:255;not null;" json:"conteudo"`
	Autor        Usuario   `json:"autor"`
	IDDoAutor    uint32    `gorm:"not null" json:"id_do_autor"`
	CriadoEm     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"criadoEm"`
	AtualizadoEm time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"atualizadoEm"`
}

//Preparar : Prepara o Post para ser salvo
func (p *Post) Preparar() {
	p.ID = 0
	p.Titulo = html.EscapeString(strings.TrimSpace(p.Titulo))
	p.Conteudo = html.EscapeString(strings.TrimSpace(p.Conteudo))
	p.Autor = Usuario{}
	p.CriadoEm = time.Now()
	p.AtualizadoEm = time.Now()
}

//Validar : Valida todos campos de Post antes de salvar
func (p *Post) Validar() error {

	if p.Titulo == "" {
		return errors.New("Título obrigatório")
	}
	if p.Conteudo == "" {
		return errors.New("Conteúdo obrigatório")
	}
	if p.IDDoAutor < 1 {
		return errors.New("Autor obrigatório")
	}
	return nil
}

//SalvarPost : Salva Post no banco de dados
func (p *Post) SalvarPost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Usuario{}).Where("id = ?", p.IDDoAutor).Take(&p.Autor).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

//BuscarTodosPosts : Busca todos os Posts no banco de dados
func (p *Post) BuscarTodosPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&Usuario{}).Where("id = ?", posts[i].IDDoAutor).Take(&posts[i].Autor).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

//BuscarPostPorID : Busca Post por ID no banco de dados
func (p *Post) BuscarPostPorID(db *gorm.DB, pid uint64) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Usuario{}).Where("id = ?", p.IDDoAutor).Take(&p.Autor).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

//AtualizarPost : Altera Post na base de dadosdos
func (p *Post) AtualizarPost(db *gorm.DB) (*Post, error) {

	var err error

	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Titulo: p.Titulo, Conteudo: p.Conteudo, AtualizadoEm: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Usuario{}).Where("id = ?", p.IDDoAutor).Take(&p.Autor).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

//ExcluirPost : Deleta Post na base de dados
func (p *Post) ExcluirPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Post{}).Where("id = ? and id_do_autor = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post não encontrado")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
