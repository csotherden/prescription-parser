package handlerutils

import (
	"strconv"
	"strings"
)

func ParseIntParam(v string, d int) int {
	if strings.TrimSpace(v) == "" {
		return d
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return d
	}

	return i
}
