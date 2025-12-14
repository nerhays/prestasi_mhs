package route

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nerhays/prestasi_uas/app/service"
)

type AdminStudentHandler struct {
	studentSvc *service.StudentService
}

func NewAdminStudentHandler(studentSvc *service.StudentService) *AdminStudentHandler {
	return &AdminStudentHandler{studentSvc}
}

type SetAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" binding:"required"`
}

func (h *AdminStudentHandler) SetAdvisor(c *gin.Context) {
	studentID := c.Param("id")

	var req SetAdvisorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	if err := h.studentSvc.SetAdvisor(
		c.Request.Context(),
		studentID,
		req.AdvisorID,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "advisor assigned successfully",
	})
}
type AdminAchievementHandler struct {
	achievementSvc *service.AchievementService
}

func NewAdminAchievementHandler(
	achievementSvc *service.AchievementService,
) *AdminAchievementHandler {
	return &AdminAchievementHandler{achievementSvc}
}

func (h *AdminAchievementHandler) GetAllAchievements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	data, total, err := h.achievementSvc.GetAllAchievements(
		c.Request.Context(),
		page,
		limit,
		statusPtr,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *AdminAchievementHandler) GetStatistics(c *gin.Context) {
	stats, err := h.achievementSvc.GetStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}


