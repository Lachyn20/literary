package validation

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func Struct(s interface{}) error {
    return Validate.Struct(s)
}
