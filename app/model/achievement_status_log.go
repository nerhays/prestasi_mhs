package model

import "time"

type AchievementStatusLog struct {
	ID                     string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AchievementReferenceID string    `gorm:"type:uuid"`
	OldStatus              string
	NewStatus              string
	ChangedBy              string    `gorm:"type:uuid"`
	Note                   *string
	CreatedAt              time.Time
}
