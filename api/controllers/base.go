package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver

	"github.com/vitao-coder/gofullstack/api/models"
)

//Servidor : Entidade Servidor que guarda as informações bases de rotas e banco de dados.
type Servidor struct {
	DB     *gorm.DB
	Router *mux.Router
}

//Inicializar : Inicializa aplicação conectando ao banco de dados e servindo as rotas
func (servidor *Servidor) Inicializar(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		servidor.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Não foi possível conectar ao banco de dados: %s", Dbdriver)
			log.Fatal("Erro:", err)
		} else {
			fmt.Printf("Conectado ao banco de dados: %s", Dbdriver)
		}
	}
	servidor.DB.Debug().AutoMigrate(&models.Usuario{}, &models.Post{}) //database migration
	servidor.Router = mux.NewRouter()
	servidor.inicializarRotas()
}

//Rodar : Executa o servidor na porta 80 e escuta as requisições
func (servidor *Servidor) Rodar(addr string) {
	fmt.Println("Escutando na porta 8080")
	log.Fatal(http.ListenAndServe(addr, servidor.Router))
}
