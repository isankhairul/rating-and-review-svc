package util

import "strings"

func MaskDisplayName(displayName string) string {
	value := ""
	splitName := strings.Fields(displayName)
	for _, name := range splitName {
		if value == "" {
			value = name[0:1] + "***"
		} else {
			value = value + " " + name[0:1] + "***"
		}
	}
	return value
}
