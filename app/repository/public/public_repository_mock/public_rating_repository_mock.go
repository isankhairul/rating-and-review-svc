package public_repository_mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/app/model/request/public"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PublicRatingRepositoryMock) GetRatingsBySourceTypeAndActor(sourceType, sourceUID string, filter publicrequest.GetRatingBySourceTypeAndActorFilter) ([]entity.RatingsCol, error) {
	arguments := repository.Mock.Called(sourceType, sourceUID, filter)

	return arguments.Get(0).([]entity.RatingsCol), nil
}

func (repository *PublicRatingRepositoryMock) GetRatingTypeLikertById(id primitive.ObjectID) (*entity.RatingTypesLikertCol, error) {
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

func (repository *PublicRatingRepositoryMock) GetRatingTypeNumById(id primitive.ObjectID) (*entity.RatingTypesNumCol, error) {
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

func (repository *PublicRatingRepositoryMock) CreateRatingSubHelpful(input request.CreateRatingSubHelpfulRequest) (*entity.RatingSubHelpfulCol, error) {
	ratingSubHelpfulCol := entity.RatingSubHelpfulCol{}
	objectId, _ := primitive.ObjectIDFromHex("629dce7bf1f26275e0d84826")
	ratingSubHelpfulCol.ID = objectId
	return &ratingSubHelpfulCol, nil
}

func (repository *PublicRatingRepositoryMock) UpdateStatusRatingSubHelpful(id primitive.ObjectID, currentStatus bool) error {
	return nil
}

func (repository *PublicRatingRepositoryMock) GetRatingSubHelpfulByRatingSubAndActor(ratingSubId, userIdLegacy string) (*entity.RatingSubHelpfulCol, error) {
	arguments := repository.Mock.Called(ratingSubId, userIdLegacy)
	if arguments.Get(0) == nil {
		return nil, nil
	} else {
		ratingSubHelp := arguments.Get(0).(entity.RatingSubHelpfulCol)
		return &ratingSubHelp, nil
	}
}

func (repository *PublicRatingRepositoryMock) UpdateCounterRatingSubmission(id primitive.ObjectID, currentCounter int64) error {
	return nil
}

func (repository *PublicRatingRepositoryMock) GetPublicRatingsByParams(limit, page, dir int, sort string, filter publicrequest.FilterRatingSummary) ([]entity.RatingsCol, *base.Pagination, error) {
	if sort == "failed" {
		return nil, nil, errors.New("Errors")
	}
	arguments := repository.Mock.Called(limit, page, sort, filter)
	return arguments.Get(0).([]entity.RatingsCol), arguments.Get(1).(*base.Pagination), nil
}

func (repository *PublicRatingRepositoryMock) GetRatingSubsByRatingId(ratingId string) ([]entity.RatingSubmisson, error) {
	if ratingId == "62c3e57b457ed515928c3690" {
		return nil, errors.New("Errors")
	}
	arguments := repository.Mock.Called(ratingId)

	return arguments.Get(0).([]entity.RatingSubmisson), nil
}

func (repository *PublicRatingRepositoryMock) CountRatingSubsByRatingIdAndValue(ratingId, value string) (int64, error) {
	if ratingId == "62c3e57b457ed515928c3690" {
		return 0, errors.New("Errors")
	}
	arguments := repository.Mock.Called(ratingId)
	return int64(arguments.Int(0)), nil
}

func (repository *PublicRatingRepositoryMock) GetPublicRatingSubmissions(limit, page, dir int, sort string, filter publicrequest.FilterRatingSubmission) ([]entity.RatingSubmisson, *base.Pagination, error) {
	if sort == "failed" {
		return nil, nil, errors.New("Errors")
	}
	arguments := repository.Mock.Called(limit, page, sort, filter)
	return arguments.Get(0).([]entity.RatingSubmisson), arguments.Get(1).(*base.Pagination), nil
}

func (repository *PublicRatingRepositoryMock) CreatePublicRatingSubmission(input []request.SaveRatingSubmission) ([]entity.RatingSubmisson, error) {
	for _, arg := range input {
		if *arg.UserID != "629dce7bf1f26275e0d84826" {
			return nil, errors.New("can not be created")
		}
	}
	arguments := repository.Mock.Called(input)
	return arguments.Get(0).([]entity.RatingSubmisson), nil
}

func (repository *PublicRatingRepositoryMock) GetRatingFormulaByRatingTypeIdAndSourceType(ratingTypeId, sourceType string) (*entity.RatingFormulaCol, error) {
	if ratingTypeId == "62c3e57b457ed515928c3690" {
		return nil, errors.New("Errors")
	}
	arguments := repository.Mock.Called(ratingTypeId, sourceType)

	return arguments.Get(0).(*entity.RatingFormulaCol), nil
}

func (repository *PublicRatingRepositoryMock) UpdateRatingSubDisplayNameByIdLegacy(input request.UpdateRatingSubDisplayNameRequest) error {
	if input.DisplayName != "Error" {
		return errors.New("can not be updated")
	}
	return nil
}

func (repository *PublicRatingRepositoryMock) GetListRatingBySourceTypeAndUID(sourceType, sourceUID string) ([]entity.RatingsCol, error) {
	arguments := repository.Mock.Called(sourceType, sourceUID)

	return arguments.Get(0).([]entity.RatingsCol), nil
}
