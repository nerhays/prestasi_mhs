package route

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/nerhays/prestasi_uas/app/repository"
	"github.com/nerhays/prestasi_uas/app/service"
	"github.com/nerhays/prestasi_uas/middleware"
)

/*
ADMIN ROUTES
Base path: /api/v1/admin
Middleware: Auth + RoleOnly(Admin)
*/

func SetupAdminRoutes(
	rg *gin.RouterGroup,
	db *gorm.DB,
	mongoDB *mongo.Database,
) {
	// repositories
	studentRepo := repository.NewStudentRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)
	userRepo := repository.NewUserRepository(db)
	refRepo := repository.NewAchievementReferenceRepository(db)
	logRepo := repository.NewAchievementStatusLogRepository(db)
	achievementRepo := repository.NewAchievementRepository(mongoDB)

	// services
	studentSvc := service.NewStudentService(studentRepo, lecturerRepo)
	achievementSvc := service.NewAchievementService(
		achievementRepo,
		studentRepo,
		refRepo,
		userRepo,
		lecturerRepo,
		logRepo,
	)

	// handlers
	studentHandler := NewAdminStudentHandler(studentSvc)
	achievementHandler := NewAdminAchievementHandler(achievementSvc)

	admin := rg.Group("/admin")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.RoleOnly("Admin"),
	)

	// === STUDENT MANAGEMENT ===
	admin.PUT("/students/:id/advisor", studentHandler.SetAdvisor)

	// === ACHIEVEMENT MANAGEMENT ===
	admin.GET("/achievements", achievementHandler.GetAllAchievements)
	admin.GET("/reports/statistics", achievementHandler.GetStatistics)

}

