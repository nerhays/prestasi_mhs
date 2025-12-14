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
	// === repositories ===
	studentRepo := repository.NewStudentRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	refRepo := repository.NewAchievementReferenceRepository(db)
	logRepo := repository.NewAchievementStatusLogRepository(db)
	achievementRepo := repository.NewAchievementRepository(mongoDB)

	// === services ===
	studentSvc := service.NewStudentService(studentRepo, lecturerRepo)
	userSvc := service.NewUserService(userRepo, roleRepo)
	lecturerSvc := service.NewLecturerService(lecturerRepo, studentRepo)
	achievementSvc := service.NewAchievementService(
		achievementRepo,
		studentRepo,
		refRepo,
		userRepo,
		lecturerRepo,
		logRepo,
	)

	// === handlers ===
	studentHandler := NewAdminStudentHandler(studentSvc)
	studentQueryHandler := NewAdminStudentQueryHandler(studentSvc, achievementSvc)
	lecturerHandler := NewAdminLecturerHandler(lecturerSvc)
	userHandler := NewAdminUserHandler(userSvc)
	achievementHandler := NewAdminAchievementHandler(achievementSvc)

	

	admin := rg.Group("/admin")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.RoleOnly("Admin"),
	)

	// === USERS ===
	admin.GET("/users", userHandler.GetAll)
	admin.GET("/users/:id", userHandler.GetByID)
	admin.POST("/users", userHandler.Create)
	admin.PUT("/users/:id", userHandler.Update)
	admin.DELETE("/users/:id", userHandler.Delete)
	admin.PUT("/users/:id/role", userHandler.UpdateRole)

	// === STUDENTS ===
	admin.PUT("/students/:id/advisor", studentHandler.SetAdvisor)
	admin.GET("/students", studentQueryHandler.GetAll)
	admin.GET("/students/:id", studentQueryHandler.GetByID)
	admin.GET("/students/:id/achievements", studentQueryHandler.GetAchievements)
	admin.GET("/reports/student/:id", achievementHandler.GetStudentReport)

	// === ACHIEVEMENTS ===
	admin.GET("/achievements", achievementHandler.GetAllAchievements)

	// === REPORTS ===
	admin.GET("/reports/statistics", achievementHandler.GetStatistics)
	admin.GET("/lecturers", lecturerHandler.GetAll)
	admin.GET("/lecturers/:id/advisees", lecturerHandler.GetAdvisees)
}


