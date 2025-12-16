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
	FindDeletedByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error)
	FindByIDs(ctx context.Context, ids []string) ([]model.Achievement, error)
	AddAttachment(ctx context.Context, mongoID string, att model.Attachment) error
	FindByID(ctx context.Context, id string) (*model.Achievement, error)
	CountByType(ctx context.Context) (map[string]int64, error)
	
	FindAll(ctx context.Context) ([]model.Achievement, error)
	Update(ctx context.Context, id string, payload *model.Achievement) (*model.Achievement, error)
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
func (r *achievementRepository) FindDeletedByStudentID(ctx context.Context, studentID string) ([]model.Achievement, error) {
    filter := bson.M{
        "studentId": studentID,
        "isDeleted": true,
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

    return results, nil
}
func (r *achievementRepository) FindByIDs(ctx context.Context, ids []string) ([]model.Achievement, error) {
    if len(ids) == 0 {
        return []model.Achievement{}, nil
    }
    // Convert to ObjectID when IDs are hex; if stored as hex strings of ObjectId, use primitive.ObjectIDFromHex
    // In your model, mongo_achievement_id stored as hex string; we can query by _id: ObjectId(hex)
    objIDs := make([]primitive.ObjectID, 0, len(ids))
    for _, h := range ids {
        oid, err := primitive.ObjectIDFromHex(h)
        if err != nil {
            continue // skip invalid
        }
        objIDs = append(objIDs, oid)
    }
    filter := bson.M{"_id": bson.M{"$in": objIDs}, "isDeleted": bson.M{"$ne": true}}
    cur, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)
    var out []model.Achievement
    for cur.Next(ctx) {
        var a model.Achievement
        if err := cur.Decode(&a); err != nil {
            return nil, err
        }
        out = append(out, a)
    }
    return out, nil
}
func (r *achievementRepository) AddAttachment(
	ctx context.Context,
	mongoID string,
	att model.Attachment,
) error {
	objID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateByID(ctx, objID, bson.M{
		"$push": bson.M{"attachments": att},
		"$set":  bson.M{"updatedAt": time.Now()},
	})
	return err
}
func (r *achievementRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Achievement, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var ac model.Achievement
	err = r.collection.FindOne(
		ctx,
		bson.M{"_id": objID},
	).Decode(&ac)

	if err != nil {
		return nil, err
	}

	return &ac, nil
}
func (r *achievementRepository) CountByType(ctx context.Context) (map[string]int64, error) {
	cursor, err := r.collection.Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id":   "$achievementType",
			"count": bson.M{"$sum": 1},
		}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]int64)
	for cursor.Next(ctx) {
		var row struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		result[row.ID] = row.Count
	}
	return result, nil
}
func (r *achievementRepository) Update(
	ctx context.Context,
	id string,
	payload *model.Achievement,
) (*model.Achievement, error) {

	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := r.collection.UpdateByID(
		ctx,
		objID,
		bson.M{"$set": payload},
	)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}
func (r *achievementRepository) FindAll(
	ctx context.Context,
) ([]model.Achievement, error) {

	cur, err := r.collection.Find(ctx, bson.M{
		"isDeleted": bson.M{"$ne": true},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var res []model.Achievement
	for cur.Next(ctx) {
		var a model.Achievement
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, nil
}

