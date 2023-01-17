package upload_request

import "os"

// swagger:parameters ReqUploadForm
type ReqUploadForm struct {
	// Image file
	// in: formData
	// swagger:file
	Image os.File `json:"image"`
	// Source
	// in: formData
	SourceType string `json:"source_type"`
	// Source UID
	// in: formData
	SourceUID string `json:"source_uid"`
}

type UploadImageRequest struct {
	// Filename
	FileName string `json:"file_name"`
	// Image byte
	Image []byte `json:"image"`
	// Source Type
	SourceType string `json: "source_type"`
	// Source Uid
	SourceUid string `json: "source_uid"`
	// Mime Type
	MimeType string `json: "mime_type"`
}

// func (r ReqUploadForm) Validate() error {
// 	f := r.Image
// 	statF, _ := f.Stat()
// 	var sizeAllowed int64 = 1024 * 1024 * 2
// 	var extensionAllowed = map[string]bool{
// 		"image/jpeg": true,
// 		"image/png": true,
// 	}

// }
