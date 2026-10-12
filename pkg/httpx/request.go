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
	msg   string
	cause error
}

func (e *DecodeError) Error() string {
	return e.msg
}

func (e *DecodeError) Unwrap() error {
	return e.cause
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
		return nil, &DecodeError{msg: validationMessage[T](err), cause: err}
	}
	return &val, nil
}

func validationMessage[T any](err error) string {
	var fieldErrs validator.ValidationErrors
	if !errors.As(err, &fieldErrs) {
		return "invalid request body"
	}
	fields := make([]string, 0, len(fieldErrs))
	for _, fe := range fieldErrs {
		fields = append(fields, jsonFieldName[T](fe.StructField()))
	}
	return "invalid fields: " + strings.Join(fields, ", ")
}

func jsonFieldName[T any](structField string) string {
	sf, ok := reflect.TypeFor[T]().FieldByName(structField)
	if !ok {
		return structField
	}
	name, _, _ := strings.Cut(sf.Tag.Get("json"), ",")
	if name == "" || name == "-" {
		return structField
	}
	return name
}
