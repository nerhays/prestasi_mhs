package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/nerhays/prestasi_uas/middleware"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// health check (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Prestasi API is running",
		})
	})

	api := r.Group("/api/v1")

	// PUBLIC ROUTES
	SetupAuthRoutes(api, db) // /auth/login

	// PROTECTED ROUTES (JWT)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	SetupRoleRoutes(protected, db)
	SetupStudentRoutes(protected, db)
	// SetupAchievementRoutes(protected, db, mongo)

	return r
}
