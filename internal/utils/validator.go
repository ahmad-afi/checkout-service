package utils

import (
	"bytes"
	"checkout-service/internal/helper"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Message []string
}

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

func (v ValidationError) Error() string {
	buff := bytes.NewBufferString("")

	for i := 0; i < len(v.Message); i++ {
		errMessage := v.Message[i]
		buff.WriteString(errMessage)
		buff.WriteString(", ")
	}
	finalStr := strings.TrimSpace(buff.String())
	return finalStr
}

func Validator(data any) *helper.ErrorStruct {
	var validate = validator.New()

	if err := validate.Struct(data); err != nil {
		errs := err.(validator.ValidationErrors)
		var newErr ValidationError

		for _, val := range errs {
			// fmt.Printf("error :  %#v \n\n ", val)
			// fmt.Printf("error :  %#v \n\n ", val.StructNamespace() )

			// getField, _ := reflect.TypeOf(data).FieldByName(val.Field())
			// jsonTag := getField.Tag.Get("json")

			var message string
			switch val.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", val.Field())
			default:
				message = fmt.Sprintf("validation error for '%s', Tag: %s", val.Field(), val.Tag())
			}

			newErr.Message = append(newErr.Message, message)
		}

		return &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  newErr,
		}
	}

	return nil
}
