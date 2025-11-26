package repository

import (
	"github.com/nerhays/prestasi_uas/app/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error)
	GetPermissionsByUserID(userID string) ([]model.Permission, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user model.User
	if err := r.db.
		Preload("Role").
		Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetPermissionsByUserID(userID string) ([]model.Permission, error) {
	var perms []model.Permission

	// join users ke roles ke role_permissions ke permissions
	err := r.db.Table("permissions p").
		Select("p.id, p.name, p.resource, p.action, p.description").
		Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN roles r ON r.id = rp.role_id").
		Joins("JOIN users u ON u.role_id = r.id").
		Where("u.id = ?", userID).
		Scan(&perms).Error

	if err != nil {
		return nil, err
	}
	return perms, nil
}
