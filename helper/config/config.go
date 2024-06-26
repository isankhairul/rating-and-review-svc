package config

import (
	"os"
	"strconv"
	"strings"
)

func GetConfigString(conf string) string {
	if strings.HasPrefix(conf, "${") && strings.HasSuffix(conf, "}") {
        return os.Getenv(strings.TrimSuffix(strings.TrimPrefix(conf, "${"), "}"))
    }

	return conf
}

func GetConfigInt(conf string) int {
	if strings.HasPrefix(conf, "${") && strings.HasSuffix(conf, "}") {
		result, _ := strconv.Atoi(os.Getenv(strings.TrimSuffix(strings.TrimPrefix(conf, "${"), "}")))
        return result
    }

	result, _ := strconv.Atoi(conf)
	return result
}