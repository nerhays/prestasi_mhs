package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Prestasi API is running",
		})
	})

	api := r.Group("/api/v1")
	{
		SetupRoleRoutes(api, db)
		SetupAuthRoutes(api, db)
		// SetupAchievementRoutes(api, db, mongo), dst.
	}

	return r
}
