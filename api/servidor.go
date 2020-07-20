package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vitao-coder/gofullstack/api/controllers"
	"github.com/vitao-coder/gofullstack/api/seed"
)

var servidor = controllers.Servidor{}

//Executar : Executa o serviço da api
func Executar() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	servidor.Inicializar(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	seed.CarregarDB(servidor.DB)
	servidor.Rodar(":8080")
}
