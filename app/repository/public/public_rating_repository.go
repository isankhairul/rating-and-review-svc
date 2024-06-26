package publicrepository

import (
	"context"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/request/public"
	"go-klikdokter/pkg/util"
	"math"
	"reflect"
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
	GetRatingsBySourceTypeAndActor(sourceType, sourceUID string, filter publicrequest.GetRatingBySourceTypeAndActorFilter) ([]entity.RatingsCol, error)
	GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error)
	GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error)
	CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) (*entity.RatingSubHelpfulCol, error)
	UpdateStatusRatingSubHelpful(id primitive.ObjectID, currentStatus bool) error
	GetRatingSubHelpfulByRatingSubAndActor(ratingSubId, userIdLegacy string) (*entity.RatingSubHelpfulCol, error)
	UpdateCounterRatingSubmission(id primitive.ObjectID, currentCounter int64) error
	GetPublicRatingsByParams(limit, page, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsCol, *base.Pagination, error)
	GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmisson, error)
	CountRatingSubsByRatingIdAndValue(ratingId, value string) (int64, error)
	GetPublicRatingSubmissions(limit, page, dir int, sort string, filter publicrequest.FilterRatingSubmission) ([]entity.RatingSubmisson, *base.Pagination, error)
	GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId, sourceType string) (*entity.RatingFormulaCol, error)
	UpdateRatingSubDisplayNameByIdLegacy(input request.UpdateRatingSubDisplayNameRequest) error
	GetListRatingBySourceTypeAndUID(sourceType, sourceUID string) ([]entity.RatingsCol, error)
}

func NewPublicRatingRepository(db *mongo.Database) PublicRatingRepository {
	return &publicRatingRepo{db}
}

var bsonStatus = bson.D{{"status", true}}

func (r *publicRatingRepo) GetRatingsBySourceTypeAndActor(sourceType, sourceUID string, filter publicrequest.GetRatingBySourceTypeAndActorFilter) ([]entity.RatingsCol, error) {
	var results []entity.RatingsCol

	bsonSourceType := bson.D{{Key: "source_type", Value: sourceType}}
	bsonSourceUid := bson.D{{Key: "source_uid", Value: sourceUID}}
	bsonStatus := bson.D{{Key: "status", Value: true}}
	bsonRatingType := bson.D{}

	if len(filter.RatingType) > 0 {
		bsonRatingType = bson.D{{Key: "rating_type", Value: bson.D{{Key: "$in", Value: filter.RatingType}}}}
	}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonStatus,
			bsonSourceType,
			bsonSourceUid,
			bsonRatingType,
		}},
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	// sort desc by rating_type
	cursor, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).
		Find(ctx, bsonFilter,
			&options.FindOptions{
				Sort: bson.D{bson.E{Key: "rating_type", Value: -1}},
			})
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
			"status":               true,
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

func (r *publicRatingRepo) UpdateStatusRatingSubHelpful(id primitive.ObjectID, currentStatus bool) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: !currentStatus}}}}
	_, err := r.db.Collection(entity.RatingSubHelpfulCol{}.CollectionName()).UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *publicRatingRepo) GetRatingSubHelpfulByRatingSubAndActor(ratingSubId, userIdLegacy string) (*entity.RatingSubHelpfulCol, error) {
	var ratingSubHelpful entity.RatingSubHelpfulCol
	bsonRatingSubId := bson.D{{Key: "rating_submission_id", Value: ratingSubId}}
	bsonUserIdLegacy := bson.D{{Key: "user_id_legacy", Value: userIdLegacy}}

	filter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonRatingSubId,
			bsonUserIdLegacy,
		},
	}}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection(entity.RatingSubHelpfulCol{}.CollectionName()).FindOne(ctx, filter).Decode(&ratingSubHelpful)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &ratingSubHelpful, nil
}

func (r *publicRatingRepo) UpdateCounterRatingSubmission(id primitive.ObjectID, currentCounter int64) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	// timeUpdate := time.Now().In(util.Loc)

	helpfulCounter, err := countRatingSubHelpful(r, id.Hex())
	if err != nil {
		return err
	}
	if currentCounter != helpfulCounter {
		currentCounter = helpfulCounter
	}

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": bson.M{"like_counter": int(currentCounter)}}

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err1 := r.db.Collection("ratingSubCol").FindOneAndUpdate(context.Background(), filter, update, &options.FindOneAndUpdateOptions{})
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
	bsonRatingSubId := bson.D{{Key: "rating_submission_id", Value: ratingSubId}}
	filter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonRatingSubId,
			bsonStatus,
		},
	}}
	counter, err := r.db.Collection("ratingSubHelpfulCol").CountDocuments(context.Background(), filter, &options.CountOptions{})
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (r *publicRatingRepo) GetPublicRatingsByParams(limit, page, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsCol, *base.Pagination, error) {
	var results []entity.RatingsCol
	limit64 := int64(limit)
	bsonSourceUid := bson.D{}
	bsonSourceType := bson.D{}
	bsonRatingType := bson.D{}

	if len(filter.SourceUid) > 0 {
		bsonSourceUid = bson.D{{Key: "source_uid", Value: bson.D{{Key: "$in", Value: filter.SourceUid}}}}
	}

	if filter.SourceType != "" {
		if filter.SourceType != "all" {
			bsonSourceType = bson.D{{Key: "source_type", Value: filter.SourceType}}
		}
	}

	if len(filter.RatingType) > 0 {
		bsonRatingType = bson.D{{Key: "rating_type", Value: bson.D{{Key: "$in", Value: filter.RatingType}}}}
	}

	var bsonFilter = bson.D{{Key: "$and",
		Value: bson.A{
			bsonSourceType,
			bsonSourceUid,
			bsonStatus,
			bsonRatingType,
		},
	},
	}

	collectionName := "ratingsCol"
	skip := int64(page)*limit64 - limit64
	cursor, err := r.db.Collection(collectionName).
		Find(context.Background(), bsonFilter,
			&options.FindOptions{
				Sort:  bson.D{bson.E{Key: sort, Value: dir}},
				Limit: &limit64,
				Skip:  &skip,
			})
	if err != nil {
		return nil, nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}
	pagination := paginate(r, page, limit64, results, collectionName, bsonFilter)
	return results, pagination, nil
}

func (r *publicRatingRepo) GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmisson, error) {
	var results []entity.RatingSubmisson
	bsonRatingId := bson.D{{Key: "rating_id", Value: ratingId}}
	bsonCancelled := bson.D{{Key: "cancelled", Value: false}}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonRatingId,
			bsonCancelled,
		},
	}}

	cursor, err := r.db.Collection("ratingSubCol").Find(context.Background(), bsonFilter)
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

func (r *publicRatingRepo) CountRatingSubsByRatingIdAndValue(ratingId, value string) (int64, error) {
	bsonRatingId := bson.D{{Key: "rating_id", Value: ratingId}}
	bsonValue := bson.M{"value": bson.M{"$regex": value, "$options": "im"}}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonRatingId,
			bsonValue,
		}},
	}
	counter, err := r.db.Collection("ratingSubCol").CountDocuments(context.Background(), bsonFilter, &options.CountOptions{})
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (r *publicRatingRepo) GetPublicRatingSubmissions(limit, page, dir int, sort string, filter publicrequest.FilterRatingSubmission) ([]entity.RatingSubmisson, *base.Pagination, error) {
	var results []entity.RatingSubmisson
	limit64 := int64(limit)
	bsonCancelled := bson.D{{Key: "cancelled", Value: false}}

	bsonUserIdLegacy := bson.D{}
	bsonSourceTransId := bson.D{}
	bsonValue := bson.D{}
	bsonRangeDate := bson.D{}
	bsonCreatedAt := bson.D{}

	if len(filter.UserIdLegacy) > 0 {
		bsonUserIdLegacy = bson.D{{Key: "user_id_legacy", Value: bson.D{{Key: "$in", Value: filter.UserIdLegacy}}}}
	}
	if len(filter.SourceTransID) > 0 {
		bsonSourceTransId = bson.D{{Key: "source_trans_id", Value: bson.D{{Key: "$in", Value: filter.SourceTransID}}}}
	}
	if filter.Value != "" {
		bsonValue = bson.D{{Key: "value", Value: filter.Value}}
	}
	if filter.StartDate != "" {
		startDate, _ := util.ConvertToDateTime(filter.StartDate + " 00:00:00")
		bsonRangeDate = append(bsonRangeDate, bson.E{Key: "$gte", Value: primitive.NewDateTimeFromTime(startDate)})
	}
	if filter.EndDate != "" {
		endDate, _ := util.ConvertToDateTime(filter.EndDate + " 23:59:59")
		bsonRangeDate = append(bsonRangeDate, bson.E{Key: "$lt", Value: primitive.NewDateTimeFromTime(endDate)})
	}
	if len(bsonRangeDate) > 0 {
		bsonCreatedAt = bson.D{{Key: "created_at", Value: bsonRangeDate}}
	}

	var bsonFilter = bson.D{}
	if filter.LikertFilter.RatingId != "" && len(filter.LikertFilter.Value) != 0 {
		bsonRatingType := bson.D{{Key: "tagging.rating_id", Value: filter.LikertFilter.RatingId}}
		bsonLikertVal := bson.D{{Key: "tagging.value", Value: bson.D{{Key: "$in", Value: filter.LikertFilter.Value}}}}

		bsonFilter = bson.D{{Key: "$and", Value: bson.A{
			bsonRatingType,
			bsonLikertVal,
			bsonCancelled,
			bsonCreatedAt,
		}}}
	} else {
		bsonRatingID := bson.D{{Key: "rating_id", Value: bson.D{{Key: "$in", Value: filter.RatingID}}}}
		bsonFilter = bson.D{{Key: "$and", Value: bson.A{
			bsonRatingID,
			bsonCancelled,
			bsonCreatedAt,
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonUserIdLegacy,
				}}},
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonSourceTransId,
				}}},
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonValue,
				}}},
		}}}
	}

	collectionName := "ratingSubCol"
	skip := int64(page)*limit64 - limit64
	cursor, err := r.db.Collection(collectionName).
		Find(context.Background(), bsonFilter,
			&options.FindOptions{
				Sort:  bson.D{bson.E{Key: sort, Value: dir}},
				Limit: &limit64,
				Skip:  &skip,
			})
	if err != nil {
		return nil, nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}
	pagination := paginate(r, page, limit64, results, collectionName, bsonFilter)
	return results, pagination, nil
}

func (r *publicRatingRepo) GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId, sourceType string) (*entity.RatingFormulaCol, error) {
	var ratingFormula entity.RatingFormulaCol

	bsonRatingTypeId := bson.D{{Key: "rating_type_id", Value: ratingTypeId}}
	bsonSourceType := bson.D{{Key: "source_type", Value: sourceType}}
	bsonStatus := bson.D{{Key: "status", Value: true}}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonRatingTypeId,
			bsonSourceType,
			bsonStatus,
		}},
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection("ratingFormulaCol").FindOne(ctx, bsonFilter).Decode(&ratingFormula)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &ratingFormula, nil
}

func (r *publicRatingRepo) UpdateRatingSubDisplayNameByIdLegacy(input request.UpdateRatingSubDisplayNameRequest) error {
	filter := bson.D{{Key: "user_id_legacy", Value: input.UserIdLegacy}}
	update := bson.M{"$set": bson.M{"display_name": input.DisplayName}}

	_, err := r.db.Collection("ratingSubCol").UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *publicRatingRepo) GetListRatingBySourceTypeAndUID(sourceType, sourceUID string) ([]entity.RatingsCol, error) {
	var results []entity.RatingsCol
	ratingType := [2]string{"rating_like_dislike", "review_for_layanan"}

	bsonFilter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "source_type", Value: sourceType}},
		bson.D{{Key: "source_uid", Value: sourceUID}},
		bson.D{{Key: "rating_type", Value: bson.D{{Key: "$in", Value: ratingType}}}},
		bson.D{{Key: "status", Value: true}},
	}}}

	cursor, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).Find(context.Background(), bsonFilter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, nil
		}
		return nil, err
	}
	if err = cursor.All(context.Background(), &results); err != nil {
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func paginate(r *publicRatingRepo, page int, limit int64, result interface{}, collectionName string, filter bson.D) *base.Pagination {
	var pagination base.Pagination
	s := reflect.ValueOf(result)
	totalRecord, err := r.db.Collection(collectionName).CountDocuments(context.Background(), filter, &options.CountOptions{})
	if err != nil {
		return nil
	}
	pagination.Page = page
	pagination.Limit = int(limit)
	pagination.TotalRecords = totalRecord
	pagination.TotalPage = int(math.Ceil(float64(pagination.TotalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(s.Len())
	return &pagination
}
