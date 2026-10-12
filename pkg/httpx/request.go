package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type DecodeError struct {
	msg    string
	fields Fields
	cause  error
}

func (e *DecodeError) Error() string {
	return e.msg
}

func (e *DecodeError) Unwrap() error {
	return e.cause
}

func (e *DecodeError) Fields() Fields {
	return e.fields
}

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(jsonTagName)
	return v
}

func jsonTagName(field reflect.StructField) string {
	name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	return name
}

func DecodeJSON[T any](r *http.Request, v *validator.Validate) (*T, error) {
	var val T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&val)
	if err != nil {
		return nil, &DecodeError{msg: "invalid request body", cause: err}
	}
	err = v.Struct(&val)
	if err != nil {
		return nil, validationError(err)
	}
	return &val, nil
}

func validationError(err error) *DecodeError {
	var fieldErrs validator.ValidationErrors
	if !errors.As(err, &fieldErrs) {
		return &DecodeError{msg: "invalid request body", cause: err}
	}
	fields := make(Fields, len(fieldErrs))
	for _, fe := range fieldErrs {
		fields[fe.Field()] = FieldError{Reason: fe.ActualTag(), Param: fe.Param()}
	}
	return &DecodeError{msg: "invalid fields", fields: fields, cause: err}
}
