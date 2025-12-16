package main

import (
	"log"

	"github.com/nerhays/prestasi_uas/config"
	"github.com/nerhays/prestasi_uas/database"
	"github.com/nerhays/prestasi_uas/route"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/nerhays/prestasi_uas/docs"
)

// @title Prestasi Mahasiswa API
// @version 1.0
// @description Sistem Manajemen Prestasi Mahasiswa
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@prestasi.ac.id

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1
// @schemes http

// ===== SECURITY =====
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	pgDB := database.NewPostgres(cfg.PostgresDSN)
	mongo := database.NewMongo(cfg.MongoURI, cfg.MongoDB)
	
	r := route.SetupRouter(pgDB, mongo.DB)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Printf("[APP] Server running on :%s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
