package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/config"
)

const (
	UniqueProjectConf = "unique-project-cfg"
)

func ValidateUniqueProjectConf() func(level validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		projects, ok := fl.Field().Interface().([]config.Project)
		if !ok {
			return false
		}

		for i := 0; i < len(projects); i++ {
			for j := i + 1; j < len(projects); j++ {
				if projects[i].PrefixPath == projects[j].PrefixPath &&
					projects[i].Hostname == projects[j].Hostname {
					return false
				}
			}
		}
		return true
	}
}
