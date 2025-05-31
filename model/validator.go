package model

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func IsValid(i interface{}) (bool, error) {
	errStrs := []string{}
	val := reflect.ValueOf(i)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(i)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := val.Type().FieldByName(err.Field())
			errStr := fmt.Sprintf("%s: %s", val.Type().Field(field.Index[0]).Tag.Get("json"), err.Tag())
			errStrs = append(errStrs, errStr)
		}
	}

	if len(errStrs) > 0 {
		return false, errors.New(strings.Join(errStrs, ", "))
	}

	return true, err
}
