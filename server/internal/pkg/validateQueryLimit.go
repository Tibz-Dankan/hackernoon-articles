package pkg

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/constants"
)

// ValidateQueryLimit function takes in limitParam as an argument  and
// checks to ensure the query limit is not exceeded
func ValidateQueryLimit(limitParam string) (float64, error) {
	defaultLimit := float64(12)
	minimum := constants.QUERY_MINIMUM_LIMIT
	maximum := constants.QUERY_MAXIMUM_LIMIT

	if limitParam == "" {
		return defaultLimit, nil
	}

	limit, err := strconv.ParseFloat(limitParam, 64)
	if err != nil || limit < minimum || limit > maximum {
		errorMsg := fmt.Sprintf("Provided limit is out of range, min %d and max %d", int(minimum), int(maximum))
		return defaultLimit, errors.New(errorMsg)
	}

	return limit, nil
}
