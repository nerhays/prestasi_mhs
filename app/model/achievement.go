package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `bson:"studentId" json:"studentId"`
	AchievementType string             `bson:"achievementType" json:"achievementType"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         map[string]any     `bson:"details" json:"details"`
	Points          float64            `bson:"points" json:"points"`
	Attachments     []any              `bson:"attachments" json:"attachments,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}
type AchievementStatistics struct {
	Total            int64            `json:"total"`
	ByType           map[string]int64 `json:"by_type"`
	ByStatus         map[string]int64 `json:"by_status"`
}
