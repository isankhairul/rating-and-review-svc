package repository_mock

import (
	"errors"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PublicRatingRepositoryMock struct {
	Mock mock.Mock
}

func (repository *PublicRatingRepositoryMock) GetRatingsBySourceTypeAndActor(sourceType, sourceUID string) ([]entity.RatingsCol, error) {
	arguments := repository.Mock.Called(sourceType, sourceUID)

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

func (repository *PublicRatingRepositoryMock) UpdateCounterRatingSubmission(id primitive.ObjectID, currentCounter int) error {
	return nil
}
