package route

import (
	"github.com/gin-gonic/gin"
	"github.com/nerhays/prestasi_uas/app/service"
)
type AdminUserHandler struct {
	userSvc *service.UserService
}

func NewAdminUserHandler(userSvc *service.UserService) *AdminUserHandler {
	return &AdminUserHandler{userSvc}
}
func (h *AdminUserHandler) GetAll(c *gin.Context) {
	users, _ := h.userSvc.GetAllUsers()
	c.JSON(200, gin.H{"status": "success", "data": users})
}
func (h *AdminUserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userSvc.GetUserByID(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "user not found"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": user})
}
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	RoleID   string `json:"role_id"`
}

func (h *AdminUserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	c.ShouldBindJSON(&req)

	err := h.userSvc.CreateUser(
		req.Username,
		req.Email,
		req.Password,
		req.FullName,
		req.RoleID,
	)

	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "success"})
}
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func (h *AdminUserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	c.ShouldBindJSON(&req)

	if err := h.userSvc.UpdateUser(id, req.Username, req.Email, req.FullName); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
func (h *AdminUserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.userSvc.DeleteUser(id); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
type UpdateRoleRequest struct {
	RoleID string `json:"role_id"`
}

func (h *AdminUserHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req UpdateRoleRequest
	c.ShouldBindJSON(&req)

	if err := h.userSvc.UpdateUserRole(id, req.RoleID); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
