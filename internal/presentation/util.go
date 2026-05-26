package presentation

import (
	"strconv"
)

func parseID(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
