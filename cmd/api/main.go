package main

import (
	"log"

	"github.com/rubewafula/edairy-go-26/internal/db"
	initializers "github.com/rubewafula/edairy-go-26/internal/initializers"
	"github.com/rubewafula/edairy-go-26/internal/routes"
)

func init() {
	initializers.LoadEnvVariables()
	db.ConnectToDatabase()
}

func main() {

	r := routes.SetupRouter()

	addr := initializers.GetEnv("PORT", "8000")

	log.Printf("listening on :%s", addr)
	if err := r.Run("0.0.0.0:" + addr); err != nil {
		log.Fatal(err)
	}
}
