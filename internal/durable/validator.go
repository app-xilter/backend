package durable

import (
	"backend/internal/model"
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateStruct(s interface{}) error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	return nil
}

func ValidateUrl(s string) error {
	pattern := `.*\/status\/.*`

	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("invalid url")
	}

	return nil
}

func ContainsTag(tags []model.Tag, tag int) bool {
	for _, t := range tags {
		if t.Id == tag {
			return true
		}
	}

	return false
}
