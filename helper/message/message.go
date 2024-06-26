package message

var (
	JSONParseFailCode   = 412001
	ValidationFailCode  = 412002
	UnauthorizedCode    = 412001
	FailConnectCode     = 512003
	TimeOutCode         = 512005
	SuccessCode         = 212000
	DataNotFoundCode    = 212004
	ErrDataNotFoundCode = 412003
)

// Message wrapper.
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var TelErrUserNotFound = Message{Code: ErrDataNotFoundCode, Message: "Not found"}
var ErrDataExists = Message{Code: ValidationFailCode, Message: "Data already exists"}
var ErrBadRouting = Message{Code: FailConnectCode, Message: "Inconsistent mapping between route and handler"}
var ErrInternalError = Message{Code: ValidationFailCode, Message: "Error has been occured while processing request"}
var ErrUnmarshalRequest = Message{Code: ValidationFailCode, Message: "Error can not unmarshal"}
var ErrNoAuth = Message{Code: UnauthorizedCode, Message: "No Authorization"}
var ErrInvalidHeader = Message{Code: 34005, Message: "Invalid header"}
var ErrDB = Message{Code: FailConnectCode, Message: "Error has been occured while processing database request"}
var ErrLTNumState = Message{Code: ValidationFailCode, Message: "Error Num of Statements less Than required Num Statements"}
var ErrGTNumState = Message{Code: ValidationFailCode, Message: "Error Num of Statements greater than  required Num Statements"}
var ErrNoData = Message{Code: ErrDataNotFoundCode, Message: "Data is not found"}
var ErrSaveData = Message{Code: ValidationFailCode, Message: "Data cannot be saved, please check your request"}
var ErrMatchNumState = Message{Code: ValidationFailCode, Message: "Error num_statements does not match number of valid statements"}
var ErrReq = Message{Code: ValidationFailCode, Message: "Required field"}
var ErrTypeReq = Message{Code: ValidationFailCode, Message: "Type required field"}
var ErrTypeFormatReq = Message{Code: ValidationFailCode, Message: "Type is wrong format"}
var ErrIdFormatReq = Message{Code: ValidationFailCode, Message: "Id is wrong format"}
var ErrScaleValueReq = Message{Code: ValidationFailCode, Message: "Scale is wrong value"}
var ErrDuplicateType = Message{Code: ValidationFailCode, Message: "Duplicate type, please check your request"}
var ErrIntervalsValueReq = Message{Code: ValidationFailCode, Message: "Intervals is wrong value"}
var UserAgentTooLong = Message{Code: ValidationFailCode, Message: "The maximum length of user_agent allowed is 200 characters"}
var ErrIPFormatReq = Message{Code: ValidationFailCode, Message: "Wrong IP format"}
var UserUIDRequired = Message{Code: ValidationFailCode, Message: "One of the following display_name, user_id and user_id_legacy must be filled"}
var UserRated = Message{Code: ValidationFailCode, Message: "Duplicate submissions by the same user id, rating id and source_trans_id is not allowed"}
var ErrRatingNotFound = Message{Code: ErrDataNotFoundCode, Message: "Rating not found"}
var ErrRatingNumericTypeNotFound = Message{Code: ErrDataNotFoundCode, Message: "Rating Numeric Type Not Found"}
var RatingSubmissionNotFound = Message{Code: ErrDataNotFoundCode, Message: "Rating submission not found"}
var WrongScoreFilter = Message{Code: ValidationFailCode, Message: "Wrong score filter format"}
var WrongFilter = Message{Code: ValidationFailCode, Message: "Wrong filter"}
var ErrValueFormatForNumericType = Message{Code: ValidationFailCode, Message: "Wrong value format for numeric type"}
var ErrLikertTypeNotFound = Message{Code: ValidationFailCode, Message: "Likert type not found"}
var ErrValueFormat = Message{Code: ValidationFailCode, Message: "Value is wrong format"}
var ErrRatingTypeNotExist = Message{Code: ValidationFailCode, Message: "Rating type not exist"}
var ErrDuplicateRatingName = Message{Code: ValidationFailCode, Message: "Rating name has already existed"}
var ErrSourceNotExist = Message{Code: ValidationFailCode, Message: "Source not exist"}
var ErrFailedToCallGetMedicalFacility = Message{Code: ValidationFailCode, Message: "Failed to call get medical facility"}
var ErrThisRatingTypeIsInUse = Message{Code: ValidationFailCode, Message: "This rating type is in use and has submission"}
var ErrUnmarshalFilterListRatingRequest = Message{Code: ValidationFailCode, Message: "Error can not unmarshal filter param"}
var ErrDataNotFound = Message{Code: ErrDataNotFoundCode, Message: "Data not found"}
var ErrRatingHasRatingSubmission = Message{Code: ValidationFailCode, Message: "Rating has rating submission"}
var ErrMinScoreReq = Message{Code: ValidationFailCode, Message: "Min Score required field"}
var ErrMaxScoreReq = Message{Code: ValidationFailCode, Message: "Max Score required field"}
var ErrScaleReq = Message{Code: ValidationFailCode, Message: "Scale required field"}
var ErrCannotModifiedStatus = Message{Code: ValidationFailCode, Message: "Status can not modified because this rating type in use"}
var ErrCannotModifiedRatingType = Message{Code: ValidationFailCode, Message: "Rating type can not modified because this rating submission in use"}
var ErrCannotModifiedRatingTypeId = Message{Code: ValidationFailCode, Message: "Rating type id not modified because this rating submission in use"}
var ErrCannotModifiedMinScore = Message{Code: ValidationFailCode, Message: "Min Score cannot be modified because this rating type is in use and has submission"}
var ErrCannotModifiedMaxScore = Message{Code: ValidationFailCode, Message: "Max Score cannot be modified because this rating type is in use and has submission"}
var ErrCannotModifiedScale = Message{Code: ValidationFailCode, Message: "Scale cannot be modified because this rating type is in use and has submission"}
var ErrCannotModifiedInterval = Message{Code: ValidationFailCode, Message: "Interval cannot be modified because this rating type is in use and has submission"}
var ErrCannotModifiedStatement = Message{Code: ValidationFailCode, Message: "Statement cannot be modified because this rating type is in use and has submission"}
var ErrCannotModifiedNumStatement = Message{Code: ValidationFailCode, Message: "Num Statement cannot be modified because this rating type is in use and has submission"}
var ErrCannotModifiedType = Message{Code: ValidationFailCode, Message: "Type cannot be modified because this rating type is in use"}
var ErrSourceUidRequire = Message{Code: ValidationFailCode, Message: "Source_uid is required"}
var ErrStoreUidRequire = Message{Code: ValidationFailCode, Message: "store_uid is required"}
var ErrStoreUidMax = Message{Code: ValidationFailCode, Message: "max store_uid is 20"}
var ErrSourceUidMax = Message{Code: ValidationFailCode, Message: "max source_uid is 50"}
var ErrMaxMin = Message{Code: ValidationFailCode, Message: "min_score can not be less than max_score"}
var ErrTypeNotFound = Message{Code: ValidationFailCode, Message: "Rating type not found"}
var ErrCannotSameRatingId = Message{Code: ValidationFailCode, Message: "Cannot create rating submission with same rating"}
var ErrExistingRatingTypeIdSourceUidAndSourceType = Message{Code: ValidationFailCode, Message: "Rating type id, source uid and source type have already existed"}
var ErrRatingSubmissionNotFound = Message{Code: ValidationFailCode, Message: "Rating submission not found"}
var ErrCanNotUpdateSourceTypeOrSoureUid = Message{Code: ValidationFailCode, Message: "Can not update source uid or source type if rating has rating submission"}
var ErrFailedToCalculate = Message{Code: ValidationFailCode, Message: "Failed to calculate rating value"}
var ErrFailedToGetFormula = Message{Code: ValidationFailCode, Message: "Failed to get formula rating"}
var ErrFailedSummaryRatingNumeric = Message{Code: ValidationFailCode, Message: "Failed to summary rating numeric"}
var ErrDisplayNameRequired = Message{Code: ValidationFailCode, Message: "Display Name is required"}
var ErrUserNotFound = Message{Code: ValidationFailCode, Message: "User not found"}
var ErrUserPermissionUpdate = Message{Code: ValidationFailCode, Message: "Not allowed update rating submission of another user"}
var ErrRangeDate = Message{Code: ValidationFailCode, Message: "end_date can not before start_date"}
var ErrInvalidDate = Message{Code: ValidationFailCode, Message: "invalid format date, format should be 2006-01-02"}

// Code 39000 - 39999 Server error
var ErrRevocerRoute = Message{Code: 39000, Message: "Terjadi kesalahan routing"}
var ErrPageNotFound = Message{Code: 39404, Message: "Halaman Tidak ditemukan"}
var SuccessMsg = Message{Code: SuccessCode, Message: "Success"}
var FailedMsg = Message{Code: ValidationFailCode, Message: "Failed"}
var RecordNotFound = Message{Code: ValidationFailCode, Message: "Collection tidak ditemukan"}
var ErrReqParam = Message{Code: ValidationFailCode, Message: "Invalid Request Parameter(s)"}
var ErrFailedRequestToPayment = Message{Code: ValidationFailCode, Message: "Failed to call request to payment service"}

// msg in api get booking Medical facility
var GetMedicalFacilitySuccess = Message{Code: 200, Message: "OK"}
var GetMedicalFacilityNotFound = Message{Code: 400, Message: "Data tidak ditemukan"}

// error code upload media
var ErrUploadMedia = Message{Code: 109400, Message: "Upload Failed"}
