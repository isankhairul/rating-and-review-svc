package message

var (
	BadRequestCode     = 412002
	SuccessCode        = 212000
	DataNotFoundCode   = 201001
	InternalServerCode = 501000
)

// Message wrapper.
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var TelErrUserNotFound = Message{Code: 34000, Message: "Not found"}
var ErrDataExists = Message{Code: 34001, Message: "Data already exists"}
var ErrBadRouting = Message{Code: 34002, Message: "Inconsistent mapping between route and handler"}
var ErrInternalError = Message{Code: InternalServerCode, Message: "Error has been occured while processing request"}
var ErrUnmarshalRequest = Message{Code: 412001, Message: "Error can not unmarshal"}
var ErrNoAuth = Message{Code: 34004, Message: "No Authorization"}
var ErrInvalidHeader = Message{Code: 34005, Message: "Invalid header"}
var ErrDB = Message{Code: 34005, Message: "Error has been occured while processing database request"}
var ErrNoData = Message{Code: 212004, Message: "Data is not found"}
var ErrSaveData = Message{Code: 412002, Message: "Data cannot be saved, please check your request"}
var ErrReq = Message{Code: 34005, Message: "Required field"}
var ErrTypeReq = Message{Code: 401001, Message: "Type required field"}
var ErrTypeFormatReq = Message{Code: 401002, Message: "Type is wrong format"}
var ErrIdFormatReq = Message{Code: 401003, Message: "Id is wrong format"}
var ErrScaleValueReq = Message{Code: 401004, Message: "Scale is wrong value"}
var ErrIntervalsValueReq = Message{Code: 401005, Message: "Intervals is wrong value"}
var UserAgentTooLong = Message{Code: BadRequestCode, Message: "The maximum length of user_agent allowed is 200 characters"}
var ErrIPFormatReq = Message{Code: InternalServerCode, Message: "Wrong IP format"}
var UserUIDRequired = Message{Code: BadRequestCode, Message: "One of the following user_id and user_id_legacy must be filled"}
var UserRated = Message{Code: BadRequestCode, Message: "Duplicate submissions by the same user for rating is not allowed"}
var ErrRatingNotFound = Message{Code: DataNotFoundCode, Message: "Rating not found"}
var ErrRatingNumericTypeNotFound = Message{Code: DataNotFoundCode, Message: "Rating Numeric Type Not Found"}
var ErrValueFormat = Message{Code: BadRequestCode, Message: "Wrong format Rating submission value"}

// Code 39000 - 39999 Server error
var ErrRevocerRoute = Message{Code: 39000, Message: "Terjadi kesalahan routing"}
var ErrPageNotFound = Message{Code: 39404, Message: "Halaman Tidak ditemukan"}
var SuccessMsg = Message{Code: SuccessCode, Message: "Success"}
var FailedMsg = Message{Code: BadRequestCode, Message: "Failed"}
var ErrReqParam = Message{Code: BadRequestCode, Message: "Invalid Request Parameter(s)"}
