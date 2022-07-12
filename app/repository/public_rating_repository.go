package repository

import (
	"context"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/pkg/util"
	"math"
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
	GetPublicRatingsByParams(limit, page, dir int, sort string, filter request.FilterRatingSummary) ([]entity.RatingsCol, *base.Pagination, error)
	GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmisson, error)
	GetPublicRatingSubmissions(limit, page, dir int, sort string, filter request.FilterRatingSubmission) ([]entity.RatingSubmisson, *base.Pagination, error)
}

func NewPublicRatingRepository(db *mongo.Database) PublicRatingRepository {
	return &publicRatingRepo{db}
}

func (r *publicRatingRepo) GetRatingsBySourceTypeAndActor(sourceType, sourceUID string) ([]entity.RatingsCol, error) {
	var results []entity.RatingsCol

	bsonSourceType := bson.D{{Key: "source_type", Value: sourceType}}
	bsonSourceUid := bson.D{{Key: "source_uid", Value: sourceUID}}
	bsonStatus := bson.D{{Key: "status", Value: true}}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
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
	timeUpdate := time.Now().In(util.Loc)
	counter := int64(currentCounter) + 1

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
	filter := bson.D{{Key: "_id", Value: id}}
	data := bson.D{{Key: "$set", Value: ratingSubmiss}}

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
	filterHelpful := bson.D{{Key: "rating_submission_id", Value: ratingSubId}}
	counter, err := r.db.Collection("ratingSubHelpfulCol").CountDocuments(context.Background(), filterHelpful, &options.CountOptions{})
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (r *publicRatingRepo) GetPublicRatingsByParams(limit, page, dir int, sort string, filter request.FilterRatingSummary) ([]entity.RatingsCol, *base.Pagination, error) {
	var results []entity.RatingsCol
	var allResults []bson.D
	var pagination base.Pagination

	bsonSourceUid := bson.D{}
	bsonSourceType := bson.D{}
	bsonRatingType := bson.D{}

	if len(filter.SourceUid) > 0 {
		bsonSourceUid = bson.D{{Key: "source_uid", Value: filter.SourceUid}}
	}
	if filter.SourceType != "" {
		if filter.SourceType != "all" {
			bsonSourceType = bson.D{{Key: "source_type", Value: filter.SourceType}}
		}
	}
	if len(filter.RatingType) > 0 {
		bsonRatingType = bson.D{{Key: "rating_type", Value: bson.D{{Key: "$in", Value: filter.RatingType}}}}
	}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonStatus,
			bsonSourceType,
			bsonSourceUid,
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonRatingType,
				}}},
		},
	},
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	cursor, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).
		Find(ctx, bsonFilter,
			newMongoPaginate(limit, page).getPaginatedOpts().
				SetSort(bson.D{{Key: sort, Value: dir}}))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, nil, nil
		}
		return nil, nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}

	crsr, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).Find(ctx, bsonFilter)
	if err != nil {
		return nil, nil, err
	}
	if err = crsr.All(ctx, &allResults); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}

	totalRecords := int64(len(allResults))
	pagination.Limit = limit
	pagination.Page = page
	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(len(results))

	return results, &pagination, nil
}

func (r *publicRatingRepo) GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmisson, error) {
	var results []entity.RatingSubmisson

	cursor, err := r.db.Collection("ratingSubCol").Find(context.Background(), bson.M{"rating_id": ratingId})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, nil
		}
		return nil, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *publicRatingRepo) GetPublicRatingSubmissions(limit, page, dir int, sort string, filter request.FilterRatingSubmission) ([]entity.RatingSubmisson, *base.Pagination, error) {
	var results []entity.RatingSubmisson
	var allResults []bson.D
	var pagination base.Pagination

	bsonRatingID := bson.D{{Key: "rating_id", Value: bson.D{{Key: "$in", Value: filter.RatingID}}}}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	cursor, err := r.db.Collection("ratingSubCol").
		Find(ctx, bsonRatingID,
			newMongoPaginate(limit, page).getPaginatedOpts().
				SetSort(bson.D{{Key: sort, Value: dir}}))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, nil, nil
		}
		return nil, nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}
	crsr, err := r.db.Collection("ratingSubCol").Find(ctx, bsonRatingID)
	if err != nil {
		return nil, nil, err
	}
	if err = crsr.All(ctx, &allResults); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}

	totalRecords := int64(len(allResults))
	pagination.Limit = limit
	pagination.Page = page
	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(len(results))

	return results, &pagination, nil
}
