package repository

import (
	"context"
	"errors"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type ratingRepo struct {
	db *mongo.Database
}

var bsonStatus = bson.D{{"status", true}}

type RatingRepository interface {
	// Rating type num
	CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) (*entity.RatingTypesNumCol, error)
	UpdateRatingTypeNum(id primitive.ObjectID, input request.CreateRatingTypeNumRequest) error
	GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error)
	DeleteRatingTypeNum(id primitive.ObjectID) error
	GetRatingTypeNums(filter request.Filter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesNumCol, *base.Pagination, error)

	// Rating submission
	CreateRatingSubmission(input request.CreateRatingSubmissonRequest) (*entity.RatingSubmisson, error)
	UpdateRatingSubmission(input request.UpdateRatingSubmissonRequest, id primitive.ObjectID) error
	DeleteSubmission(id primitive.ObjectID) error
	GetRatingSubmissionById(id primitive.ObjectID) (*entity.RatingSubmisson, error)
	GetListRatingSubmissions(filter request.RatingSubmissionFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingSubmisson, *base.Pagination, error)

	CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, error)
	GetRatingById(id primitive.ObjectID) (*entity.RatingsCol, error)
	UpdateRating(id primitive.ObjectID, input request.SaveRatingRequest) (*entity.RatingsCol, error)
	DeleteRating(id primitive.ObjectID) error
	GetRatingsByParams(limit, page, dir int, sort string, filter request.RatingFilter) ([]entity.RatingsCol, *base.Pagination, error)
	GetRatingByName(name string) (*entity.RatingsCol, error)

	CreateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) error
	GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error)
	UpdateRatingTypeLikert(id primitive.ObjectID, input request.SaveRatingTypeLikertRequest) error
	DeleteRatingTypeLikert(id primitive.ObjectID) error
	GetRatingTypeLikerts(filter request.FilterRatingTypeLikert, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesLikertCol, *base.Pagination, error)

	FindRatingSubmissionByUserIDAndRatingID(userId *string, ratingId string) (*entity.RatingSubmisson, error)
	FindRatingSubmissionByUserIDLegacyAndRatingID(userIdLegacy *string, ratingId string) (*entity.RatingSubmisson, error)
	FindRatingByRatingID(ratingId primitive.ObjectID) (*entity.RatingsCol, error)
	FindRatingNumericTypeByRatingTypeID(ratingTypeId primitive.ObjectID) (*entity.RatingTypesNumCol, error)
}

func NewRatingRepository(db *mongo.Database) RatingRepository {
	return &ratingRepo{db}
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

func (r *ratingRepo) UpdateRatingTypeNum(id primitive.ObjectID, input request.CreateRatingTypeNumRequest) error {
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
	var pagination base.Pagination
	var typeIds []primitive.ObjectID
	for _, ele := range filter.TypeId {
		objectId, _ := primitive.ObjectIDFromHex(ele)
		typeIds = append(typeIds, objectId)
	}
	bsonMinScore := bson.D{}
	bsonMaxScore := bson.D{}
	bsonTypeIdsScore := bson.D{}
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

	cursor, err := r.db.Collection("ratingTypesNumCol").
		Find(context.Background(), filter1,
			&options.FindOptions{
				Sort:  bson.D{bson.E{Key: sort, Value: dir}},
				Limit: &limit,
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
	pagination.Page = page
	pagination.Limit = int(limit)
	pagination.TotalRecords = int64(len(results))
	pagination.TotalPage = int(math.Ceil(float64(pagination.TotalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(len(results))
	return results, &pagination, nil
}

func (r *ratingRepo) CreateRatingSubmission(input request.CreateRatingSubmissonRequest) (*entity.RatingSubmisson, error) {

	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var ratingSubmission entity.RatingSubmisson

	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		result, err := ratingSubmissionColl.InsertOne(ctx, bson.M{
			"rating_id":      input.RatingID,
			"user_id":        input.UserID,
			"user_id_legacy": input.UserIDLegacy,
			"comment":        input.Comment,
			"value":          input.Value,
			"ip_address":     input.IPAddress,
			"user_agent":     input.UserAgent,
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
		ratingSubmission.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	})
	if errTransaction != nil {
		return nil, errTransaction
	}
	return &ratingSubmission, nil
}

func (r *ratingRepo) UpdateRatingSubmission(input request.UpdateRatingSubmissonRequest, id primitive.ObjectID) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	var timeUpdate time.Time
	timeUpdate = time.Now().In(util.Loc)
	ratingSubmiss := entity.RatingSubmisson{
		RatingID:     input.RatingID,
		UserID:       input.UserID,
		UserIDLegacy: input.UserIDLegacy,
		Comment:      input.Comment,
		Value:        input.Value,
		UpdatedAt:    timeUpdate,
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

func (r *ratingRepo) FindRatingSubmissionByUserIDAndRatingID(userId *string, ratingId string) (*entity.RatingSubmisson, error) {
	var ratingSubmission entity.RatingSubmisson
	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"user_id": &userId, "rating_id": ratingId}).Decode(&ratingSubmission)
	if err != nil {

		return nil, err
	}
	return &ratingSubmission, nil
}

func (r *ratingRepo) FindRatingSubmissionByUserIDLegacyAndRatingID(userIdLegacy *string, ratingId string) (*entity.RatingSubmisson, error) {
	var ratingSubmission entity.RatingSubmisson
	ratingSubmissionColl := r.db.Collection("ratingSubCol")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := ratingSubmissionColl.FindOne(ctx, bson.M{"user_id_legacy": &userIdLegacy, "rating_id": ratingId}).Decode(&ratingSubmission)
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
	update := bson.D{{"$set", bson.D{{"status", false}}}}
	_, err := r.db.Collection("ratingSubCol").UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *ratingRepo) GetListRatingSubmissions(filter request.RatingSubmissionFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingSubmisson, *base.Pagination, error) {
	var results []entity.RatingSubmisson
	var pagination base.Pagination
	bsonUserUid := bson.D{}
	bsonScore := bson.D{}
	if len(filter.UserUID) > 0 {
		bsonUserUid = bson.D{{"user_uid", bson.D{{"$in", filter.UserUID}}}}
	}
	if len(filter.Score) > 0 {
		bsonScore = bson.D{{"value", bson.D{{"$in", filter.Score}}}}
	}

	filter1 := bson.D{{"$and",

		bson.A{
			bsonStatus,
			bson.D{{"$or",
				bson.A{
					bsonUserUid,
				}}},
			bson.D{{"$or",
				bson.A{
					bsonScore,
				}}},
		},
	},
	}

	cursor, err := r.db.Collection("ratingSubCol").
		Find(context.Background(), filter1,
			&options.FindOptions{
				Sort:  bson.D{bson.E{Key: sort, Value: dir}},
				Limit: &limit,
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
	pagination.Page = page
	pagination.Limit = int(limit)
	pagination.TotalRecords = int64(len(results))
	pagination.TotalPage = int(math.Ceil(float64(pagination.TotalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(pagination.Limit*(pagination.Page-1)) + int64(len(results))
	return results, &pagination, nil
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
	bsonFilter := bson.D{{"$and",
		bson.A{
			bsonStatus,
			bson.M{"_id": id},
		},
	},
	}
	err := r.db.Collection(entity.RatingsCol{}.CollectionName()).FindOne(ctx, bsonFilter).Decode(&rating)
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepo) UpdateRating(id primitive.ObjectID, input request.SaveRatingRequest) (*entity.RatingsCol, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	rating := entity.RatingsCol{
		Name:           input.Name,
		Description:    input.Description,
		SourceUid:      input.SourceUid,
		SourceType:     input.SourceType,
		RatingType:     input.RatingType,
		RatingTypeId:   input.RatingTypeId,
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
			return err
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

	if len(filter.SourceUid) > 0 {
		bsonSourceUid = bson.D{{"source_uid", bson.D{{"$in", filter.SourceUid}}}}
	}
	if len(filter.RatingTypeId) > 0 {
		bsonRatingTypeId = bson.D{{"rating_type_id", bson.D{{"$in", filter.RatingTypeId}}}}
	}

	bsonFilter := bson.D{{"$and",

		bson.A{
			bsonStatus,
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

	crsr, err := r.db.Collection(entity.RatingsCol{}.CollectionName()).Find(ctx, filter)
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
		Statement01:   input.Statement01,
		Statement02:   input.Statement02,
		Statement03:   input.Statement03,
		Statement04:   input.Statement04,
		Statement05:   input.Statement05,
		Statement06:   input.Statement06,
		Statement07:   input.Statement07,
		Statement08:   input.Statement08,
		Statement09:   input.Statement09,
		Statement10:   input.Statement10,
		Status:        input.Status,
		UpdatedAt:     &timeUpdate,
	}
	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", ratingTypeLikert}}
	// transaction
	errTransaction := r.db.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err1 := r.db.Collection("ratingTypesLikertCol").FindOneAndUpdate(context.Background(), filter, data, &options.FindOneAndUpdateOptions{})
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
			bson.D{{"$or",
				bson.A{
					bsonTypeIdsScore,
				}}},
		},
	},
	}
	skip := int64(page)*limit - limit
	cursor, err := r.db.Collection("ratingTypesLikertCol").
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
	pagination := getPagination(page, limit, results)
	return results, pagination, nil
}

func getPagination(page int, limit int64, result interface{}) *base.Pagination {
	var pagination base.Pagination
	s := reflect.ValueOf(result)
	pagination.Page = page
	pagination.Limit = int(limit)
	pagination.TotalRecords = int64(pagination.Limit*(pagination.Page-1)) + int64(s.Len())
	pagination.TotalPage = int(math.Ceil(float64(pagination.TotalRecords) / float64(pagination.GetLimit())))
	pagination.Records = int64(s.Len())
	return &pagination
}
