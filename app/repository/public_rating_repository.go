package repository

import (
	"context"
	"errors"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/pkg/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type publicRatingRepo struct {
	db *mongo.Database
}

type PublicRatingRepository interface {
	GetRatingsBySourceTypeAndActor(sourceType, sourceUID string) ([]entity.RatingsCol, error)
	GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error)
	GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error)
	CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) (*entity.RatingSubHelpfulCol, error)
	UpdateCounterRatingSubmission(id primitive.ObjectID, currentCounter int) error
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

func (r *publicRatingRepo) CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) (*entity.RatingSubHelpfulCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var ratingSubHelpful entity.RatingSubHelpfulCol

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := r.db.Collection(entity.RatingSubHelpfulCol{}.CollectionName()).InsertOne(ctx, bson.M{
			"rating_submission_id": input.RatingSubmissionID,
			"user_id":              input.UserID,
			"user_id_legacy":       input.UserIDLegacy,
			"ip_address":           input.IPAddress,
			"user_agent":           input.UserAgent,
			"created_at":           time.Now().In(util.Loc),
			"updated_at":           time.Now().In(util.Loc),
		})

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		ratingSubHelpful.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	})
	if errTransaction != nil {
		return nil, errTransaction
	}
	return &ratingSubHelpful, nil
}

func (r *publicRatingRepo) UpdateCounterRatingSubmission(id primitive.ObjectID, currentCounter int) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var timeUpdate time.Time
	timeUpdate = time.Now().In(util.Loc)
	var counter int64
	counter = int64(currentCounter) + 1

	helpfulCounter, err := countRatingSubHelpful(r, id.Hex())
	if err != nil {
		return err
	}
	if counter != helpfulCounter {
		return nil
	}

	ratingSubmiss := entity.RatingSubmisson{
		LikeCounter: int(counter),
		UpdatedAt:   timeUpdate,
	}
	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", ratingSubmiss}}

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err1 := r.db.Collection("ratingSubCol").FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
		if err1 != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err1.Err()
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	})
	if errTransaction != nil {
		return errTransaction
	}
	return nil
}

func countRatingSubHelpful(r *publicRatingRepo, ratingSubId string) (int64, error) {
	filterHelpful := bson.D{{"rating_submission_id", ratingSubId}}
	counter, err := r.db.Collection("ratingSubHelpfulCol").CountDocuments(context.Background(), filterHelpful, &options.CountOptions{})
	if err != nil {
		return 0, err
	}
	return counter, nil
}
