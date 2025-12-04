package route

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nerhays/prestasi_uas/app/model"
	"github.com/nerhays/prestasi_uas/app/repository"
	"github.com/nerhays/prestasi_uas/app/service"
	"github.com/nerhays/prestasi_uas/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type AchievementHandler struct {
	svc *service.AchievementService
}

func NewAchievementHandler(svc *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{svc: svc}
}

func (h *AchievementHandler) Create(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)

	var req model.Achievement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	ac, ref, err := h.svc.CreateAchievementForUser(context.Background(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"achievement": ac,
			"reference":   ref,
		},
	})
}
func (h *AchievementHandler) GetMyAchievements(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)

	acs, err := h.svc.GetMyAchievements(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   acs,
	})
}
func SetupAchievementRoutes(rg *gin.RouterGroup, db *gorm.DB, mongoDB *mongo.Database) {
	achievementRepo := repository.NewAchievementRepository(mongoDB)
	studentRepo := repository.NewStudentRepository(db)
	refRepo := repository.NewAchievementReferenceRepository(db)
	achievementSvc := service.NewAchievementService(achievementRepo, studentRepo, refRepo)
	handler := NewAchievementHandler(achievementSvc)

	authRequired := rg.Group("/achievements")
	authRequired.Use(middleware.AuthMiddleware())
	authRequired.POST("/", handler.Create)
	authRequired.GET("/me", handler.GetMyAchievements)
}
