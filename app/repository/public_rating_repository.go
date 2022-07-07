package repository

import (
	"context"
	"errors"
	"go-klikdokter/app/model/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type publicRatingRepo struct {
	db *mongo.Database
}

type PublicRatingRepository interface {
	GetRatingsBySourceTypeAndActor(sourceType, sourceUID string) ([]entity.RatingsCol, error)
	GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error)
	GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error)
}

func NewPublicRatingRepository(db *mongo.Database) PublicRatingRepository {
	return &publicRatingRepo{db}
}

func (r *publicRatingRepo) GetRatingsBySourceTypeAndActor(sourceType, sourceUID string) ([]entity.RatingsCol, error) {
	var results []entity.RatingsCol

	bsonSourceType := bson.D{{"source_type", sourceType}}
	bsonSourceUid := bson.D{{"source_uid", sourceUID}}
	bsonStatus := bson.D{{"status", true}}

	bsonFilter := bson.D{{"$and",
		bson.A{
			bsonStatus,
			bsonSourceType,
			bsonSourceUid,
		}},
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	cursor, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).Find(ctx, bsonFilter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, nil
		}
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *publicRatingRepo) GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
	var ratingTypeLikert entity.RatingTypesLikertCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection("ratingTypesLikertCol").FindOne(ctx, bson.M{"_id": id}).Decode(&ratingTypeLikert)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &ratingTypeLikert, nil
}

func (r *publicRatingRepo) GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	var ratingTypeNum entity.RatingTypesNumCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection("ratingTypesNumCol").FindOne(ctx, bson.M{"_id": id}).Decode(&ratingTypeNum)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &ratingTypeNum, nil
}
