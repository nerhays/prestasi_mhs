package main

import (
	"log"

	"github.com/nerhays/prestasi_uas/config"
	"github.com/nerhays/prestasi_uas/database"
	"github.com/nerhays/prestasi_uas/route"
)

func main() {
	cfg := config.LoadConfig()

	pgDB := database.NewPostgres(cfg.PostgresDSN)

	database.Seed(pgDB)

	_ = database.NewMongo(cfg.MongoURI, cfg.MongoDB) // dipakai achievements

	r := route.SetupRouter(pgDB)

	log.Printf("[APP] Server running on :%s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
