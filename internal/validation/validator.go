package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		structName := strings.Split(err.Namespace(), ".")[0]

		fields[err.Field()] = fmt.Sprintf("failed '%s' tag check (value '%s' is not valid for %s struct)",
			err.Tag(), err.Value(), structName)
	}

	return fields
}
