package repository

import (
	"time"

	"github.com/nerhays/prestasi_uas/app/model"
	"gorm.io/gorm"
)

type AchievementReferenceRepository interface {
	CreateDraft(studentID string, mongoAchievementID string) (*model.AchievementReference, error)
	GetByID(id string) (*model.AchievementReference, error)
	Save(ref *model.AchievementReference) error
	CountByStudentIDs(studentIDs []string, status *model.AchievementStatus) (int64, error)
    FindByStudentIDs(studentIDs []string, status *model.AchievementStatus, limit, offset int) ([]model.AchievementReference, error)
	
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
func (r *achievementReferenceRepository) GetByID(id string) (*model.AchievementReference, error) {
	var ref model.AchievementReference
	if err := r.db.First(&ref, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *achievementReferenceRepository) Save(ref *model.AchievementReference) error {
	ref.UpdatedAt = time.Now()
	return r.db.Save(ref).Error
}

func (r *achievementReferenceRepository) CountByStudentIDs(studentIDs []string, status *model.AchievementStatus) (int64, error) {
    var count int64
    q := r.db.Model(&model.AchievementReference{}).Where("student_id IN ?", studentIDs)
    if status != nil {
        q = q.Where("status = ?", *status)
    }
    if err := q.Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}

func (r *achievementReferenceRepository) FindByStudentIDs(studentIDs []string, status *model.AchievementStatus, limit, offset int) ([]model.AchievementReference, error) {
    var refs []model.AchievementReference
    q := r.db.Where("student_id IN ?", studentIDs)
    if status != nil {
        q = q.Where("status = ?", *status)
    }
    if err := q.Order("submitted_at DESC NULLS LAST, created_at DESC").Limit(limit).Offset(offset).Find(&refs).Error; err != nil {
        return nil, err
    }
    return refs, nil
}

