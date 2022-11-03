package request_dapr

// swagger:parameters PublisherRequest
type PublisherRequest struct {
	// in: formData
	// example: queuing.rnr.ratingsubcol
	Topic string `json:"topic" schema:"topic" binding:"omitempty"`

	// in: formData
	// example: {"comment_allowed":true,"description":"Rating Group for Rumah Sakit RS Pondok Indah Bintaro Jaya","name":"Rumah Sakit - RS Pondok Indah Bintaro Jaya","rating_type":"standard-0.0-to-5.0","rating_type_id":"629dc84ff16b9b21f357a609","source_type":"hospital","source_uid":"2729","status":true}
	Data string `json:"data" schema:"data" binding:"omitempty"`
}

// swagger:parameters SubscriberRatingsubcolRequest
type SubscriberRatingsubcolRequest struct {
	// in: body
	Body BodySubscriberRatingsubcolRequest `json:"body"`
}

type BodySubscriberRatingsubcolRequest struct {
	// in: body
	// example: Foo Bar
	Data string `json:"data"`
}
