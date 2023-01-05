package repository

import (
	"context"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	publicrequest "go-klikdokter/app/model/request/public"
	publicresponse "go-klikdokter/app/model/response/public"
	"go-klikdokter/pkg/util"
	"math"
	"reflect"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ratingMpRepo struct {
	db *mongo.Database
}

type RatingMpRepository interface {
	// Rating submission
	CreateRatingSubmission(input []request.SaveRatingSubmissionMp) (*[]entity.RatingSubmissionMp, error)
	UpdateRatingSubmission(input entity.RatingSubmissionMp, id primitive.ObjectID) error
	GetRatingSubmissionById(id primitive.ObjectID) (*entity.RatingSubmissionMp, error)
	GetRatingSubmissionByIdAndUser(id primitive.ObjectID, userIDLegacy string) (*entity.RatingSubmissionMp, error)
	GetListRatingSubmissions(filter request.RatingSubmissionMpFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingSubmissionMp, *base.Pagination, error)
	FindRatingSubmissionByUserIDAndRatingID(userId *string, ratingId string, sourceTransId string) (*entity.RatingSubmissionMp, error)
	FindRatingSubmissionBySourceTransID(sourceTransId string) (*entity.RatingSubmissionMp, error)
	GetPublicRatingsByParams(limit, page, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsMpCol, *base.Pagination, error)
	GetSumCountRatingSubsByRatingId(ratingId string) (*publicresponse.PublicSumCountRatingSummaryMp, error)
	CountRatingSubsByRatingIdAndValue(ratingId, value string) (int64, error)
	GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId, sourceType string) (*entity.RatingFormulaCol, error)

	// rating
	FindRatingBySourceUIDAndRatingType(sourceUID, ratingType string) (*entity.RatingsMpCol, error)
	GetRatingById(id primitive.ObjectID) (*entity.RatingsMpCol, error)
	GetRatingTypeLikertByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error)
	GetRatingByRatingTypeSourceUidAndSourceType(ratingTypeId, sourceUid, sourceType string) (*entity.RatingsMpCol, error)
	GetRatingTypeNumByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesNumCol, error)
	CreateRating(input request.SaveRatingRequest) (*entity.RatingsMpCol, error)
	GetRatingSubmissionByRatingId(id string) (*entity.RatingSubmissionMp, error)
	UpdateRating(id primitive.ObjectID, input request.BodyUpdateRatingRequest) (*entity.RatingsMpCol, error)
	DeleteRating(id primitive.ObjectID) error
	GetRatingsByParams(limit, page, dir int, sort string, filter request.RatingFilter) ([]entity.RatingsMpCol, *base.Pagination, error)

	// rating type
	FindRatingTypeNumByRatingType(ratingType string) (*entity.RatingTypesNumCol, error)
	FindRatingTypeNumByRatingTypeID(ratingTypeID primitive.ObjectID) (*entity.RatingTypesNumCol, error)
}

func NewRatingMpRepository(db *mongo.Database) RatingMpRepository {
	return &ratingMpRepo{db}
}

func (r *ratingMpRepo) GetRatingTypeLikertByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
	var ratingTypeLikert entity.RatingTypesLikertCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonStatus,
			bson.M{"_id": id},
		}},
	}
	err := r.db.Collection("ratingTypesLikertCol").FindOne(ctx, bsonFilter).Decode(&ratingTypeLikert)
	if err != nil {
		return nil, err
	}
	return &ratingTypeLikert, nil
}

func (r *ratingMpRepo) FindRatingBySourceUIDAndRatingType(sourceUID, ratingType string) (*entity.RatingsMpCol, error) {
	var rating entity.RatingsMpCol

	ratingColl := r.db.Collection(entity.RatingsMpCol{}.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingColl.FindOne(ctx, bson.M{"source_uid": sourceUID, "rating_type": ratingType}).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingMpRepo) CreateRatingSubmission(input []request.SaveRatingSubmissionMp) (*[]entity.RatingSubmissionMp, error) {
	ratingSubmissionColl := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var ratingSubmission []entity.RatingSubmissionMp
	var docs []interface{}
	for _, args := range input {
		dateNow := time.Now().In(util.Loc)
		docs = append(docs, bson.D{
			// {Key: "rating_id", Value: args.RatingID},
			{Key: "user_id", Value: args.UserID},
			{Key: "user_id_legacy", Value: args.UserIDLegacy},
			{Key: "display_name", Value: args.DisplayName},
			{Key: "comment", Value: args.Comment},
			{Key: "value", Value: args.Value},
			{Key: "avatar", Value: args.Avatar},
			{Key: "ip_address", Value: args.IPAddress},
			{Key: "user_agent", Value: args.UserAgent},
			{Key: "source_trans_id", Value: args.SourceTransID},
			{Key: "user_platform", Value: args.UserPlatform},
			{Key: "like_counter", Value: 0},
			{Key: "created_at", Value: dateNow},
			{Key: "updated_at", Value: dateNow},
			{Key: "source_uid", Value: args.SourceUID},
			{Key: "cancelled", Value: false},
			{Key: "cancelled_reason", Value: ""},
			{Key: "is_anonymous", Value: args.IsAnonymous},
			{Key: "source_type", Value: args.SourceType},
			{Key: "media_path", Value: args.MediaPath},
			{Key: "is_with_media", Value: args.IsWithMedia},
			{Key: "order_number", Value: args.OrderNumber},
			{Key: "rating_type_id", Value: args.RatingTypeID},
		})
	}
	if len(docs) < 1 {
		return nil, mongo.ErrNilValue
	}

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := ratingSubmissionColl.InsertMany(ctx, docs)

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		for _, args := range result.InsertedIDs {
			data := entity.RatingSubmissionMp{
				ID: args.(primitive.ObjectID),
			}
			ratingSubmission = append(ratingSubmission, data)
		}
		return nil
	})
	if errTransaction != nil {
		return nil, errTransaction
	}
	return &ratingSubmission, nil
}

func (r *ratingMpRepo) UpdateRatingSubmission(input entity.RatingSubmissionMp, id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)

	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", input}}

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err1 := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName()).FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
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

func (r *ratingMpRepo) GetRatingSubmissionById(id primitive.ObjectID) (*entity.RatingSubmissionMp, error) {
	var ratingSubmission entity.RatingSubmissionMp
	ratingSubmissionColl := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"_id": id}).Decode(&ratingSubmission)
	if err != nil {
		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingMpRepo) GetRatingSubmissionByIdAndUser(id primitive.ObjectID, userIDLegacy string) (*entity.RatingSubmissionMp, error) {
	var ratingSubmission entity.RatingSubmissionMp
	ratingSubmissionColl := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"_id": id, "user_id_legacy": &userIDLegacy}).Decode(&ratingSubmission)
	if err != nil {
		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingMpRepo) GetListRatingSubmissions(filter request.RatingSubmissionMpFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingSubmissionMp, *base.Pagination, error) {
	var results []entity.RatingSubmissionMp
	startDate, _ := time.Parse(util.LayoutDateOnly, filter.StartDate)
	endDate, _ := time.Parse(util.LayoutDateOnly, filter.EndDate)
	if errD := endDate.Before(startDate); errD == true {
		return nil, nil, errors.New("end_date can not before start_date")
	}
	bsonUserUid := bson.D{}
	bsonRating := bson.D{}
	bsonDate := bson.D{}
	bsonTransId := bson.D{}
	bsonIsWithMedia := bson.D{}
	if len(filter.UserIDLegacy) > 0 {
		bsonUserUid = bson.D{{Key: "user_id_legacy", Value: bson.D{{Key: "$in", Value: filter.UserIDLegacy}}}}
	}
	if len(filter.RatingID) > 0 {
		bsonRating = bson.D{{Key: "rating_id", Value: bson.D{{Key: "$in", Value: filter.RatingID}}}}
	}
	if len(filter.StartDate) > 0 && len(filter.EndDate) > 0 {
		bsonDate = bson.D{{Key: "created_at", Value: bson.D{{Key: "$gt", Value: startDate}, {Key: "$lt", Value: endDate.AddDate(0, 0, 1)}}}}
	}
	if filter.SourceTransID != "" {
		bsonTransId = bson.D{{Key: "source_trans_id", Value: filter.SourceTransID}}
	}
	if filter.IsWithMedia != nil {
		bsonIsWithMedia = bson.D{{Key: "is_with_media", Value: filter.IsWithMedia}}
	}

	arrSourceType := viper.GetStringSlice("source-type-mp")
	bsonSourceType := bson.D{{Key: "source_type", Value: bson.D{{Key: "$in", Value: arrSourceType}}}}

	filter1 := bson.D{{Key: "$and",
		Value: bson.A{
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonUserUid,
				}}},
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonRating,
				}}},
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonTransId,
				}}},
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonDate,
				},
			}},
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonIsWithMedia,
				},
			}},
			bsonSourceType,
		},
	},
	}

	limitPage := int(limit)
	cursor, err := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName()).
		Find(context.Background(), filter1,
			newMongoPaginate(limitPage, page).getPaginatedOpts().
				SetSort(bson.D{{Key: sort, Value: dir}}))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return results, nil, nil
		}
		return nil, nil, err
	}
	if err = cursor.All(context.Background(), &results); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}
	pagination := getPaginationMp(r.db, page, limit, results, "ratingSubMpCol", filter1)
	return results, pagination, nil
}

func (r *ratingMpRepo) FindRatingSubmissionByUserIDAndRatingID(userId *string, ratingId string, sourceTransId string) (*entity.RatingSubmissionMp, error) {
	var ratingSubmission entity.RatingSubmissionMp
	ratingSubmissionColl := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"user_id": &userId, "rating_id": ratingId, "source_trans_id": sourceTransId}).Decode(&ratingSubmission)
	if err != nil {

		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingMpRepo) FindRatingSubmissionBySourceTransID(sourceTransId string) (*entity.RatingSubmissionMp, error) {
	var ratingSubmission entity.RatingSubmissionMp
	ratingSubmissionColl := r.db.Collection(entity.RatingSubmissionMp{}.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"source_trans_id": sourceTransId}).Decode(&ratingSubmission)
	if err != nil {

		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingMpRepo) GetPublicRatingsByParams(limit, page, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsMpCol, *base.Pagination, error) {
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
	},
	}

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
	pagination := getPaginationMp(r.db, page, limit64, results, collectionName, bsonFilter)
	return results, pagination, nil
}

func (r *ratingMpRepo) GetSumCountRatingSubsByRatingId(ratingId string) (*publicresponse.PublicSumCountRatingSummaryMp, error) {
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

func (r *ratingMpRepo) CountRatingSubsByRatingIdAndValue(ratingId, value string) (int64, error) {
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

// RATING
func (r *ratingMpRepo) CreateRating(input request.SaveRatingRequest) (*entity.RatingsMpCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var rating entity.RatingsMpCol

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).InsertOne(ctx, bson.M{
			"name":            input.Name,
			"description":     input.Description,
			"source_uid":      input.SourceUid,
			"source_type":     input.SourceType,
			"rating_type":     input.RatingType,
			"rating_type_id":  input.RatingTypeId,
			"comment_allowed": input.CommentAllowed,
			"status":          input.Status,
			"created_at":      time.Now().In(util.Loc),
			"updated_at":      time.Now().In(util.Loc),
		})

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		rating.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	})
	if errTransaction != nil {
		return nil, errTransaction
	}
	return &rating, nil
}

func (r *ratingMpRepo) GetRatingById(id primitive.ObjectID) (*entity.RatingsMpCol, error) {
	var rating entity.RatingsMpCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).FindOne(ctx, bson.M{"_id": id}).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingMpRepo) GetRatingSubmissionByRatingId(id string) (*entity.RatingSubmissionMp, error) {
	var ratingSubmission entity.RatingSubmissionMp
	ratingSubmissionColl := r.db.Collection(ratingSubmission.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"rating_id": id}).Decode(&ratingSubmission)
	if err != nil {
		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingMpRepo) UpdateRating(id primitive.ObjectID, input request.BodyUpdateRatingRequest) (*entity.RatingsMpCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	rating := entity.RatingsMpCol{
		Name:           input.Name,
		Description:    input.Description,
		SourceUid:      input.SourceUid,
		SourceType:     input.SourceType,
		CommentAllowed: input.CommentAllowed,
		UpdatedAt:      time.Now().In(util.Loc),
	}
	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", rating}}
	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err1 := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
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
		return nil, errTransaction
	}
	return &rating, nil
}

func (r *ratingMpRepo) DeleteRating(id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	update := bson.D{{"$set", bson.D{{"status", false}}}}
	_, err := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ratingMpRepo) GetRatingsByParams(limit, page, dir int, sort string, filter request.RatingFilter) ([]entity.RatingsMpCol, *base.Pagination, error) {
	var results []entity.RatingsMpCol
	var allResults []bson.D
	var pagination base.Pagination

	bsonSourceUid := bson.D{}
	bsonRatingTypeId := bson.D{}
	bsonSourceType := bson.D{}

	if filter.SourceType != "" {
		if filter.SourceType != "all" {
			bsonSourceType = bson.D{{"source_type", filter.SourceType}}
		}
	}
	if len(filter.SourceUid) > 0 {
		bsonSourceUid = bson.D{{"source_uid", bson.D{{"$in", filter.SourceUid}}}}
	}
	if len(filter.RatingTypeId) > 0 {
		bsonRatingTypeId = bson.D{{"rating_type_id", bson.D{{"$in", filter.RatingTypeId}}}}
	}

	bsonFilter := bson.D{{"$and",
		bson.A{
			bsonStatus,
			bsonSourceType,
			bson.D{{"$or",
				bson.A{
					bsonSourceUid,
				}}},
			bson.D{{"$or",
				bson.A{
					bsonRatingTypeId,
				}}},
		},
	},
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	cursor, err := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).
		Find(ctx, bsonFilter,
			newMongoPaginate(limit, page).getPaginatedOpts().
				SetSort(bson.D{{sort, dir}}))

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

	crsr, err := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).Find(ctx, bsonFilter)
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

func (r *ratingMpRepo) GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId, sourceType string) (*entity.RatingFormulaCol, error) {
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

func (r *ratingMpRepo) GetRatingByRatingTypeSourceUidAndSourceType(ratingTypeId, sourceUid, sourceType string) (*entity.RatingsMpCol, error) {
	var rating entity.RatingsMpCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	bsonFilter := bson.D{{"$and", bson.A{
		bson.D{{"status", true}},
		bson.D{{"rating_type_id", ratingTypeId}},
		bson.D{{"source_uid", sourceUid}},
		bson.D{{"source_type", sourceType}},
	}}}
	err := r.db.Collection(entity.RatingsMpCol{}.CollectionName()).FindOne(ctx, bsonFilter).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingMpRepo) GetRatingTypeNumByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	var ratingTypeNum entity.RatingTypesNumCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	bsonFilter := bson.D{{"$and",
		bson.A{
			bsonStatus,
			bson.M{"_id": id},
		},
	},
	}
	err := r.db.Collection(ratingTypeNum.CollectionName()).FindOne(ctx, bsonFilter).Decode(&ratingTypeNum)
	if err != nil {

		return nil, err
	}
	return &ratingTypeNum, nil
}

func (r *ratingMpRepo) FindRatingTypeNumByRatingType(ratingType string) (*entity.RatingTypesNumCol, error) {
	var ratingTypeNum entity.RatingTypesNumCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	bsonFilter := bson.D{{"type", ratingType}}
	err := r.db.Collection(ratingTypeNum.CollectionName()).FindOne(ctx, bsonFilter).Decode(&ratingTypeNum)
	if err != nil {

		return nil, err
	}
	return &ratingTypeNum, nil
}

func (r *ratingMpRepo) FindRatingTypeNumByRatingTypeID(ratingTypeId primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	var ratingTypeNum entity.RatingTypesNumCol
	ratingTypeNumColl := r.db.Collection(ratingTypeNum.CollectionName())
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingTypeNumColl.FindOne(ctx, bson.M{"_id": ratingTypeId}).Decode(&ratingTypeNum)
	if err != nil {
		return nil, err
	}
	return &ratingTypeNum, nil
}

func getPaginationMp(db *mongo.Database, page int, limit int64, result interface{}, collectionName string, filter bson.D) *base.Pagination {
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
