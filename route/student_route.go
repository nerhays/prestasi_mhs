package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nerhays/prestasi_uas/app/repository"
	"github.com/nerhays/prestasi_uas/app/service"
	"github.com/nerhays/prestasi_uas/middleware"
	"gorm.io/gorm"
)

type StudentHandler struct {
	studentService *service.StudentService
}

func NewStudentHandler(studentService *service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

func (h *StudentHandler) GetMyProfile(c *gin.Context) {
	userID := c.GetString("userID") // dari JWT middleware

	student, err := h.studentService.GetProfileByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   student,
	})
}

func SetupStudentRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	studentRepo := repository.NewStudentRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)

	studentSvc := service.NewStudentService(
		studentRepo,
		lecturerRepo,
	)

	handler := NewStudentHandler(studentSvc)

	authRequired := rg.Group("/students", middleware.AuthMiddleware())
	authRequired.GET("/me", handler.GetMyProfile)
}
