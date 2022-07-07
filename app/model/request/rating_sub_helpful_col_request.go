package request

import (
	"go-klikdokter/helper/message"
	"regexp"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters ReqRatingSubHelpfulBody
type ReqRatingSubHelpfulBody struct {
	//  in: body
	Body CreateRatingSubHelpfulRequest `json:"body"`
}

type CreateRatingSubHelpfulRequest struct {
	RatingSubmissionID string `json:"rating_submission_id"`
	UserID             string `json:"user_id"`
	UserIDLegacy       string `json:"user_id_legacy"`
	IPAddress          string `json:"ip_address"`
	UserAgent          string `json:"user_agent"`
}

func (req CreateRatingSubHelpfulRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.RatingSubmissionID, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.IPAddress, validation.Match(regexp.MustCompile(regexIP)).Error(message.ErrIPFormatReq.Message)),
	)
}
