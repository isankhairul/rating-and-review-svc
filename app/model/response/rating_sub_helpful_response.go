package response

type RatingSubHelpfulResponse struct {
	RatingSubmissionId string `json:"rating_submission_id"`
	UserIdLegacy       string `json:"user_id_legacy"`
	LikeCounter        int    `json:"like_counter"`
	Status             string `json:"status"`
}
