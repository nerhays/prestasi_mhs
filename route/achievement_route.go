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
func (h *AchievementHandler) Submit(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	refID := c.Param("id")

	ref, err := h.svc.SubmitAchievement(context.Background(), userID, refID)
	if err != nil {
		switch err {
		case service.ErrStudentProfileNotFound, service.ErrRefNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case service.ErrInvalidStatus, service.ErrNotOwner:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": ref})
}

type rejectRequest struct {
	Note string `json:"note" binding:"required"`
}

func (h *AchievementHandler) Verify(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	refID := c.Param("id")

	ref, err := h.svc.VerifyAchievement(context.Background(), userID, refID)
	if err != nil {
		switch err {
		case service.ErrRefNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case service.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": ref})
}

func (h *AchievementHandler) Reject(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	refID := c.Param("id")

	var req rejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	ref, err := h.svc.RejectAchievement(context.Background(), userID, refID, req.Note)
	if err != nil {
		switch err {
		case service.ErrRefNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case service.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": ref})
}
func (h *AchievementHandler) Delete(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	refID := c.Param("id")

	err := h.svc.DeleteDraftAchievement(context.Background(), userID, refID)
	if err != nil {
		switch err {
		case service.ErrRefNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case service.ErrInvalidStatus:
			c.JSON(http.StatusBadRequest, gin.H{"message": "only draft can be deleted"})
		case service.ErrNotOwner:
			c.JSON(http.StatusForbidden, gin.H{"message": "unauthorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "achievement deleted"})
}

func SetupAchievementRoutes(rg *gin.RouterGroup, db *gorm.DB, mongoDB *mongo.Database) {
	achievementRepo := repository.NewAchievementRepository(mongoDB)
	studentRepo := repository.NewStudentRepository(db)
	refRepo := repository.NewAchievementReferenceRepository(db)
	achievementSvc := service.NewAchievementService(achievementRepo, studentRepo, refRepo)
	handler := NewAchievementHandler(achievementSvc)

	ach := rg.Group("/achievements")
	ach.Use(middleware.AuthMiddleware())

	// mahasiswa
	ach.POST("/", handler.Create)
	ach.GET("/me", handler.GetMyAchievements)
	ach.POST("/:id/submit", handler.Submit)
	ach.DELETE("/:id", handler.Delete)

	// dosen wali / admin
	dosen := ach.Group("")
	dosen.Use(middleware.RequireRole("Dosen Wali", "Admin"))
	dosen.POST("/:id/verify", handler.Verify)
	dosen.POST("/:id/reject", handler.Reject)
}
