package publicrepository

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	publicrequest "go-klikdokter/app/model/request/public"
	publicresponse "go-klikdokter/app/model/response/public"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"reflect"
	"time"
)

type publicRatingMpRepo struct {
	db *mongo.Database
}

type PublicRatingMpRepository interface {
	GetListRatingBySourceTypeAndUID(sourceType, sourceUID string) ([]entity.RatingsMpCol, error)
	GetPublicRatingSubmissions(limit, page, dir int, sort string, filter publicrequest.FilterRatingSubmissionMp) ([]entity.RatingSubmissionMp, *base.Pagination, error)
	GetPublicRatingsByParams(limit, page, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsMpCol, *base.Pagination, error)
	CountRatingSubsByRatingIdAndValue(ratingId, value string) (int64, error)
	GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmissionMp, error)
	GetSumCountRatingSubsByRatingId(ratingId string) (*publicresponse.PublicSumCountRatingSummaryMp, error)
	GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId, sourceType string) (*entity.RatingFormulaCol, error)
}

func NewPublicRatingMpRepository(db *mongo.Database) PublicRatingMpRepository {
	return &publicRatingMpRepo{db}
}

func (r *publicRatingMpRepo) GetListRatingBySourceTypeAndUID(sourceType, sourceUID string) ([]entity.RatingsMpCol, error) {
	var results []entity.RatingsMpCol
	arrRatingType := viper.GetStringSlice("rating-type-mp")

	bsonFilter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "source_type", Value: sourceType}},
		bson.D{{Key: "source_uid", Value: sourceUID}},
		bson.D{{Key: "rating_type", Value: bson.D{{Key: "$in", Value: arrRatingType}}}},
		bson.D{{Key: "status", Value: true}},
	}}}

	cursor, err := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).Find(context.Background(), bsonFilter)
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

func (r *publicRatingMpRepo) GetPublicRatingSubmissions(limit, page, dir int, sort string, filter publicrequest.FilterRatingSubmissionMp) ([]entity.RatingSubmissionMp, *base.Pagination, error) {
	var results []entity.RatingSubmissionMp
	limit64 := int64(limit)
	bsonCancelled := bson.D{{Key: "cancelled", Value: false}}

	bsonUserIdLegacy := bson.D{}
	bsonSourceTransId := bson.D{}
	bsonValue := bson.D{}
	bsonIsWithMedia := bson.D{}

	if len(filter.UserIdLegacy) > 0 {
		bsonUserIdLegacy = bson.D{{Key: "user_id_legacy", Value: bson.D{{Key: "$in", Value: filter.UserIdLegacy}}}}
	}
	if len(filter.SourceTransID) > 0 {
		bsonSourceTransId = bson.D{{Key: "source_trans_id", Value: bson.D{{Key: "$in", Value: filter.SourceTransID}}}}
	}
	if filter.Value != "" {
		bsonValue = bson.D{{Key: "value", Value: filter.Value}}
	}
	if filter.IsWithMedia != nil {
		bsonIsWithMedia = bson.D{{Key: "is_with_media", Value: filter.IsWithMedia}}
	}

	var bsonFilter = bson.D{}
	if filter.LikertFilter.RatingId != "" && len(filter.LikertFilter.Value) != 0 {
		bsonRatingType := bson.D{{Key: "tagging.rating_id", Value: filter.LikertFilter.RatingId}}
		bsonLikertVal := bson.D{{Key: "tagging.value", Value: bson.D{{Key: "$in", Value: filter.LikertFilter.Value}}}}

		bsonFilter = bson.D{{Key: "$and", Value: bson.A{
			bsonRatingType,
			bsonLikertVal,
			bsonCancelled,
			bsonIsWithMedia,
		}}}
	} else {
		bsonRatingID := bson.D{{Key: "rating_id", Value: bson.D{{Key: "$in", Value: filter.RatingID}}}}
		bsonFilter = bson.D{{Key: "$and", Value: bson.A{
			bsonRatingID,
			bsonCancelled,
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
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonIsWithMedia,
				}}},
		}}}
	}

	collectionName := entity.RatingSubmissionMp{}.CollectionName()
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
	pagination := paginateMp(r.db, page, limit64, results, collectionName, bsonFilter)
	return results, pagination, nil
}

func (r *publicRatingMpRepo) GetPublicRatingsByParams(limit, page, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsMpCol, *base.Pagination, error) {
	var results []entity.RatingsMpCol
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
	}}

	collectionName := entity.RatingsMpCol{}.CollectionName()
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
	pagination := paginateMp(r.db, page, limit64, results, collectionName, bsonFilter)
	return results, pagination, nil
}

func (r *publicRatingMpRepo) CountRatingSubsByRatingIdAndValue(ratingId, value string) (int64, error) {
	bsonRatingId := bson.D{{Key: "rating_id", Value: ratingId}}
	bsonValue := bson.M{"value": bson.M{"$regex": value, "$options": "im"}}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonRatingId,
			bsonValue,
		}},
	}
	counter, err := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName()).CountDocuments(context.Background(), bsonFilter, &options.CountOptions{})
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (r *publicRatingMpRepo) GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmissionMp, error) {
	var results []entity.RatingSubmissionMp
	bsonRatingId := bson.D{{Key: "rating_id", Value: ratingId}}
	bsonCancelled := bson.D{{Key: "cancelled", Value: false}}

	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonRatingId,
			bsonCancelled,
		},
	}}

	cursor, err := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName()).Find(context.Background(), bsonFilter)
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

func (r *publicRatingMpRepo) GetSumCountRatingSubsByRatingId(ratingId string) (*publicresponse.PublicSumCountRatingSummaryMp, error) {
	var results []publicresponse.PublicSumCountRatingSummaryMp
	bsonRatingIdAndCancelled := bson.D{
		{Key: "rating_id", Value: ratingId},
		{Key: "cancelled", Value: false},
	}

	pipeline := bson.A{
		bson.D{{Key: "$match",
			Value: bsonRatingIdAndCancelled,
		}},
		bson.D{{
			Key:   "$addFields",
			Value: bson.D{{"convertedValue", bson.D{{"$toInt", "$value"}}}}}},
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: primitive.Null{}},
				{Key: "sum", Value: bson.D{{"$sum", bson.D{{"$multiply", bson.A{"$convertedValue"}}}}}},
				{Key: "count", Value: bson.D{{"$sum", 1}}},
			}}},
	}

	cursor, err := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName()).Aggregate(context.Background(), pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &publicresponse.PublicSumCountRatingSummaryMp{}, nil
		}
		return nil, err
	}

	if err = cursor.All(context.Background(), &results); err != nil {
		if err != nil {
			return nil, err
		}
	}

	if len(results) == 0 {
		return nil, errors.New("repository data RatingSubmission not found")
	}

	return &results[0], nil
}

func (r *publicRatingMpRepo) GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId, sourceType string) (*entity.RatingFormulaCol, error) {
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
	err := r.db.Collection(entity.RatingFormulaCol{}.CollectionName()).FindOne(ctx, bsonFilter).Decode(&ratingFormula)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &ratingFormula, nil
}

func paginateMp(db *mongo.Database, page int, limit int64, result interface{}, collectionName string, filter bson.D) *base.Pagination {
	var pagination base.Pagination
	s := reflect.ValueOf(result)
	totalRecord, err := db.Collection(collectionName).CountDocuments(context.Background(), filter, &options.CountOptions{})
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
