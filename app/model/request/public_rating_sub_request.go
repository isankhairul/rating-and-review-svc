package request

// swagger:parameters ReqUpdateRatingSubDisplayNameBody
type ReqUpdateRatingSubDisplayNameBody struct {
	// in: path
	// required: true
	UserIdLegacy string `json:"user_id_legacy"`
	// in: body
	Body UpdateRatingSubDisplayNameRequest `json:"body"`
}

type UpdateRatingSubDisplayNameRequest struct {
	UserIdLegacy string `json:"-"`
	DisplayName  string `json:"display_name,omitempty"`
}
