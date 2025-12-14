package repository

import (
	"github.com/nerhays/prestasi_uas/app/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error)
	GetPermissionsByUserID(userID string) ([]model.Permission, error)
	FindByID(id string) (*model.User, error)
	FindAll() ([]model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id string) error
	UpdateRole(userID, roleID string) error
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
func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.
		Preload("Role").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Preload("Role").Find(&users).Error
	return users, err
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *userRepository) UpdateRole(userID, roleID string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("role_id", roleID).Error
}
