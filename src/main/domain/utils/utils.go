package utils

import (
	"regexp"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

func ValidateHexID(ids []string) apierrors.ApiError {

	// Regular expression to check if a string is a valid hex

	regex := regexp.MustCompile("^[a-fA-F0-9]{24}$")

	for _, id := range ids {
		val := regex.MatchString(id)
		if !val {
			return apierrors.NewApiError("one or more of the provided ids are not a valid hex string", "bad_request", 400, apierrors.CauseList{})
		}
	}

	return nil
}
