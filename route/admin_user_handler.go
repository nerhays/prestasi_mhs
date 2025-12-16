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

// GetAllUsers godoc
// @Summary Get all users
// @Description Admin can retrieve all users
// @Tags Admin - Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "List of users"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/users [get]
func (h *AdminUserHandler) GetAll(c *gin.Context) {
	users, _ := h.userSvc.GetAllUsers()
	c.JSON(200, gin.H{"status": "success", "data": users})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Admin can retrieve user detail by ID
// @Tags Admin - Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User detail"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/users/{id} [get]
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

// CreateUser godoc
// @Summary Create new user
// @Description Admin creates a new user
// @Tags Admin - Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "Create user payload"
// @Success 201 {object} map[string]string "User created"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/users [post]
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

// UpdateUser godoc
// @Summary Update user
// @Description Admin updates user information
// @Tags Admin - Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body UpdateUserRequest true "Update user payload"
// @Success 200 {object} map[string]string "User updated"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/users/{id} [put]
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

// DeleteUser godoc
// @Summary Delete user
// @Description Admin deletes a user
// @Tags Admin - Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "User deleted"
// @Failure 400 {object} map[string]string "Delete failed"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/users/{id} [delete]
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

// UpdateUserRole godoc
// @Summary Update user role
// @Description Admin updates role of a user
// @Tags Admin - Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body UpdateRoleRequest true "Update role payload"
// @Success 200 {object} map[string]string "Role updated"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/users/{id}/role [put]
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
