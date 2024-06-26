// Code generated by mockery v2.13.1. DO NOT EDIT.

package public_repository_mock

import (
	base "go-klikdokter/app/model/base"
	entity "go-klikdokter/app/model/entity"

	mock "github.com/stretchr/testify/mock"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	publicrequest "go-klikdokter/app/model/request/public"

	request "go-klikdokter/app/model/request"
)

// PublicRatingRepository is an autogenerated mock type for the PublicRatingRepository type
type PublicRatingRepository struct {
	mock.Mock
}

// CountRatingSubsByRatingIdAndValue provides a mock function with given fields: ratingId, value
func (_m *PublicRatingRepository) CountRatingSubsByRatingIdAndValue(ratingId string, value string) (int64, error) {
	ret := _m.Called(ratingId, value)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, string) int64); ok {
		r0 = rf(ratingId, value)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(ratingId, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRatingSubHelpful provides a mock function with given fields: input
func (_m *PublicRatingRepository) CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) (*entity.RatingSubHelpfulCol, error) {
	ret := _m.Called(input)

	var r0 *entity.RatingSubHelpfulCol
	if rf, ok := ret.Get(0).(func(request.CreateRatingSubHelpfulRequest) *entity.RatingSubHelpfulCol); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingSubHelpfulCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(request.CreateRatingSubHelpfulRequest) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetListRatingBySourceTypeAndUID provides a mock function with given fields: sourceType, sourceUID
func (_m *PublicRatingRepository) GetListRatingBySourceTypeAndUID(sourceType string, sourceUID string) ([]entity.RatingsCol, error) {
	ret := _m.Called(sourceType, sourceUID)

	var r0 []entity.RatingsCol
	if rf, ok := ret.Get(0).(func(string, string) []entity.RatingsCol); ok {
		r0 = rf(sourceType, sourceUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingsCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(sourceType, sourceUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPublicRatingSubmissions provides a mock function with given fields: limit, page, dir, sort, filter
func (_m *PublicRatingRepository) GetPublicRatingSubmissions(limit int, page int, dir int, sort string, filter publicrequest.FilterRatingSubmission) ([]entity.RatingSubmisson, *base.Pagination, error) {
	ret := _m.Called(limit, page, dir, sort, filter)

	var r0 []entity.RatingSubmisson
	if rf, ok := ret.Get(0).(func(int, int, int, string, publicrequest.FilterRatingSubmission) []entity.RatingSubmisson); ok {
		r0 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingSubmisson)
		}
	}

	var r1 *base.Pagination
	if rf, ok := ret.Get(1).(func(int, int, int, string, publicrequest.FilterRatingSubmission) *base.Pagination); ok {
		r1 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*base.Pagination)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, int, int, string, publicrequest.FilterRatingSubmission) error); ok {
		r2 = rf(limit, page, dir, sort, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPublicRatingsByParams provides a mock function with given fields: limit, page, dir, sort, filter
func (_m *PublicRatingRepository) GetPublicRatingsByParams(limit int, page int, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsCol, *base.Pagination, error) {
	ret := _m.Called(limit, page, dir, sort, filter)

	var r0 []entity.RatingsCol
	if rf, ok := ret.Get(0).(func(int, int, int, string, publicrequest.FilterRatingSummary) []entity.RatingsCol); ok {
		r0 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingsCol)
		}
	}

	var r1 *base.Pagination
	if rf, ok := ret.Get(1).(func(int, int, int, string, publicrequest.FilterRatingSummary) *base.Pagination); ok {
		r1 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*base.Pagination)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, int, int, string, publicrequest.FilterRatingSummary) error); ok {
		r2 = rf(limit, page, dir, sort, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetRatingFormulaByRatingTypeIdAndSourceType provides a mock function with given fields: ratingTypeId, sourceType
func (_m *PublicRatingRepository) GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId string, sourceType string) (*entity.RatingFormulaCol, error) {
	ret := _m.Called(ratingTypeId, sourceType)

	var r0 *entity.RatingFormulaCol
	if rf, ok := ret.Get(0).(func(string, string) *entity.RatingFormulaCol); ok {
		r0 = rf(ratingTypeId, sourceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingFormulaCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(ratingTypeId, sourceType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingSubHelpfulByRatingSubAndActor provides a mock function with given fields: ratingSubId, userIdLegacy
func (_m *PublicRatingRepository) GetRatingSubHelpfulByRatingSubAndActor(ratingSubId string, userIdLegacy string) (*entity.RatingSubHelpfulCol, error) {
	ret := _m.Called(ratingSubId, userIdLegacy)

	var r0 *entity.RatingSubHelpfulCol
	if rf, ok := ret.Get(0).(func(string, string) *entity.RatingSubHelpfulCol); ok {
		r0 = rf(ratingSubId, userIdLegacy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingSubHelpfulCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(ratingSubId, userIdLegacy)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingSubsByRatingId provides a mock function with given fields: ratingId
func (_m *PublicRatingRepository) GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmisson, error) {
	ret := _m.Called(ratingId)

	var r0 []entity.RatingSubmisson
	if rf, ok := ret.Get(0).(func(string) []entity.RatingSubmisson); ok {
		r0 = rf(ratingId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingSubmisson)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ratingId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingTypeLikertById provides a mock function with given fields: id
func (_m *PublicRatingRepository) GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
	ret := _m.Called(id)

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

// GetRatingTypeNumById provides a mock function with given fields: id
func (_m *PublicRatingRepository) GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
	ret := _m.Called(id)

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

// GetRatingsBySourceTypeAndActor provides a mock function with given fields: sourceType, sourceUID, filter
func (_m *PublicRatingRepository) GetRatingsBySourceTypeAndActor(sourceType string, sourceUID string, filter publicrequest.GetRatingBySourceTypeAndActorFilter) ([]entity.RatingsCol, error) {
	ret := _m.Called(sourceType, sourceUID, filter)

	var r0 []entity.RatingsCol
	if rf, ok := ret.Get(0).(func(string, string, publicrequest.GetRatingBySourceTypeAndActorFilter) []entity.RatingsCol); ok {
		r0 = rf(sourceType, sourceUID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingsCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, publicrequest.GetRatingBySourceTypeAndActorFilter) error); ok {
		r1 = rf(sourceType, sourceUID, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCounterRatingSubmission provides a mock function with given fields: id, currentCounter
func (_m *PublicRatingRepository) UpdateCounterRatingSubmission(id primitive.ObjectID, currentCounter int64) error {
	ret := _m.Called(id, currentCounter)

	var r0 error
	if rf, ok := ret.Get(0).(func(primitive.ObjectID, int64) error); ok {
		r0 = rf(id, currentCounter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRatingSubDisplayNameByIdLegacy provides a mock function with given fields: input
func (_m *PublicRatingRepository) UpdateRatingSubDisplayNameByIdLegacy(input request.UpdateRatingSubDisplayNameRequest) error {
	ret := _m.Called(input)

	var r0 error
	if rf, ok := ret.Get(0).(func(request.UpdateRatingSubDisplayNameRequest) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateStatusRatingSubHelpful provides a mock function with given fields: id, currentStatus
func (_m *PublicRatingRepository) UpdateStatusRatingSubHelpful(id primitive.ObjectID, currentStatus bool) error {
	ret := _m.Called(id, currentStatus)

	var r0 error
	if rf, ok := ret.Get(0).(func(primitive.ObjectID, bool) error); ok {
		r0 = rf(id, currentStatus)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPublicRatingRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPublicRatingRepository creates a new instance of PublicRatingRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPublicRatingRepository(t mockConstructorTestingTNewPublicRatingRepository) *PublicRatingRepository {
	mock := &PublicRatingRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
