package strictjson

import (
	"encoding/json"
	"reflect"
	"strings"
)

func Unmarshal(b []byte, dest interface{}) Errors {
	var errs Errors

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return errs.WithErr(&json.InvalidUnmarshalError{
			Type: reflect.TypeOf(dest),
		})
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		// TODO support more than structure as root type
		return errs.WithErr(&json.InvalidUnmarshalError{
			Type: reflect.TypeOf(dest),
		})
	}

	return unmarshal(b, val, "")
}

func unmarshal(b []byte, dest reflect.Value, path string) Errors {
	kind := dest.Kind()
	v := dest
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
		kind = v.Kind()
	}

	var errs Errors
	switch kind {
	case reflect.Struct:
		return unmarshalStruct(b, dest, path)
	case
		reflect.Slice,
		reflect.Map,
		reflect.String,
		reflect.Int, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:

		if err := json.Unmarshal(b, dest.Interface()); err != nil {
			if e, ok := err.(*json.UnmarshalTypeError); ok {
				return errs.WithInvalidType(path, e.Value, e.Type)
			}
			return errs.WithErr(err)
		}
	case reflect.Invalid:
		panic("something went wrong")
	default:
		return errs.WithErr(&json.UnsupportedTypeError{
			Type: reflect.TypeOf(dest),
		})
	}
	return nil
}

func unmarshalStruct(b []byte, dest reflect.Value, path string) (errs Errors) {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &raw); err != nil {
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			return errs.WithInvalidType(path, e.Value, e.Type)
		}
		return errs.WithErr(err)
	}

	v := dest
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	tp := v.Type()

	fields := make(map[string]struct{})

	for i := 0; i < v.NumField(); i++ {
		def := definition(tp.Field(i))
		fv := v.Field(i)

		fields[def.name] = struct{}{}

		fraw, ok := raw[def.name]
		if !ok {
			if def.required {
				errs = errs.WithRequired(fieldPath(path, def.name))
			}
			continue
		}

		if fv.Kind() == reflect.Ptr {
			v := reflect.New(fv.Type().Elem())
			v.Elem().Set(reflect.Zero(v.Type().Elem()))
			fv.Set(v)
		}

		if es := unmarshal(fraw, fv.Addr(), fieldPath(path, def.name)); es != nil {
			errs = append(errs, es...)
			continue
		}

		// field might be present but the value is "null"
		if def.required && def.IsEmpty(fv) {
			errs = errs.WithRequired(def.name)
		}
	}

	for name := range raw {
		if _, ok := fields[name]; !ok {
			errs = errs.WithNotAllowed(fieldPath(path, name))
		}
	}

	return errs
}

type fieldDefinition struct {
	required bool
	ignore   bool
	name     string
}

// definition return JSON metadata for given struct field
func definition(f reflect.StructField) *fieldDefinition {
	tags := strings.Split(f.Tag.Get("json"), ",")
	ignore := false
	if len(tags) > 0 || tags[0] == "-" {
		ignore = true
	}

	return &fieldDefinition{
		required: isRequired(f, tags),
		name:     jsonName(f),
		ignore:   ignore,
	}
}

// isEmpty return true if given value is nil or empty.
func (f *fieldDefinition) IsEmpty(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}

// isRequired return true if non-zero field value must be provided.
func isRequired(f reflect.StructField, tags []string) bool {
	if f.Type.Kind() == reflect.Ptr {
		return false
	}
	if tags[0] == "-" {
		return false
	}
	if len(tags) >= 2 {
		for _, t := range tags[1:] {
			if t == "omitempty" {
				return false
			}
		}
	}

	if f.Type.Kind() == reflect.Slice {
		return f.Type.Elem().Kind() != reflect.Ptr
	}
	return true
}

// jsonName return name that should be used as JSON representation.
func jsonName(v reflect.StructField) string {
	tags := v.Tag.Get("json")
	if tags != "" {
		name := strings.SplitN(tags, ",", 2)[0]
		if name == "" {
			return v.Name
		}
		return name
	}
	return v.Name
}

// fieldPath return full field path combined from root path and field name.
func fieldPath(root, name string) string {
	if root != "" {
		return root + "." + name
	}
	return name
}
