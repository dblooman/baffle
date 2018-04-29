package validate

import (
	"regexp"

	"github.com/dblooman/baffle/server/backends"
	"github.com/dblooman/baffle/server/logger"
)

func Ensure(data backends.CreateSecret) (bool, error) {

	r, err := regexp.Compile(data.Regex)
	if err != nil {
		return false, logger.WrapError(err, "Unable to parse regex")
	}

	match := r.MatchString(data.Secret)
	if !match {
		return false, logger.WrapError(err, "Regex failed")
	}

	return true, nil
}
