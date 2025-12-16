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

// SetAdvisor godoc
// @Summary Assign advisor to student
// @Description Assign dosen wali to a student (Admin only)
// @Tags Admin - Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body SetAdvisorRequest true "Advisor payload"
// @Success 200 {object} map[string]string "Advisor assigned"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/students/{id}/advisor [put]
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

// GetAllAchievements godoc
// @Summary Get all achievements
// @Description Admin can view all achievements with pagination and filter
// @Tags Admin - Achievements
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param status query string false "Achievement status"
// @Success 200 {object} map[string]interface{} "List of achievements"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/achievements [get]
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

// GetStatistics godoc
// @Summary Get achievement statistics
// @Description Get statistics of achievements by type and status
// @Tags Admin - Reports
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Statistics data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/reports/statistics [get]
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

type AdminStudentQueryHandler struct {
	studentSvc *service.StudentService
	achievementSvc *service.AchievementService
}

func NewAdminStudentQueryHandler(
	studentSvc *service.StudentService,
	achievementsvc *service.AchievementService,
) *AdminStudentQueryHandler {
	return &AdminStudentQueryHandler{studentSvc, achievementsvc}
}

// GetAllStudents godoc
// @Summary Get all students
// @Description Admin can view all students
// @Tags Admin - Students
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of students"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/students [get]
func (h *AdminStudentQueryHandler) GetAll(c *gin.Context) {
	data, err := h.studentSvc.GetAllStudents()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": data})
}

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Get detail student by ID
// @Tags Admin - Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{} "Student detail"
// @Failure 404 {object} map[string]string "Student not found"
// @Router /admin/students/{id} [get]
func (h *AdminStudentQueryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	data, err := h.studentSvc.GetStudentByID(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "student not found"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": data})
}

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Get all achievements for a student
// @Tags Admin - Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{} "Student achievements"
// @Failure 404 {object} map[string]string "Student not found"
// @Router /admin/students/{id}/achievements [get]
func (h *AdminStudentQueryHandler) GetAchievements(c *gin.Context) {
	studentID := c.Param("id")

	data, err := h.achievementSvc.GetAchievementsByStudentID(
		c.Request.Context(),
		studentID,
	)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": data})
}

type AdminLecturerHandler struct {
	lecturerSvc *service.LecturerService
}

func NewAdminLecturerHandler(
	lecturerSvc *service.LecturerService,
) *AdminLecturerHandler {
	return &AdminLecturerHandler{lecturerSvc}
}

// GetAllLecturers godoc
// @Summary Get all lecturers
// @Description Retrieve list of all lecturers
// @Tags Admin - Lecturers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of lecturers"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/lecturers [get]
func (h *AdminLecturerHandler) GetAll(c *gin.Context) {
	data, err := h.lecturerSvc.GetAllLecturers()
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": data})
}

// GetAdvisees godoc
// @Summary Get lecturer advisees
// @Description Get students supervised by a lecturer
// @Tags Admin - Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{} "List of advisees"
// @Failure 404 {object} map[string]string "Lecturer not found"
// @Router /admin/lecturers/{id}/advisees [get]
func (h *AdminLecturerHandler) GetAdvisees(c *gin.Context) {
	id := c.Param("id")
	data, err := h.lecturerSvc.GetAdvisees(id)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": data})
}

// GetStudentReport godoc
// @Summary Get student achievement report
// @Description Get achievement summary report for a student
// @Tags Admin - Reports
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{} "Student report"
// @Failure 404 {object} map[string]string "Student not found"
// @Router /admin/reports/student/{id} [get]
func (h *AdminAchievementHandler) GetStudentReport(c *gin.Context) {
	studentID := c.Param("id")

	data, err := h.achievementSvc.GetStudentReport(
		c.Request.Context(),
		studentID,
	)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}
