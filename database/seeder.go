package database

import (
	"log"

	"github.com/nerhays/prestasi_uas/app/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedAdminUser(db)
}

func seedAdminUser(db *gorm.DB) {
	// cek user admin
	var count int64
	if err := db.Model(&model.User{}).
		Where("username = ?", "admin").
		Count(&count).Error; err != nil {
		log.Printf("[SEED] failed to check admin user: %v", err)
		return
	}

	if count > 0 {
		log.Println("[SEED] admin user already exists, skip seeding")
		return
	}

	// cari role "Admin"
	var role model.Role
	if err := db.Where("name = ?", "Admin").First(&role).Error; err != nil {
		log.Printf("[SEED] role 'Admin' not found: %v", err)
		return
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[SEED] failed to hash password: %v", err)
		return
	}

	admin := model.User{
		Username:     "admin",
		Email:        "admin@example.com",
		FullName:     "Administrator",
		PasswordHash: string(hashed),
		RoleID:       role.ID,
		IsActive:     true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("[SEED] failed to create admin user: %v", err)
		return
	}

	log.Println("[SEED] admin user created: username=admin, password=admin123")
}
