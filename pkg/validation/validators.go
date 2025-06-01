package validation

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func init() {
	validate := validator.New(validator.WithRequiredStructEnabled())
	Validate = validate
}
