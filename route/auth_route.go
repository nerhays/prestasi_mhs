package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/nerhays/prestasi_uas/app/repository"
	"github.com/nerhays/prestasi_uas/app/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	res, err := h.authService.Login(service.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		// simple error handling, nanti bisa dibagusin
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// mapping permission ke string
	perms := make([]string, 0, len(res.Permissions))
	for _, p := range res.Permissions {
		perms = append(perms, p.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": res.Token,
			"user": gin.H{
				"id":       res.User.ID,
				"username": res.User.Username,
				"fullName": res.User.FullName,
				"role":     res.User.Role.Name,
			},
			"permissions": perms,
		},
	})
}

func SetupAuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	handler := NewAuthHandler(authSvc)

	rg.POST("/auth/login", handler.Login)
}
