package route

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// CreateAchievement godoc
// @Summary Create achievement (draft)
// @Description Mahasiswa membuat prestasi baru dengan status draft
// @Tags Achievements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.Achievement true "Achievement payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements [post]
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

// GetMyAchievements godoc
// @Summary Get my achievements
// @Description Mahasiswa melihat semua prestasi miliknya
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Achievement
// @Failure 401 {object} map[string]string
// @Router /achievements/me [get]
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

// SubmitAchievement godoc
// @Summary Submit achievement for verification
// @Description Mahasiswa submit prestasi draft untuk diverifikasi dosen wali
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} model.AchievementReference
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id}/submit [post]
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

// VerifyAchievement godoc
// @Summary Verify achievement
// @Description Dosen wali atau admin memverifikasi prestasi mahasiswa
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} model.AchievementReference
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id}/verify [post]
func (h *AchievementHandler) Verify(c *gin.Context) {
    userID := c.GetString(middleware.ContextUserIDKey)
    refID := c.Param("id")

    ref, err := h.svc.VerifyAchievement(context.Background(), userID, refID)
    if err != nil {
        log.Println("VerifyAchievement error:", err) // tetap log untuk debug

        switch {
        case errors.Is(err, service.ErrRefNotFound):
            c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        case errors.Is(err, service.ErrInvalidStatus):
            c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        case errors.Is(err, service.ErrNotAdvisor):
            c.JSON(http.StatusForbidden, gin.H{"message": "only the assigned advisor can verify this achievement"})
        case errors.Is(err, service.ErrUserNotFound):
            c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"status":"success","data": ref})
}

// RejectAchievement godoc
// @Summary Reject achievement
// @Description Dosen wali menolak prestasi dengan catatan
// @Tags Achievements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Param body body rejectRequest true "Rejection note"
// @Success 200 {object} model.AchievementReference
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id}/reject [post]
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
		switch {
        case errors.Is(err, service.ErrRefNotFound):
            c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        case errors.Is(err, service.ErrInvalidStatus):
            c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        case errors.Is(err, service.ErrNotAdvisor):
            c.JSON(http.StatusForbidden, gin.H{"message": "only the assigned advisor can verify this achievement"})
        case errors.Is(err, service.ErrUserNotFound):
            c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
        }
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": ref})
}

// DeleteAchievement godoc
// @Summary Delete draft achievement
// @Description Mahasiswa menghapus prestasi berstatus draft
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [delete]
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

// GetDeletedAchievements godoc
// @Summary Get deleted achievements
// @Description Mahasiswa melihat prestasi yang dihapus
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Achievement
// @Router /achievements/deleted [get]
func (h *AchievementHandler) GetDeleted(c *gin.Context) {
    userID := c.GetString(middleware.ContextUserIDKey)

    res, err := h.svc.GetDeletedAchievements(context.Background(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   res,
    })
}

// GetBimbinganAchievements godoc
// @Summary Get achievements under supervision
// @Description Dosen wali melihat prestasi mahasiswa bimbingannya
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param status query string false "Achievement status"
// @Success 200 {object} map[string]interface{}
// @Router /achievements/bimbingan [get]
func (h *AchievementHandler) GetBimbingan(c *gin.Context) {
    userID := c.GetString(middleware.ContextUserIDKey)

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
    var status *model.AchievementStatus
    if s := c.Query("status"); s != "" {
        st := model.AchievementStatus(s)
        status = &st
    }

    total, rows, err := h.svc.GetBimbinganAchievements(c.Request.Context(), userID, page, perPage, status)
    if err != nil {
        log.Println("GetBimbingan error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "meta": gin.H{"page": page, "per_page": perPage, "total": total},
        "data": rows,
    })
}

// UploadAttachment godoc
// @Summary Upload achievement attachment
// @Description Mahasiswa upload file bukti prestasi
// @Tags Achievements
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Param file formData file true "Attachment file"
// @Success 200 {object} model.Attachment
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id}/attachments [post]
func (h *AchievementHandler) UploadAttachment(c *gin.Context) {
	// 🔑 ambil user dari context (WAJIB pakai constant)
	userID := c.GetString(middleware.ContextUserIDKey)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	// refID = achievement_reference.id (Postgres UUID)
	refID := c.Param("id")

	// ambil file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file required"})
		return
	}

	// validasi ekstensi (simple & aman)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid file type"})
		return
	}

	// buat folder jika belum ada
	uploadDir := "uploads/achievements"
	_ = os.MkdirAll(uploadDir, 0755)

	// generate nama file aman
	filename := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, filename)

	// simpan file ke disk
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// panggil service
	attachment, err := h.svc.UploadAttachment(
		c.Request.Context(),
		userID,
		refID,
		filename,
		"/"+filePath,
		file.Header.Get("Content-Type"),
	)
	if err != nil {
		_ = os.Remove(filePath) // rollback file
		c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   attachment,
	})
}

// GetAchievementHistory godoc
// @Summary Get achievement history
// @Description Melihat riwayat perubahan status prestasi
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {array} model.AchievementStatusLog
// @Router /achievements/{id}/history [get]
func (h *AchievementHandler) GetHistory(c *gin.Context) {
	refID := c.Param("id")

	logs, err := h.svc.GetAchievementHistory(c.Request.Context(), refID)
	if err != nil {
		c.JSON(404, gin.H{"message": "history not found"})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   logs,
	})
}

// GetAchievementDetail godoc
// @Summary Get achievement detail
// @Description Detail prestasi (RBAC: Mahasiswa, Dosen Wali, Admin)
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [get]
func (h *AchievementHandler) GetDetail(c *gin.Context) {
	refID := c.Param("id")

	userID := c.GetString(middleware.ContextUserIDKey)
	role := c.GetString(middleware.ContextRoleKey)

	data, err := h.svc.GetAchievementDetail(
		c.Request.Context(),
		refID,
		userID,
		role,
	)
	if err != nil {
		c.JSON(403, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}

// UpdateAchievementDraft godoc
// @Summary Update draft achievement
// @Description Mahasiswa mengubah prestasi berstatus draft
// @Tags Achievements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Achievement Reference ID"
// @Param body body model.Achievement true "Updated payload"
// @Success 200 {object} model.Achievement
// @Failure 400 {object} map[string]string
// @Router /achievements/{id} [put]
func (h *AchievementHandler) Update(c *gin.Context) {
	refID := c.Param("id")
	userID := c.GetString(middleware.ContextUserIDKey)

	var req model.Achievement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "invalid input"})
		return
	}

	data, err := h.svc.UpdateAchievementDraft(
		c.Request.Context(),
		refID,
		userID,
		&req,
	)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": data})
}

// GetAchievementsByRole godoc
// @Summary Get achievements by role
// @Description Mengambil prestasi sesuai role user (Mahasiswa, Dosen Wali, Admin)
// @Tags Achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Achievement
// @Router /achievements [get]
func (h *AchievementHandler) GetListByRole(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)
	role := c.GetString(middleware.ContextRoleKey)

	data, err := h.svc.GetAchievementsByRole(
		c.Request.Context(),
		userID,
		role,
	)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": data})
}


func SetupAchievementRoutes(rg *gin.RouterGroup, db *gorm.DB, mongoDB *mongo.Database) {
	achievementRepo := repository.NewAchievementRepository(mongoDB)
	studentRepo := repository.NewStudentRepository(db)
	refRepo := repository.NewAchievementReferenceRepository(db)
	userRepo := repository.NewUserRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)
	logRepo := repository.NewAchievementStatusLogRepository(db)
	achievementSvc := service.NewAchievementService(achievementRepo, studentRepo, refRepo, userRepo, lecturerRepo, logRepo)
	handler := NewAchievementHandler(achievementSvc)

	ach := rg.Group("/achievements")
	ach.Use(middleware.AuthMiddleware())

	// mahasiswa
	ach.POST("/", handler.Create)
	ach.GET("/me", handler.GetMyAchievements)
	ach.POST("/:id/submit", handler.Submit)
	ach.DELETE("/:id", handler.Delete)
	ach.GET("/deleted", handler.GetDeleted)
	ach.POST("/:id/attachments", handler.UploadAttachment)
	ach.GET("/:id/history", handler.GetHistory)
	ach.GET("/:id", handler.GetDetail)
	ach.PUT("/:id", handler.Update)
	ach.GET("/", handler.GetListByRole)




	// dosen wali / admin
	dosen := ach.Group("")
	dosen.Use(middleware.RequireRole("Dosen Wali", "Admin"))
	dosen.POST("/:id/verify", handler.Verify)
	dosen.POST("/:id/reject", handler.Reject)
	dosen.GET("/bimbingan", handler.GetBimbingan)
}
