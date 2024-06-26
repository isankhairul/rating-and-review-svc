package util

import (
	"errors"
	"fmt"
	"go-klikdokter/app/model/entity"
	"go-klikdokter/app/model/request"
	"go-klikdokter/helper/message"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidValue(min int, max int, interval int, scale int) []float64 {
	var values []float64
	increase := float64(max-min) / (float64(interval) - 1)
	for i := float64(min); i <= float64(max); i += increase {
		i = RoundFloatWithPrecision(i, scale)
		values = append(values, i)
	}
	return values
}

func RoundFloatWithPrecision(number float64, precision int) float64 {
	rounded, err := strconv.ParseFloat(fmt.Sprintf("%."+fmt.Sprintf("%d", precision)+"f", number), 64)
	if err != nil {
		return number
	}
	return rounded
}

func ValidateValue(vArray []float64, v float64) bool {
	for _, args := range vArray {
		if v == args {
			return true
		}
	}
	return false
}

func ValidInterval(min int, max int, scale int) int {
	var interval = 0
	if scale == 0 || scale == 1 || scale == 2 {
		interval = (max-min)*(scale+1) + 1
	}
	return interval
}

func ValidInputUpdateRatingTypeNumRated(input request.EditRatingTypeNumRequest) message.Message {
	if input.Status != nil {
		return message.ErrCannotModifiedStatus
	}
	if input.Type != "" {
		return message.ErrCannotModifiedType
	}
	return message.SuccessMsg
}

func ValidInputUpdateRatingTypeNumSubmission(input request.EditRatingTypeNumRequest) message.Message {
	if input.MinScore != nil {
		return message.ErrCannotModifiedMinScore
	}
	if input.MaxScore != nil {
		return message.ErrCannotModifiedMaxScore
	}
	if input.Intervals != nil {
		return message.ErrCannotModifiedInterval
	}
	if input.Scale != nil {
		return message.ErrCannotModifiedScale
	}
	return message.SuccessMsg
}

func ValidInputUpdateRatingTypeLikertInRating(input request.SaveRatingTypeLikertRequest) message.Message {
	var errMsg message.Message
	if input.Status != nil {
		errMsg = message.ErrCannotModifiedStatus
		return errMsg
	}
	if input.Type != "" {
		errMsg = message.ErrCannotModifiedType
		return errMsg
	}
	return errMsg
}

func ValidInputUpdateRatingTypeLikertInSubmission(input request.SaveRatingTypeLikertRequest) message.Message {
	var errMsg message.Message
	if input.Status != nil {
		errMsg = message.ErrCannotModifiedStatus
		return errMsg
	}
	if input.Type != "" {
		errMsg = message.ErrCannotModifiedType
		return errMsg
	}
	if input.NumStatements != 0 {
		errMsg = message.ErrCannotModifiedNumStatement
		return errMsg
	}
	if input.Statement01 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement02 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement03 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement04 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement05 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement06 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement07 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement08 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement09 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Statement10 != nil {
		errMsg = message.ErrCannotModifiedStatement
		return errMsg
	}
	if input.Status != nil {
		errMsg = message.ErrCannotModifiedStatus
		return errMsg
	}
	return errMsg
}

func ValidateTypeNumeric(input *entity.RatingTypesNumCol, value float64) message.Message {
	// The value must be valid according to requirements of rating type
	values := ValidValue(*input.MinScore, *input.MaxScore, *input.Intervals, *input.Scale)
	isInclude := ValidateValue(values, value)
	if isInclude == false {
		return message.Message{
			Code:    message.ValidationFailCode,
			Message: "value is only 1 and must be included in : " + fmt.Sprintf("%v", values),
		}
	}
	return message.SuccessMsg
}

func ValidateLikertType(input *entity.RatingTypesLikertCol, value []string) (error, []int) {
	wrongValue := "wrong value"
	validValue := make([]int, 0)
	emtpy := ""
	if input.Statement01 != nil && *input.Statement01 != emtpy {
		validValue = append(validValue, 1)
	}
	if input.Statement02 != nil && *input.Statement02 != emtpy {
		validValue = append(validValue, 2)
	}
	if input.Statement03 != nil && *input.Statement03 != emtpy {
		validValue = append(validValue, 3)
	}
	if input.Statement04 != nil && *input.Statement04 != emtpy {
		validValue = append(validValue, 4)
	}
	if input.Statement05 != nil && *input.Statement05 != emtpy {
		validValue = append(validValue, 5)
	}
	if input.Statement06 != nil && *input.Statement06 != emtpy {
		validValue = append(validValue, 6)
	}
	if input.Statement07 != nil && *input.Statement07 != emtpy {
		validValue = append(validValue, 7)
	}
	if input.Statement08 != nil && *input.Statement08 != emtpy {
		validValue = append(validValue, 8)
	}
	if input.Statement09 != nil && *input.Statement09 != emtpy {
		validValue = append(validValue, 9)
	}
	if input.Statement10 != nil && *input.Statement10 != emtpy {
		validValue = append(validValue, 10)
	}
	for _, args := range value {
		v, err := strconv.ParseFloat(args, 64)
		if err != nil {
			return errors.New(wrongValue), validValue
		}
		if v-float64(int(v)) > 0 || v == 0 {
			return errors.New(wrongValue), validValue
		}
		if isInclude := IsInclude(validValue, v); isInclude == false {
			return errors.New(wrongValue), validValue
		}
	}

	return nil, validValue
}

func IsInclude(arrValue []int, value float64) bool {
	for _, args := range arrValue {
		if int(value) == args {
			return true
		}
	}
	return false
}

func ValidateUserIdAndUserIdLegacy(sourceTransId, ratingId string, id *string, idLegacy *string, ratingSubmission *entity.RatingSubmisson, err error) bool {
	if ratingSubmission != nil {
		if ratingSubmission.RatingID == ratingId && ratingSubmission.UserID != nil {
			if *ratingSubmission.UserID == *id && sourceTransId == ratingSubmission.SourceTransID {
				return true
			}
		}
		if ratingSubmission.RatingID == ratingId && ratingSubmission.UserIDLegacy != nil {
			if *ratingSubmission.UserIDLegacy == *idLegacy && sourceTransId == ratingSubmission.SourceTransID {
				return true
			}
		}
		if err == nil {
			return true
		}
	}
	return false
}

func ValidateUserIdAndUserIdLegacyMp(sourceTransId, ratingId string, id *string, idLegacy *string, ratingSubmission *entity.RatingSubmissionMp, err error) bool {
	if ratingSubmission != nil {
		if ratingSubmission.RatingID == ratingId && ratingSubmission.UserID != nil {
			if *ratingSubmission.UserID == *id && sourceTransId == ratingSubmission.SourceTransID {
				return true
			}
		}
		if ratingSubmission.RatingID == ratingId && ratingSubmission.UserIDLegacy != nil {
			if *ratingSubmission.UserIDLegacy == *idLegacy && sourceTransId == ratingSubmission.SourceTransID {
				return true
			}
		}
		if err == nil {
			return true
		}
	}
	return false
}

func ValidateUserIdAndUserIdLegacyForUpdate(input request.UpdateRatingSubmissionRequest, subId primitive.ObjectID, sourceTransId string, id *string, idLegacy *string, ratingSubmission *entity.RatingSubmisson, err error) bool {
	if ratingSubmission != nil {
		if ratingSubmission.ID == subId {
			return false
		}

		if input.RatingID == ratingSubmission.RatingID && ratingSubmission.UserID != nil {
			if *ratingSubmission.UserID == *id && ratingSubmission.SourceTransID == sourceTransId {
				return true
			}
		}

		if input.RatingID == ratingSubmission.RatingID && ratingSubmission.UserIDLegacy != nil {
			if *ratingSubmission.UserIDLegacy == *idLegacy && ratingSubmission.SourceTransID == sourceTransId {
				return true
			}
		}

		if err != nil {
			return true
		}
	}
	return false
}

func ValidateUserIdAndUserIdLegacyForUpdateMp(input request.UpdateRatingSubmissionRequest, subId primitive.ObjectID, sourceTransId string, id *string, idLegacy *string, ratingSubmission *entity.RatingSubmissionMp, err error) bool {
	if ratingSubmission != nil {
		if ratingSubmission.ID == subId {
			return false
		}

		if ratingSubmission.UserID != nil {
			if *ratingSubmission.UserID == *id && ratingSubmission.SourceTransID == sourceTransId {
				return true
			}
		}

		if ratingSubmission.UserIDLegacy != nil {
			if *ratingSubmission.UserIDLegacy == *idLegacy && ratingSubmission.SourceTransID == sourceTransId {
				return true
			}
		}

		if err != nil {
			return true
		}
	}
	return false
}

func ValidateUserCannotUpdateMp(userID string, userIDLegacy string, ratingSub entity.RatingSubmissionMp) bool {
	if userID != "" && *ratingSub.UserID != userID {
		return true
	}

	if userIDLegacy != "" && *ratingSub.UserIDLegacy != userIDLegacy {
		return true
	}

	return false
}