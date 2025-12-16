package route

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/nerhays/prestasi_uas/app/repository"
	"github.com/nerhays/prestasi_uas/app/service"
	"github.com/nerhays/prestasi_uas/middleware"
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

// Login godoc
// @Summary Login user
// @Description Login menggunakan username dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "Login payload"
// @Success 200 {object} map[string]interface{} "Login success"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/login [post]
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

// Refresh godoc
// @Summary Refresh JWT token
// @Description Generate token baru dari token lama
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Token refreshed"
// @Failure 401 {object} map[string]string "Invalid or expired token"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{"message": "missing token"})
		return
	}

	oldToken := strings.TrimPrefix(authHeader, "Bearer ")

	newToken, err := h.authService.RefreshToken(oldToken)
	if err != nil {
		c.JSON(401, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"token": newToken,
		},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user (JWT stateless, client-side logout)
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Logout success"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT stateless → cukup respon sukses
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "logged out successfully",
	})
}

// Profile godoc
// @Summary Get user profile
// @Description Get profile user yang sedang login
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetString(middleware.ContextUserIDKey)

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"fullName": user.FullName,
			"role":     user.Role.Name,
		},
	})
}


func SetupAuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	handler := NewAuthHandler(authSvc)

	auth := rg.Group("/auth")

	// public
	auth.POST("/login", handler.Login)
	auth.POST("/refresh", handler.Refresh)

	// protected
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/logout", handler.Logout)
	auth.GET("/profile", handler.Profile)
}
