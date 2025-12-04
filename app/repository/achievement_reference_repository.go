package repository

import (
	"github.com/nerhays/prestasi_uas/app/model"
	"gorm.io/gorm"
)

type AchievementReferenceRepository interface {
	CreateDraft(studentID string, mongoAchievementID string) (*model.AchievementReference, error)
}

type achievementReferenceRepository struct {
	db *gorm.DB
}

func NewAchievementReferenceRepository(db *gorm.DB) AchievementReferenceRepository {
	return &achievementReferenceRepository{db: db}
}

func (r *achievementReferenceRepository) CreateDraft(studentID string, mongoAchievementID string) (*model.AchievementReference, error) {
	ref := &model.AchievementReference{
		StudentID:          studentID,
		MongoAchievementID: mongoAchievementID,
		Status:             model.AchievementStatusDraft,
	}

	if err := r.db.Create(ref).Error; err != nil {
		return nil, err
	}

	return ref, nil
}
