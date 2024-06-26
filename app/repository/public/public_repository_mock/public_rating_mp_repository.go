// Code generated by mockery v2.16.0. DO NOT EDIT.

package public_repository_mock

import (
	base "go-klikdokter/app/model/base"
	entity "go-klikdokter/app/model/entity"

	mock "github.com/stretchr/testify/mock"

	publicrequest "go-klikdokter/app/model/request/public"

	publicresponse "go-klikdokter/app/model/response/public"
)

// PublicRatingMpRepository is an autogenerated mock type for the PublicRatingMpRepository type
type PublicRatingMpRepository struct {
	mock.Mock
}

// CountRatingSubsByRatingIdAndValue provides a mock function with given fields: ratingId, value
func (_m *PublicRatingMpRepository) CountRatingSubsByRatingIdAndValue(ratingId string, value string) (int64, error) {
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

// GetListRatingBySourceTypeAndUID provides a mock function with given fields: sourceType, sourceUID
func (_m *PublicRatingMpRepository) GetListRatingBySourceTypeAndUID(sourceType string, sourceUID string) ([]entity.RatingsMpCol, error) {
	ret := _m.Called(sourceType, sourceUID)

	var r0 []entity.RatingsMpCol
	if rf, ok := ret.Get(0).(func(string, string) []entity.RatingsMpCol); ok {
		r0 = rf(sourceType, sourceUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingsMpCol)
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
func (_m *PublicRatingMpRepository) GetPublicRatingSubmissions(limit int, page int, dir int, sort string, filter publicrequest.FilterRatingSubmissionMp) ([]entity.RatingSubmissionMp, *base.Pagination, error) {
	ret := _m.Called(limit, page, dir, sort, filter)

	var r0 []entity.RatingSubmissionMp
	if rf, ok := ret.Get(0).(func(int, int, int, string, publicrequest.FilterRatingSubmissionMp) []entity.RatingSubmissionMp); ok {
		r0 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingSubmissionMp)
		}
	}

	var r1 *base.Pagination
	if rf, ok := ret.Get(1).(func(int, int, int, string, publicrequest.FilterRatingSubmissionMp) *base.Pagination); ok {
		r1 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*base.Pagination)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, int, int, string, publicrequest.FilterRatingSubmissionMp) error); ok {
		r2 = rf(limit, page, dir, sort, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPublicRatingSubmissionsCustom provides a mock function with given fields: limit, page, dir, sort, filter, source
func (_m *PublicRatingMpRepository) GetPublicRatingSubmissionsCustom(limit int, page int, dir int, sort string, filter publicrequest.FilterRatingSubmissionMp, source string) ([]entity.RatingSubmissionMp, *base.Pagination, error) {
	ret := _m.Called(limit, page, dir, sort, filter, source)

	var r0 []entity.RatingSubmissionMp
	if rf, ok := ret.Get(0).(func(int, int, int, string, publicrequest.FilterRatingSubmissionMp, string) []entity.RatingSubmissionMp); ok {
		r0 = rf(limit, page, dir, sort, filter, source)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingSubmissionMp)
		}
	}

	var r1 *base.Pagination
	if rf, ok := ret.Get(1).(func(int, int, int, string, publicrequest.FilterRatingSubmissionMp, string) *base.Pagination); ok {
		r1 = rf(limit, page, dir, sort, filter, source)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*base.Pagination)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, int, int, string, publicrequest.FilterRatingSubmissionMp, string) error); ok {
		r2 = rf(limit, page, dir, sort, filter, source)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPublicRatingSubmissionsGroupBySource provides a mock function with given fields: limit, page, dir, sort, filter
func (_m *PublicRatingMpRepository) GetPublicRatingSubmissionsGroupBySource(filter publicrequest.FilterRatingSummary) ([]publicresponse.PublicRatingSubGroupBySourceMp, error) {
	ret := _m.Called(filter)

	var r0 []publicresponse.PublicRatingSubGroupBySourceMp
	if rf, ok := ret.Get(0).(func(publicrequest.FilterRatingSummary) []publicresponse.PublicRatingSubGroupBySourceMp); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]publicresponse.PublicRatingSubGroupBySourceMp)
		}
	}

	var r2 error
	if rf, ok := ret.Get(1).(func(publicrequest.FilterRatingSummary) error); ok {
		r2 = rf(filter)
	} else {
		r2 = ret.Error(1)
	}

	return r0, r2
}

// GetPublicRatingSubmissionsGroupByStoreSource provides a mock function with given fields: limit, page, dir, sort, filter
func (_m *PublicRatingMpRepository) GetPublicRatingSubmissionsGroupByStoreSource(filter publicrequest.FilterRatingSummary) ([]publicresponse.PublicRatingSubGroupByStoreSourceMp, error) {
	ret := _m.Called(filter)

	var r0 []publicresponse.PublicRatingSubGroupByStoreSourceMp
	if rf, ok := ret.Get(0).(func(publicrequest.FilterRatingSummary) []publicresponse.PublicRatingSubGroupByStoreSourceMp); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]publicresponse.PublicRatingSubGroupByStoreSourceMp)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(publicrequest.FilterRatingSummary) error); ok {
		r2 = rf(filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r2
}

// GetPublicRatingsByParams provides a mock function with given fields: limit, page, dir, sort, filter
func (_m *PublicRatingMpRepository) GetPublicRatingsByParams(limit int, page int, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsMpCol, *base.Pagination, error) {
	ret := _m.Called(limit, page, dir, sort, filter)

	var r0 []entity.RatingsMpCol
	if rf, ok := ret.Get(0).(func(int, int, int, string, publicrequest.FilterRatingSummary) []entity.RatingsMpCol); ok {
		r0 = rf(limit, page, dir, sort, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingsMpCol)
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

// GetRatingFormulaBySourceType provides a mock function with given fields: sourceType
func (_m *PublicRatingMpRepository) GetRatingFormulaBySourceType(sourceType string) (*entity.RatingFormulaCol, error) {
	ret := _m.Called(sourceType)

	var r0 *entity.RatingFormulaCol
	if rf, ok := ret.Get(0).(func(string) *entity.RatingFormulaCol); ok {
		r0 = rf(sourceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.RatingFormulaCol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sourceType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRatingSubsByRatingId provides a mock function with given fields: ratingId
func (_m *PublicRatingMpRepository) GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmissionMp, error) {
	ret := _m.Called(ratingId)

	var r0 []entity.RatingSubmissionMp
	if rf, ok := ret.Get(0).(func(string) []entity.RatingSubmissionMp); ok {
		r0 = rf(ratingId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.RatingSubmissionMp)
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

// GetSumCountRatingSubsByRatingId provides a mock function with given fields: ratingId
func (_m *PublicRatingMpRepository) GetSumCountRatingSubsByRatingId(ratingId string) (*publicresponse.PublicSumCountRatingSummaryMp, error) {
	ret := _m.Called(ratingId)

	var r0 *publicresponse.PublicSumCountRatingSummaryMp
	if rf, ok := ret.Get(0).(func(string) *publicresponse.PublicSumCountRatingSummaryMp); ok {
		r0 = rf(ratingId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*publicresponse.PublicSumCountRatingSummaryMp)
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

// GetSumCountRatingSubsBySource provides a mock function with given fields: sourceUID, sourceType
func (_m *PublicRatingMpRepository) GetSumCountRatingSubsBySource(sourceUID string, sourceType string) (*publicresponse.PublicSumCountRatingSummaryMp, error) {
	ret := _m.Called(sourceUID, sourceType)

	var r0 *publicresponse.PublicSumCountRatingSummaryMp
	if rf, ok := ret.Get(0).(func(string, string) *publicresponse.PublicSumCountRatingSummaryMp); ok {
		r0 = rf(sourceUID, sourceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*publicresponse.PublicSumCountRatingSummaryMp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(sourceUID, sourceType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *PublicRatingMpRepository) GetRatingSubsGroupByValue(sourceUid string, sourceType string) ([]interface{}, error) {
	ret := _m.Called(sourceUid, sourceType)

	var r0 []interface{}
	if rf, ok := ret.Get(0).(func(string, string) []interface{}); ok {
		r0 = rf(sourceUid, sourceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interface{})
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

type mockConstructorTestingTNewPublicRatingMpRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPublicRatingMpRepository creates a new instance of PublicRatingMpRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPublicRatingMpRepository(t mockConstructorTestingTNewPublicRatingMpRepository) *PublicRatingMpRepository {
	mock := &PublicRatingMpRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
