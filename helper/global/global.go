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
