package repository

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"gorm.io/gorm"
)

type AchievementStatusLogRepository interface {
	Create(log *model.AchievementStatusLog) error
	FindByReferenceID(refID string) ([]model.AchievementStatusLog, error)
}

type achievementStatusLogRepo struct {
	db *gorm.DB
}

func NewAchievementStatusLogRepository(db *gorm.DB) AchievementStatusLogRepository {
	return &achievementStatusLogRepo{db}
}

func (r *achievementStatusLogRepo) Create(log *model.AchievementStatusLog) error {
	return r.db.Create(log).Error
}

func (r *achievementStatusLogRepo) FindByReferenceID(refID string) ([]model.AchievementStatusLog, error) {
	var logs []model.AchievementStatusLog
	err := r.db.
		Where("achievement_reference_id = ?", refID).
		Order("created_at ASC").
		Find(&logs).Error
	return logs, err
}
