package global

import (
	"html"
	"reflect"
)

func HtmlEscape(req interface{}) {
	value := reflect.ValueOf(req).Elem()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type() != reflect.TypeOf("") {
			continue
		}

		str := field.Interface().(string)
		field.SetString(html.EscapeString(str))
	}
}

func GetSourceTypeByRatingType(ratingType string) string {
	sourceType := "doctor"
	if ratingType == "review_for_layanan" {
		sourceType = "layanan"
	} else if ratingType == "rating_for_product" {
		sourceType = "product"
	} else if ratingType == "rating_for_store" {
		sourceType = "store"
	}

	return sourceType
}

func GetMaximumValueBySourceType(sourceType string) string {
	// default value for product
	var maximumValue string = "5.0"
	if sourceType == "store" {
		maximumValue = "3.0"
	}
	return maximumValue
}

func GetListRatingValueBySourceType(sourceType string) []string {
	// default value for product
	var arrRatingValue []string = []string{"5", "4", "3", "2", "1"}
	if sourceType == "store" {
		arrRatingValue = []string{"3", "2", "1"}
	}
	return arrRatingValue
}
