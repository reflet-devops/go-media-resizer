package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/context"
)

func New(ctx *context.Context, options ...validator.Option) *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation(UniqueProjectConf, ValidateUniqueProjectConf())
	return validate
}
