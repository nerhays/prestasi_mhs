package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/nerhays/prestasi_uas/app/repository"
	"github.com/nerhays/prestasi_uas/app/service"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   roles,
	})
}

func SetupRoleRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	roleRepo := repository.NewRoleRepository(db)
	roleSvc := service.NewRoleService(roleRepo)
	handler := NewRoleHandler(roleSvc)

	rg.GET("/roles", handler.GetAll)
}
