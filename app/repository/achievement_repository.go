package repository

import (
	"context"
	"time"

	"github.com/nerhays/prestasi_uas/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	Create(ctx context.Context, ac *model.Achievement) (*model.Achievement, error)
	FindByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error)
	SoftDelete(ctx context.Context, mongoID string) error
}

type achievementRepository struct {
	collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) AchievementRepository {
	return &achievementRepository{
		collection: db.Collection("achievements"),
	}
}

func (r *achievementRepository) Create(ctx context.Context, ac *model.Achievement) (*model.Achievement, error) {
	ac.CreatedAt = time.Now()
	ac.UpdatedAt = time.Now()

	res, err := r.collection.InsertOne(ctx, ac)
	if err != nil {
		return nil, err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		ac.ID = oid
	}

	return ac, nil
}

func (r *achievementRepository) FindByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error) {
	filter := bson.M{
		"studentId": studentID,
		"isDeleted": bson.M{"$ne": true}, //tidak tampilkan yang soft delete
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []model.Achievement
	for cur.Next(ctx) {
		var ac model.Achievement
		if err := cur.Decode(&ac); err != nil {
			return nil, err
		}
		results = append(results, ac)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *achievementRepository) SoftDelete(ctx context.Context, mongoID string) error {
	objID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"isDeleted": true,
				"deletedAt": time.Now(),
			},
		},
	)

	return err
}
