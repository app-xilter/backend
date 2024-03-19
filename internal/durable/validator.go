package durable

import (
	"backend/internal/model"
	"github.com/go-playground/validator/v10"
)

func ValidateStruct(t model.Tweets) error {
	validate := validator.New()
	if err := validate.Struct(t); err != nil {
		return err
	}

	return nil
}
