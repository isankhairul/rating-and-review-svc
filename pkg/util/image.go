package util

import (
	"errors"
	"net/http"
	"strings"
)

var DefaultMimeAllowed = [] string {
	"image/jpeg",
	"image/png",
}

var DefaultSizeAllowed int64 = 2000000

type StructValidationImage struct {
	FileName string
	FileBytes []byte
	MimeAllowed []string
	SizeAllowed int64
}

func NewValidationImage(fileName string, fileBytes []byte, mimeAllowed *[]string, sizeAllowed *int64) *StructValidationImage {
	mime := *mimeAllowed
	size := *sizeAllowed

	if mimeAllowed == nil {
		mime = DefaultMimeAllowed
	}

	if sizeAllowed == nil {
		size = DefaultSizeAllowed
	}

	return &StructValidationImage{
		FileName: fileName,
		FileBytes: fileBytes,
		MimeAllowed: mime,
		SizeAllowed: size,
	}
}

func (s StructValidationImage) ValidateSize() error {
	// validation max size 2 mb
	if int64(len(s.FileBytes)) > s.SizeAllowed {
		return errors.New("Max file size 2 mb")
	}
	return nil
}

func (s StructValidationImage) ValidateMime() error {
	var mime = http.DetectContentType(s.FileBytes)
	var isMimeAllowed bool
	//validation mim
	for _, mimeAllowed := range s.MimeAllowed {
		if mime == mimeAllowed {
			isMimeAllowed = true
		}
	}

	if !isMimeAllowed {
		return errors.New("mime/extension only allow " + strings.Join(s.MimeAllowed, ", "))
	}

	return nil
}

func (s StructValidationImage) ValidateSizeAndMime() error {
	// validation extension
	if err := s.ValidateMime(); err != nil {
		return err
	}

	// validation size max 2 mb

	if err := s.ValidateSize(); err != nil {
		return err
	}
	return nil
}
