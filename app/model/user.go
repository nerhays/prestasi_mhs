package model

import "time"

type User struct {
	ID           string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username     string    `gorm:"size:50;unique;not null" json:"username"`
	Email        string    `gorm:"size:100;unique;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	FullName     string    `gorm:"size:100;not null" json:"full_name"`
	RoleID       string    `gorm:"type:uuid;not null" json:"role_id"`
	Role         Role      `gorm:"foreignKey:RoleID" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
