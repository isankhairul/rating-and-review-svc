package util

import (
	"fmt"
	"go-klikdokter/app/model/request"
	"go-klikdokter/helper/message"
	"strconv"
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
	var interval int
	if scale == 0 || scale == 1 || scale == 2 {
		interval = (max-min)*(scale+1) + 1
	} else {
		return 0
	}
	return interval
}

func ValidInputUpdateRatingTypeNum(input request.EditRatingTypeNumRequest) message.Message {
	if input.Status != nil {
		return message.ErrCannotModifiedStatus
	}
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
