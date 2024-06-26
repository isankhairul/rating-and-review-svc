package repository

import (
	"context"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/pkg/util"
	"math"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gorm.io/gorm"
)

type ratingRepo struct {
	db *mongo.Database
}

var bsonStatus = bson.D{{"status", true}}

type RatingRepository interface {
	// Rating type num
	CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) (*entity.RatingTypesNumCol, error)
	UpdateRatingTypeNum(id primitive.ObjectID, input request.EditRatingTypeNumRequest) error
	GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error)
	DeleteRatingTypeNum(id primitive.ObjectID) error
	GetRatingTypeNums(filter request.Filter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesNumCol, *base.Pagination, error)

	// Rating submission
	CreateRatingSubmission(input []request.SaveRatingSubmission) (*[]entity.RatingSubmisson, error)
	UpdateRatingSubmission(input request.UpdateRatingSubmissionRequest, id primitive.ObjectID) error
	DeleteSubmission(id primitive.ObjectID) error
	GetRatingSubmissionById(id primitive.ObjectID) (*entity.RatingSubmisson, error)
	CancelRatingSubmissionByIds(ids []primitive.ObjectID, reason string) error

	GetListRatingSubmissions(filter request.RatingSubmissionFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingSubmisson, *base.Pagination, error)
	GetRatingSubmissionByRatingId(id string) (*entity.RatingSubmisson, error)
	FindRatingSubmissionByUserIDAndRatingID(userId *string, ratingId string, sourceTransId string) (*entity.RatingSubmisson, error)
	FindRatingSubmissionByUserIDLegacyAndRatingID(userIdLegacy *string, ratingId string, sourceTransId string) (*entity.RatingSubmisson, error)
	FindRatingByRatingID(ratingId primitive.ObjectID) (*entity.RatingsCol, error)
	FindRatingBySourceUIDAndRatingType(sourceUID, ratingType string) (*entity.RatingsCol, error)
	FindRatingNumericTypeByRatingTypeID(ratingTypeId primitive.ObjectID) (*entity.RatingTypesNumCol, error)

	CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, error)
	GetRatingById(id primitive.ObjectID) (*entity.RatingsCol, error)
	UpdateRating(id primitive.ObjectID, input request.BodyUpdateRatingRequest) (*entity.RatingsCol, error)
	DeleteRating(id primitive.ObjectID) error
	GetRatingsByParams(limit, page, dir int, sort string, filter request.RatingFilter) ([]entity.RatingsCol, *base.Pagination, error)
	GetRatingByName(name string) (*entity.RatingsCol, error)
	GetRatingTypeNumByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesNumCol, error)
	GetRatingTypeLikertByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error)
	GetRatingByType(id string) (*entity.RatingsCol, error)
	GetRatingByRatingTypeSourceUidAndSourceType(ratingTypeId, sourceUid, sourceType string) (*entity.RatingsCol, error)

	CreateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) error
	GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error)
	UpdateRatingTypeLikert(id primitive.ObjectID, input request.SaveRatingTypeLikertRequest) error
	DeleteRatingTypeLikert(id primitive.ObjectID) error
	GetRatingTypeLikerts(filter request.FilterRatingTypeLikert, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesLikertCol, *base.Pagination, error)
	Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB

	CreateRatingFormula(input request.SaveRatingFormula) (*entity.RatingFormulaCol, error)
	UpdateRatingFormula(id primitive.ObjectID, input request.SaveRatingFormula) error
	GetRatingFormulaById(id primitive.ObjectID) (*entity.RatingFormulaCol, error)
	DeleteRatingFormula(id primitive.ObjectID) error
	GetRatingFormulas(filter request.RatingFormulaFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingFormulaCol, *base.Pagination, error)
}

func NewRatingRepository(db *mongo.Database) RatingRepository {
	return &ratingRepo{db}
}

func getPagination(r *ratingRepo, page int, limit int64, result interface{}, collectionName string, filter bson.D) *base.Pagination {
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

func (r *ratingRepo) CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) (*entity.RatingTypesNumCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var ratingTypeNum entity.RatingTypesNumCol

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := r.db.Collection("ratingTypesNumCol").InsertOne(ctx, bson.M{
			"type":        input.Type,
			"description": input.Description,
			"min_score":   input.MinScore,
			"max_score":   input.MaxScore,
			"scale":       input.Scale,
			"intervals":   input.Intervals,
			"status":      input.Status,
			"created_at":  time.Now().In(util.Loc),
			"updated_at":  time.Now().In(util.Loc),
		})

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		ratingTypeNum.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	})
	if errTransaction != nil {
		return nil, errTransaction
	}
	return &ratingTypeNum, nil
}

func (r *ratingRepo) UpdateRatingTypeNum(id primitive.ObjectID, input request.EditRatingTypeNumRequest) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var timeUpdate time.Time
	timeUpdate = time.Now().In(util.Loc)
	ratingTypeLikert := entity.RatingTypesNumCol{
		Type:        input.Type,
		Description: input.Description,
		MinScore:    input.MinScore,
		MaxScore:    input.MaxScore,
		Scale:       input.Scale,
		Intervals:   input.Intervals,
		Status:      input.Status,
		UpdatedAt:   timeUpdate,
	}
	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", ratingTypeLikert}}
	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err1 := r.db.Collection("ratingTypesNumCol").FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
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

func (r *ratingRepo) GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	var ratingTypeNum entity.RatingTypesNumCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection("ratingTypesNumCol").FindOne(ctx, bson.M{"_id": id}).Decode(&ratingTypeNum)
	if err != nil {

		return nil, err
	}
	return &ratingTypeNum, nil
}

func (r *ratingRepo) DeleteRatingTypeNum(id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	update := bson.D{{"$set", bson.D{{"status", false}}}}
	_, err := r.db.Collection("ratingTypesNumCol").UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ratingRepo) GetRatingTypeNums(filter request.Filter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesNumCol, *base.Pagination, error) {
	var results []entity.RatingTypesNumCol
	var typeIds []primitive.ObjectID
	for _, ele := range filter.TypeId {
		objectId, _ := primitive.ObjectIDFromHex(ele)
		typeIds = append(typeIds, objectId)
	}
	bsonMinScore := bson.D{}
	bsonMaxScore := bson.D{}
	bsonTypeIdsScore := bson.D{}
	bsonStatus := bson.D{{"status", true}}

	if len(typeIds) > 0 {
		bsonTypeIdsScore = bson.D{{"_id", bson.D{{"$in", typeIds}}}}
	}
	if len(filter.MinScore) > 0 {
		bsonMinScore = bson.D{{"min_score", bson.D{{"$in", filter.MinScore}}}}
	}
	if len(filter.MaxScore) > 0 {
		bsonMaxScore = bson.D{{"max_score", bson.D{{"$in", filter.MaxScore}}}}
	}

	filter1 := bson.D{{"$and",
		bson.A{
			bsonStatus,
			bson.D{{"$or",
				bson.A{
					bsonMinScore,
				}}},
			bson.D{{"$or",
				bson.A{
					bsonMaxScore,
				}}},
			bson.D{{"$or",
				bson.A{
					bsonTypeIdsScore,
				}}},
		},
	},
	}

	collectionName := "ratingTypesNumCol"
	skip := int64(page)*limit - limit
	cursor, err := r.db.Collection(collectionName).
		Find(context.Background(), filter1,
			&options.FindOptions{
				Sort:  bson.D{bson.E{Key: sort, Value: dir}},
				Limit: &limit,
				Skip:  &skip,
			})
	if err != nil {
		return nil, nil, err
	}
	//
	if err = cursor.All(context.TODO(), &results); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}
	pagination := getPagination(r, page, limit, results, collectionName, filter1)
	return results, pagination, nil
}

func (r *ratingRepo) CreateRatingSubmission(input []request.SaveRatingSubmission) (*[]entity.RatingSubmisson, error) {
	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var ratingSubmission []entity.RatingSubmisson
	var docs []interface{}
	for _, args := range input {
		dateNow := time.Now().In(util.Loc)
		if args.Tagging.RatingId != "" {
			docs = append(docs, bson.D{
				{Key: "rating_id", Value: args.RatingID},
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
				{Key: "tagging", Value: args.Tagging},
				{Key: "source_type", Value: args.SourceType},
			})
		} else {
			docs = append(docs, bson.D{
				{Key: "rating_id", Value: args.RatingID},
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
			})
		}
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
			data := entity.RatingSubmisson{
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

func (r *ratingRepo) UpdateRatingSubmission(input request.UpdateRatingSubmissionRequest, id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var timeUpdate time.Time
	timeUpdate = time.Now().In(util.Loc)
	ratingSubmiss := entity.RatingSubmisson{
		RatingID:  input.RatingID,
		Comment:   &input.Comment,
		Value:     *input.Value,
		UpdatedAt: timeUpdate,
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

func (r *ratingRepo) GetRatingSubmissionById(id primitive.ObjectID) (*entity.RatingSubmisson, error) {
	var ratingSubmission entity.RatingSubmisson
	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"_id": id}).Decode(&ratingSubmission)
	if err != nil {

		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingRepo) CancelRatingSubmissionByIds(ids []primitive.ObjectID, reason string) error {
	timeUpdate := time.Now().In(util.Loc)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)

	ratingSub := entity.RatingSubmisson{
		Cancelled:       true,
		CancelledReason: reason,
		UpdatedAt:       timeUpdate,
	}

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}
	data := bson.D{{Key: "$set", Value: ratingSub}}

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		_, errUpd := r.db.Collection("ratingSubCol").UpdateMany(context.Background(), filter, data)
		if errUpd != nil {
			sessionContext.AbortTransaction(sessionContext)
			return errUpd
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

func (r *ratingRepo) GetRatingSubmissionByRatingId(id string) (*entity.RatingSubmisson, error) {
	var ratingSubmission entity.RatingSubmisson
	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"rating_id": id}).Decode(&ratingSubmission)
	if err != nil {
		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingRepo) FindRatingSubmissionByUserIDAndRatingID(userId *string, ratingId string, sourceTransId string) (*entity.RatingSubmisson, error) {
	var ratingSubmission entity.RatingSubmisson
	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"user_id": &userId, "rating_id": ratingId, "source_trans_id": sourceTransId}).Decode(&ratingSubmission)
	if err != nil {

		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingRepo) FindRatingSubmissionByUserIDLegacyAndRatingID(userIdLegacy *string, ratingId string, sourceTransId string) (*entity.RatingSubmisson, error) {
	var ratingSubmission entity.RatingSubmisson
	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"user_id_legacy": &userIdLegacy, "rating_id": ratingId, "source_trans_id": sourceTransId}).Decode(&ratingSubmission)
	if err != nil {

		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingRepo) FindRatingByRatingID(ratingId primitive.ObjectID) (*entity.RatingsCol, error) {
	var rating entity.RatingsCol
	ratingColl := r.db.Collection("ratingsCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingColl.FindOne(ctx, bson.M{"_id": ratingId}).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepo) FindRatingBySourceUIDAndRatingType(sourceUID, ratingType string) (*entity.RatingsCol, error) {
	var rating entity.RatingsCol

	ratingColl := r.db.Collection("ratingsCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingColl.FindOne(ctx, bson.M{"source_uid": sourceUID, "rating_type": ratingType}).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepo) FindRatingNumericTypeByRatingTypeID(ratingTypeId primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	var ratingTypeNum entity.RatingTypesNumCol
	ratingTypeNumColl := r.db.Collection("ratingTypesNumCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingTypeNumColl.FindOne(ctx, bson.M{"_id": ratingTypeId}).Decode(&ratingTypeNum)
	if err != nil {
		return nil, err
	}
	return &ratingTypeNum, nil
}

func (r *ratingRepo) DeleteSubmission(id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	result, err := r.db.Collection("ratingSubCol").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *ratingRepo) CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var rating entity.RatingsCol

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).InsertOne(ctx, bson.M{
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

func (r *ratingRepo) GetRatingById(id primitive.ObjectID) (*entity.RatingsCol, error) {
	var rating entity.RatingsCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection(entity.RatingsCol{}.CollectionName()).FindOne(ctx, bson.M{"_id": id}).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepo) UpdateRating(id primitive.ObjectID, input request.BodyUpdateRatingRequest) (*entity.RatingsCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	rating := entity.RatingsCol{
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
		err1 := r.db.Collection(entity.RatingsCol{}.CollectionName()).FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
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

func (r *ratingRepo) DeleteRating(id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	update := bson.D{{"$set", bson.D{{"status", false}}}}
	_, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ratingRepo) GetRatingsByParams(limit, page, dir int, sort string, filter request.RatingFilter) ([]entity.RatingsCol, *base.Pagination, error) {
	var results []entity.RatingsCol
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
	cursor, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).
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

	crsr, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).Find(ctx, bsonFilter)
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

func (r *ratingRepo) GetRatingByName(name string) (*entity.RatingsCol, error) {
	var rating entity.RatingsCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection(entity.RatingsCol{}.CollectionName()).FindOne(ctx, bson.M{"name": name}).Decode(&rating)
	if err != nil {

		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepo) GetRatingByType(id string) (*entity.RatingsCol, error) {
	var rating entity.RatingsCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection(entity.RatingsCol{}.CollectionName()).FindOne(ctx, bson.M{"rating_type_id": id}).Decode(&rating)
	if err != nil {

		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepo) GetRatingTypeNumByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	var ratingTypeNum entity.RatingTypesNumCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	bsonFilter := bson.D{{"$and",
		bson.A{
			bsonStatus,
			bson.M{"_id": id},
		},
	},
	}
	err := r.db.Collection("ratingTypesNumCol").FindOne(ctx, bsonFilter).Decode(&ratingTypeNum)
	if err != nil {

		return nil, err
	}
	return &ratingTypeNum, nil
}

func (r *ratingRepo) GetRatingTypeLikertByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
	var ratingTypeLikert entity.RatingTypesLikertCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	bsonFilter := bson.D{{Key: "$and",
		Value: bson.A{
			bsonStatus,
			bson.M{"_id": id},
		},
	},
	}
	err := r.db.Collection("ratingTypesLikertCol").FindOne(ctx, bsonFilter).Decode(&ratingTypeLikert)
	if err != nil {
		return nil, err
	}
	return &ratingTypeLikert, nil
}

func (r *ratingRepo) GetRatingByRatingTypeSourceUidAndSourceType(ratingTypeId, sourceUid, sourceType string) (*entity.RatingsCol, error) {
	var rating entity.RatingsCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	bsonFilter := bson.D{{"$and", bson.A{
		bson.D{{"status", true}},
		bson.D{{"rating_type_id", ratingTypeId}},
		bson.D{{"source_uid", sourceUid}},
		bson.D{{"source_type", sourceType}},
	}}}
	err := r.db.Collection(entity.RatingsCol{}.CollectionName()).FindOne(ctx, bsonFilter).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepo) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(value).Count(&totalRecords)

	pagination.TotalRecords = totalRecords
	pagination.TotalPage = int(math.Ceil(float64(totalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(currRecord)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func (r *ratingRepo) GetListRatingSubmissions(filter request.RatingSubmissionFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingSubmisson, *base.Pagination, error) {
	var results []entity.RatingSubmisson
	startDate, _ := time.Parse(util.LayoutDateOnly, filter.StartDate)
	endDate, _ := time.Parse(util.LayoutDateOnly, filter.EndDate)
	if errD := endDate.Before(startDate); errD == true {
		return nil, nil, errors.New("end_date can not before start_date")
	}
	bsonUserUid := bson.D{}
	bsonRating := bson.D{}
	bsonDate := bson.D{}
	bsonTransId := bson.D{}
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
		},
	},
	}

	limitPage := int(limit)
	cursor, err := r.db.Collection("ratingSubCol").
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
	pagination := getPagination(r, page, limit, results, "ratingSubCol", filter1)
	return results, pagination, nil
}

func (r *ratingRepo) CreateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var ratingTypeLikert entity.RatingTypesLikertCol

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := r.db.Collection("ratingTypesLikertCol").InsertOne(ctx, bson.M{
			"type":           input.Type,
			"description":    input.Description,
			"num_statements": input.NumStatements,
			"statement_01":   input.Statement01,
			"statement_02":   input.Statement02,
			"statement_03":   input.Statement03,
			"statement_04":   input.Statement04,
			"statement_05":   input.Statement05,
			"statement_06":   input.Statement06,
			"statement_07":   input.Statement07,
			"statement_08":   input.Statement08,
			"statement_09":   input.Statement09,
			"statement_10":   input.Statement10,
			"status":         input.Status,
			"created_at":     time.Now().In(util.Loc),
			"updated_at":     time.Now().In(util.Loc),
		})

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		ratingTypeLikert.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	})
	if errTransaction != nil {
		return errTransaction
	}
	return nil
}

func (r *ratingRepo) GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
	var ratingTypeLikert entity.RatingTypesLikertCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection("ratingTypesLikertCol").FindOne(ctx, bson.M{"_id": id}).Decode(&ratingTypeLikert)
	if err != nil {

		return nil, err
	}
	return &ratingTypeLikert, nil
}

func (r *ratingRepo) UpdateRatingTypeLikert(id primitive.ObjectID, input request.SaveRatingTypeLikertRequest) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var timeUpdate time.Time
	timeUpdate = time.Now().In(util.Loc)
	ratingTypeLikert := entity.RatingTypesLikertCol{
		Type:          input.Type,
		Description:   input.Description,
		NumStatements: input.NumStatements,
		Status:        input.Status,
		UpdatedAt:     &timeUpdate,
	}
	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", ratingTypeLikert}}
	dataStatement := bson.D{{"$set", bson.D{
		{"statement_01", input.Statement01},
		{"statement_02", input.Statement02},
		{"statement_03", input.Statement03},
		{"statement_04", input.Statement04},
		{"statement_05", input.Statement05},
		{"statement_06", input.Statement06},
		{"statement_07", input.Statement07},
		{"statement_08", input.Statement08},
		{"statement_09", input.Statement09},
		{"statement_10", input.Statement10},
	}}}
	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		_, err1 := r.db.Collection("ratingTypesLikertCol").UpdateOne(context.Background(), filter, dataStatement, &options.UpdateOptions{})
		if err1 != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err1
		}

		err2 := r.db.Collection("ratingTypesLikertCol").FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
		if err2 != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err2.Err()
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

func (r *ratingRepo) DeleteRatingTypeLikert(id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	update := bson.D{{"$set", bson.D{{"status", false}}}}
	_, err := r.db.Collection("ratingTypesLikertCol").UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ratingRepo) GetRatingTypeLikerts(filter request.FilterRatingTypeLikert, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesLikertCol, *base.Pagination, error) {
	var results []entity.RatingTypesLikertCol
	var typeIds []primitive.ObjectID
	for _, ele := range filter.TypeId {
		objectId, _ := primitive.ObjectIDFromHex(ele)
		typeIds = append(typeIds, objectId)
	}
	bsonTypeIdsScore := bson.D{}
	if len(typeIds) > 0 {
		bsonTypeIdsScore = bson.D{{"_id", bson.D{{"$in", typeIds}}}}
	}

	filter1 := bson.D{{"$and",
		bson.A{
			bsonStatus,
			bson.D{{"$or",
				bson.A{
					bsonTypeIdsScore,
				}}},
		},
	},
	}
	collectionName := "ratingTypesLikertCol"
	skip := int64(page)*limit - limit
	cursor, err := r.db.Collection(collectionName).
		Find(context.Background(), filter1,
			&options.FindOptions{
				Sort:  bson.D{bson.E{Key: sort, Value: dir}},
				Limit: &limit,
				Skip:  &skip,
			})
	if err != nil {
		return nil, nil, err
	}
	//
	if err = cursor.All(context.TODO(), &results); err != nil {
		if err != nil {
			return nil, nil, err
		}
	}
	pagination := getPagination(r, page, limit, results, collectionName, filter1)
	return results, pagination, nil
}

func (r *ratingRepo) CreateRatingFormula(input request.SaveRatingFormula) (*entity.RatingFormulaCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var ratingFormula entity.RatingFormulaCol

	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := r.db.Collection("ratingFormulaCol").InsertOne(ctx, bson.M{
			"source_type":    input.SourceType,
			"formula":        input.Formula,
			"rating_type_id": input.RatingTypeId,
			"rating_type":    input.RatingType,
			"status":         input.Status,
			"created_at":     time.Now().In(util.Loc),
			"updated_at":     time.Now().In(util.Loc),
		})

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		ratingFormula.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	})
	if errTransaction != nil {
		return nil, errTransaction
	}
	return &ratingFormula, nil
}

func (r *ratingRepo) UpdateRatingFormula(id primitive.ObjectID, input request.SaveRatingFormula) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var timeUpdate time.Time
	timeUpdate = time.Now().In(util.Loc)
	ratingFormula := entity.RatingFormulaCol{
		SourceType:   input.SourceType,
		Formula:      input.Formula,
		RatingTypeId: input.RatingTypeId,
		RatingType:   input.RatingType,
		Status:       input.Status,
		UpdatedAt:    timeUpdate,
	}
	filter := bson.D{{Key: "_id", Value: id}}
	data := bson.D{{Key: "$set", Value: ratingFormula}}
	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err1 := r.db.Collection("ratingFormulaCol").FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
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

func (r *ratingRepo) GetRatingFormulaById(id primitive.ObjectID) (*entity.RatingFormulaCol, error) {
	var ratingFormula entity.RatingFormulaCol
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := r.db.Collection("ratingFormulaCol").FindOne(ctx, bson.M{"_id": id}).Decode(&ratingFormula)
	if err != nil {

		return nil, err
	}
	return &ratingFormula, nil
}

func (r *ratingRepo) DeleteRatingFormula(id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: false}}}}
	_, err := r.db.Collection("ratingFormulaCol").UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ratingRepo) GetRatingFormulas(filter request.RatingFormulaFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingFormulaCol, *base.Pagination, error) {
	var results []entity.RatingFormulaCol
	var typeIds []primitive.ObjectID
	for _, ele := range filter.TypeId {
		objectId, _ := primitive.ObjectIDFromHex(ele)
		typeIds = append(typeIds, objectId)
	}
	bsonSourceType := bson.D{}
	bsonRatingTypeIds := bson.D{}
	bsonStatus := bson.D{{Key: "status", Value: true}}

	if len(typeIds) > 0 {
		bsonRatingTypeIds = bson.D{{Key: "_id", Value: bson.D{{"$in", typeIds}}}}
	}
	if len(filter.SourceType) > 0 {
		bsonSourceType = bson.D{{Key: "min_score", Value: bson.D{{Key: "$in", Value: filter.SourceType}}}}
	}

	filter1 := bson.D{{Key: "$and",
		Value: bson.A{
			bsonStatus,
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonSourceType,
				}}},
			bson.D{{Key: "$or",
				Value: bson.A{
					bsonRatingTypeIds,
				}}},
		},
	},
	}

	collectionName := "ratingFormulaCol"
	skip := int64(page)*limit - limit
	cursor, err := r.db.Collection(collectionName).
		Find(context.Background(), filter1,
			&options.FindOptions{
				Sort:  bson.D{bson.E{Key: sort, Value: dir}},
				Limit: &limit,
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
	pagination := getPagination(r, page, limit, results, collectionName, filter1)
	return results, pagination, nil
}
