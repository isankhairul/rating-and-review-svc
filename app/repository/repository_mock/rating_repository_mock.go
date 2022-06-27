package repository_mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type RatingRepositoryMock struct {
	Mock mock.Mock
}

func (_m *RatingRepositoryMock) GetRatingBySourceUidAndSourceType(sourceUid, sourceType string) (*entity.RatingsCol, error) {
	ret := _m.Mock.Called(sourceUid, sourceType)

	var r0 *entity.RatingsCol
	if rf, ok := ret.Get(0).(func(string, string) *entity.RatingsCol); ok {
		r0 = rf(sourceUid, sourceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingsCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(sourceUid, sourceType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (repository *RatingRepositoryMock) GetRatingByType(id string) (*entity.RatingsCol, error) {
	arguments := repository.Mock.Called(id)
	rating := "629dce7bf1f26275e0d84826"
	if rating == id {
		return nil, errors.New("Error")
	}
	if arguments.Get(0) == nil {
		return nil, mongo.ErrNoDocuments
	} else {
		rating := arguments.Get(0).(entity.RatingsCol)
		return &rating, nil
	}

}

func (repository *RatingRepositoryMock) GetRatingSubmissionByRatingId(id string) (*entity.RatingSubmisson, error) {
	arguments := repository.Mock.Called(id)
	submission, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	objectId, _ := primitive.ObjectIDFromHex(id)
	if submission == objectId {
		return nil, errors.New("Error")
	}
	if arguments.Get(0) == nil {
		return nil, mongo.ErrNoDocuments
	} else {
		submission := arguments.Get(0).(entity.RatingSubmisson)
		return &submission, nil
	}
}

func (repository *RatingRepositoryMock) GetRatingByTypeNum(id string) (*entity.RatingsCol, error) {
	arguments := repository.Mock.Called(id)
	rating := "629dce7bf1f26275e0d84826"
	if rating == id {
		return nil, errors.New("Error")
	}
	if arguments.Get(0) == nil {
		return nil, mongo.ErrNoDocuments
	} else {
		rating := arguments.Get(0).(entity.RatingsCol)
		return &rating, nil
	}
}

func (repository *RatingRepositoryMock) CreateRatingTypeLikert(input request.SaveRatingTypeLikertRequest) error {
	//var err error
	if input.Type == "typeErr" {
		return errors.New("Error")
	} else if input.Type == "duplicate" {
		writeError := mongo.WriteError{
			Code:    11000,
			Message: "E11000 duplicate key error collection: ratingsNReviewsDB.ratingTypesLikertCol index: type_1 dup key: { type: \"type\" }",
		}
		writeErrors := mongo.WriteErrors{}
		writeErrors = append(writeErrors, writeError)
		mongo.WriteException{
			WriteErrors: writeErrors,
		}.Error()
		return mongo.WriteException{
			WriteErrors: writeErrors,
		}
	}
	return nil
}

func (repository *RatingRepositoryMock) GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
	arguments := repository.Mock.Called(id)
	ratingType, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	if ratingType == id {
		return nil, errors.New("Error")
	}
	if arguments.Get(0) == nil {
		return nil, mongo.ErrNoDocuments
	} else {
		ratingType := arguments.Get(0).(entity.RatingTypesLikertCol)
		return &ratingType, nil
	}
}

func (repository *RatingRepositoryMock) UpdateRatingTypeLikert(id primitive.ObjectID, input request.SaveRatingTypeLikertRequest) error {
	description := "failed"
	if input.Description == &description {
		return errors.New("Error")
	}
	return nil
}

func (repository *RatingRepositoryMock) DeleteRatingTypeLikert(id primitive.ObjectID) error {
	objectId, _ := primitive.ObjectIDFromHex("629ec0836f3c2761ba2dc869")
	if id == objectId {
		return errors.New("user not found")
	}
	objectId2, _ := primitive.ObjectIDFromHex("629ec0836f3c2761ba2dc899")
	if id == objectId2 {
		return errors.New("error")
	}
	return nil
}

func (repository *RatingRepositoryMock) GetRatingTypeLikerts(filter request.FilterRatingTypeLikert, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesLikertCol, *base.Pagination, error) {
	arguments := repository.Mock.Called(filter, page, limit, sort, dir)
	rating := entity.RatingTypesLikertCol{}
	if arguments.Get(0) == rating {
		return nil, nil, gorm.ErrRecordNotFound
	}
	if sort == "failed" {
		return nil, nil, errors.New("Errors")
	}
	return arguments.Get(0).([]entity.RatingTypesLikertCol), arguments.Get(1).(*base.Pagination), nil
}

var (
	id               primitive.ObjectID
	matchStrValuePtr = "match"
)

func (repository *RatingRepositoryMock) GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	arguments := repository.Mock.Called(id)
	ratingType, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	if ratingType == id {
		return nil, errors.New("Error")
	}
	if arguments.Get(0) == nil {
		return nil, mongo.ErrNoDocuments
	} else {
		ratingType := arguments.Get(0).(entity.RatingTypesNumCol)
		return &ratingType, nil
	}
}

func (repository *RatingRepositoryMock) DeleteRatingTypeNum(id primitive.ObjectID) error {
	objectId, _ := primitive.ObjectIDFromHex("629ec0836f3c2761ba2dc869")
	if id == objectId {
		return errors.New("user not found")
	}
	objectId2, _ := primitive.ObjectIDFromHex("629ec0836f3c2761ba2dc899")
	if id == objectId2 {
		return errors.New("error")
	}
	return nil
}

func (repository *RatingRepositoryMock) GetRatingTypeNums(filter request.Filter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingTypesNumCol, *base.Pagination, error) {
	arguments := repository.Mock.Called(filter, page, limit, sort, dir)
	rating := entity.RatingTypesNumCol{}
	if arguments.Get(0) == rating {
		return nil, nil, gorm.ErrRecordNotFound
	}
	if sort == "failed" {
		return nil, nil, errors.New("Errors")
	}
	return arguments.Get(0).([]entity.RatingTypesNumCol), arguments.Get(1).(*base.Pagination), nil
}

func (repository *RatingRepositoryMock) UpdateRatingTypeNum(id primitive.ObjectID, input request.EditRatingTypeNumRequest) error {
	objectId, _ := primitive.ObjectIDFromHex("629ec07e6f3c2761ba2dc868")
	if id == objectId {
		return errors.New("Error")
	}
	return nil
}

func (repository *RatingRepositoryMock) CreateRatingTypeNum(input request.CreateRatingTypeNumRequest) (*entity.RatingTypesNumCol, error) {
	ratingTypesNumCol := entity.RatingTypesNumCol{}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingTypesNumCol.ID = objectId
	if input.Type == "12345" {
		return nil, errors.New("errors")
	}
	return &ratingTypesNumCol, nil
}

func (repository *RatingRepositoryMock) Paginate(value interface{}, pagination *base.Pagination, db *gorm.DB, currRecord int64) func(db *gorm.DB) *gorm.DB {
	return nil
}

func (repository *RatingRepositoryMock) CreateRatingSubmission(input []request.SaveRatingSubmission) (*[]entity.RatingSubmisson, error) {
	for _, arg := range input {
		if *arg.UserID != "629dce7bf1f26275e0d84826" {
			return nil, errors.New("can not be created")
		}
	}
	return nil, nil
}
func (repository *RatingRepositoryMock) UpdateRatingSubmission(input request.UpdateRatingSubmissonRequest, id primitive.ObjectID) error {
	ratingSubmissionCol := entity.RatingSubmisson{}
	ratingSubmissionCol.ID = id
	return nil
}
func (repository *RatingRepositoryMock) DeleteSubmission(id primitive.ObjectID) error {
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	if id != objectId {
		return errors.New("user not found")
	}
	return nil
}
func (repository *RatingRepositoryMock) GetRatingSubmissionById(id primitive.ObjectID) (*entity.RatingSubmisson, error) {
	arguments := repository.Mock.Called(id)
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	if id != objectId {
		return nil, mongo.ErrNoDocuments
	}
	sub := arguments.Get(0).(entity.RatingSubmisson)
	return &sub, nil
}
func (repository *RatingRepositoryMock) FindRatingSubmissionByUserIDAndRatingID(userId *string, ratingId string) (*entity.RatingSubmisson, error) {
	arguments := repository.Mock.Called(userId, ratingId)
	if *userId == "629dce7bf1f26275e0d84826" && ratingId == "629dce7bf1f26275e0d84826" {
		return nil, errors.New("record found")
	}
	sub := arguments.Get(0).(entity.RatingSubmisson)
	return &sub, nil
}
func (repository *RatingRepositoryMock) FindRatingSubmissionByUserIDLegacyAndRatingID(userIdLegacy *string, ratingId string) (*entity.RatingSubmisson, error) {
	if *userIdLegacy == "629dce7bf1f26275e0d84826" && ratingId == "629dce7bf1f26275e0d84826" {
		return nil, gorm.ErrRecordNotFound
	}
	sub := entity.RatingSubmisson{}
	return &sub, nil
}
func (repository *RatingRepositoryMock) FindRatingByRatingID(ratingId primitive.ObjectID) (*entity.RatingsCol, error) {
	arguments := repository.Mock.Called(ratingId)
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	if ratingId != objectId {
		return nil, errors.New("user not found")
	}
	rat := arguments.Get(0).(entity.RatingsCol)
	return &rat, nil
}
func (repository *RatingRepositoryMock) FindRatingNumericTypeByRatingTypeID(ratingTypeId primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	arguments := repository.Mock.Called(ratingTypeId)
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	if ratingTypeId != objectId {
		return nil, errors.New("user not found")
	}
	num := arguments.Get(0).(entity.RatingTypesNumCol)
	return &num, nil
}

func (repository *RatingRepositoryMock) GetListRatingSubmissions(filter request.RatingSubmissionFilter, page int, limit int64, sort string, dir interface{}) ([]entity.RatingSubmisson, *base.Pagination, error) {
	arguments := repository.Mock.Called(filter, page, limit, sort, dir)
	if sort == "no data " {
		return nil, nil, mongo.ErrNoDocuments
	}
	return arguments.Get(0).([]entity.RatingSubmisson), arguments.Get(1).(*base.Pagination), nil
}

// CreateRating provides a mock function with given fields: input
func (_m *RatingRepositoryMock) CreateRating(input request.SaveRatingRequest) (*entity.RatingsCol, error) {
	ret := _m.Mock.Called(input)

	var r0 *entity.RatingsCol
	if rf, ok := ret.Get(0).(func(request.SaveRatingRequest) *entity.RatingsCol); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingsCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(request.SaveRatingRequest) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRating provides a mock function with given fields: id
func (_m *RatingRepositoryMock) DeleteRating(id primitive.ObjectID) error {
	ret := _m.Mock.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRatingById provides a mock function with given fields: id
func (_m *RatingRepositoryMock) GetRatingById(id primitive.ObjectID) (*entity.RatingsCol, error) {
	ret := _m.Mock.Called(id)

	var r0 *entity.RatingsCol
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) *entity.RatingsCol); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingsCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingByName provides a mock function with given fields: name
func (_m *RatingRepositoryMock) GetRatingByName(name string) (*entity.RatingsCol, error) {
	ret := _m.Mock.Called(name)

	var r0 *entity.RatingsCol
	if rf, ok := ret.Get(0).(func(string) *entity.RatingsCol); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingsCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingTypeLikertByIdAndStatus provides a mock function with given fields: id
func (_m *RatingRepositoryMock) GetRatingTypeLikertByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
	ret := _m.Mock.Called(id)
	objectId2, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84820")
	if id == objectId2 {
		num := ret.Get(0).(entity.RatingTypesLikertCol)
		return &num, nil
	}
	var r0 *entity.RatingTypesLikertCol
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) *entity.RatingTypesLikertCol); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingTypesLikertCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingTypeNumByIdAndStatus provides a mock function with given fields: id
func (_m *RatingRepositoryMock) GetRatingTypeNumByIdAndStatus(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	ret := _m.Mock.Called(id)

	var r0 *entity.RatingTypesNumCol
	if rf, ok := ret.Get(0).(func(primitive.ObjectID) *entity.RatingTypesNumCol); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingTypesNumCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingsByParams provides a mock function with given fields: limit, page, dir, sort, filter
func (_m *RatingRepositoryMock) GetRatingsByParams(limit int, page int, dir int, sort string, filter request.RatingFilter) ([]entity.RatingsCol, *base.Pagination, error) {
	ret := _m.Mock.Called(filter, page, limit, sort, dir)

	var r0 []entity.RatingsCol
	if rf, ok := ret.Get(0).(func(int, int, int, string, request.RatingFilter) []entity.RatingsCol); ok {
		r0 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingsCol)
		}
	}

	var r1 *base.Pagination
	if rf, ok := ret.Get(1).(func(int, int, int, string, request.RatingFilter) *base.Pagination); ok {
		r1 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*base.Pagination)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, int, int, string, request.RatingFilter) error); ok {
		r2 = rf(limit, page, dir, sort, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateRating provides a mock function with given fields: id, input
func (_m *RatingRepositoryMock) UpdateRating(id primitive.ObjectID, input request.BodyUpdateRatingRequest) (*entity.RatingsCol, error) {
	ret := _m.Mock.Called(id, input)

	var r0 *entity.RatingsCol
	if rf, ok := ret.Get(0).(func(primitive.ObjectID, request.BodyUpdateRatingRequest) *entity.RatingsCol); ok {
		r0 = rf(id, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingsCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(primitive.ObjectID, request.BodyUpdateRatingRequest) error); ok {
		r1 = rf(id, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
