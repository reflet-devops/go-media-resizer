package cli

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/context"
	validatorMediaResize "github.com/reflet-devops/go-media-resizer/validator"
)

func validateConfig(ctx *context.Context) error {
	validate := validatorMediaResize.New(ctx)
	err := validate.Struct(ctx.Config)
	if err != nil {

		var validationErrors validator.ValidationErrors
		switch {
		case errors.As(err, &validationErrors):
			for _, validationError := range validationErrors {
				ctx.Logger.Error(fmt.Sprintf("%v", validationError))
			}
			return errors.New("configuration file is not valid")
		default:
			return err
		}
	}
	return nil
}
