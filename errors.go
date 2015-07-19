package strictjson

import (
	"bytes"
	"fmt"
	"reflect"
)

type Errors []error

func (errs Errors) WithErr(err error) Errors {
	return append(errs, err)
}

func (errs Errors) WithNotAllowed(field string) Errors {
	return append(errs, &NotAllowedError{
		Path: field,
	})
}

func (errs Errors) WithRequired(field string) Errors {
	return append(errs, &RequiredError{
		Path: field,
	})
}

func (errs Errors) WithInvalidType(field, got string, expected reflect.Type) Errors {
	return append(errs, &InvalidTypeError{
		Path:     field,
		Got:      got,
		Expected: expected,
	})
}

// String provides pretty list of all errors combined
func (errs Errors) String() string {
	var b bytes.Buffer
	switch len(errs) {
	case 0:
		b.WriteString("no errors")
	case 1:
		b.WriteString("1 error:\n")
	default:
		fmt.Fprintf(&b, "%d errors:\n", len(errs))
	}

	for _, err := range errs {
		fmt.Fprintf(&b, " %T: %s\n", err, err.Error())
	}
	return b.String()
}

func (errs Errors) Error() string {
	switch len(errs) {
	case 0:
		return "no errors"
	case 1:
		return "1 error"
	default:
		return fmt.Sprintf("%d errors", len(errs))
	}
}

// RequiredError is retuned when no value or empty value was provided for field
// that is required
type RequiredError struct {
	Path string
}

func (err *RequiredError) Error() string {
	return "field is required: " + err.Path
}

// NotAllowedError is retuned when unknown value was provided
type NotAllowedError struct {
	Path string
}

func (err *NotAllowedError) Error() string {
	return "field is not allowed: " + err.Path
}

// InvalidTypeError is returned when JSON type does not match expected data
// format.
type InvalidTypeError struct {
	Path     string
	Expected reflect.Type
	Got      string
}

func (err *InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type: %s: expected %s, got %s",
		err.Path, err.Expected.String(), err.Got)
}
