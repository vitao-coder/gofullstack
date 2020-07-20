package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/vitao-coder/gofullstack/api/models"
)

var usuarios = []models.Usuario{
	models.Usuario{
		Apelido: "Vitão",
		Email:   "vitor51@gmail.com",
		Senha:   "senha",
	},
	models.Usuario{
		Apelido: "Teste",
		Email:   "teste@gmail.com",
		Senha:   "senhateste",
	},
}

var posts = []models.Post{
	models.Post{
		Titulo:   "Mussum Ipsum",
		Conteudo: "Mussum Ipsum, cacilds vidis litro abertis. Manduma pindureta quium dia nois paga.",
	},
	models.Post{
		Titulo:   "Cacilds vidis",
		Conteudo: "Mussum Ipsum, cacilds vidis litro abertis. Manduma pindureta quium dia nois paga.",
	},
}

//CarregarDB : Método para executar as migrações de banco para efetuar login
func CarregarDB(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.Usuario{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Usuario{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("id_do_autor", "usuarios(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range usuarios {
		err = db.Debug().Model(&models.Usuario{}).Create(&usuarios[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].IDDoAutor = usuarios[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
