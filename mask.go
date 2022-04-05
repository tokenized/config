package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

// MaskedJSONMarshaller provides an interface for structs to implement so that when the value is
// tagged as masked this marshaller will be called to output a related value that isn't masked.
// For example a private key class can implement this function to output the public key instead even
// though the private key is in the config.
type MaskedJSONMarshaller interface {
	MarshalJSONMasked() ([]byte, error)
}

// Mask returns a map representing a config that is safe to print.
func Mask(v interface{}) map[string]interface{} {
	// a map to contain a masked version of the struct
	out := map[string]interface{}{}

	// get the type
	rt := reflect.TypeOf(v).Elem()

	for i := 0; i < rt.NumField(); i++ {
		// get the struct field
		field := rt.Field(i)

		// field name
		name := field.Name

		// get the value
		fieldValue := reflect.ValueOf(v).Elem().Field(i)

		if field.Type.Kind() == reflect.Struct {
			// This field is a struct, which may contain masked values.
			v := Mask(fieldValue.Addr().Interface())

			if !isEmptyMap(v) {
				out[name] = v
			}

			continue
		}

		// make sure the value is "printable"
		value := fmt.Sprintf("%v", fieldValue)

		if len(value) == 0 {
			// don't add empty elements
			continue
		}

		if len(value) > 0 && field.Tag.Get("masked") == "true" {
			// field is marked to be masked
			value = masked
		}

		// value can be shown as it is
		out[name] = value
	}

	return out
}

// isEmptyMap returns true if a map contains zero keys, false otherwise.
func isEmptyMap(v map[string]interface{}) bool {
	for range v {
		return false
	}

	return true
}

// MarshalJSONMaskedRaw marshals a config into a JSON RawMessage and excludes any "masked" values.
// This can be used to pass into a function that wants a JSON marshaler like a logging function that
// will escape the JSON if it isn't passed as a marshaler object.
func MarshalJSONMaskedRaw(value interface{}) (*json.RawMessage, error) {
	b, err := MarshalJSONMasked(value)
	if err != nil {
		return nil, errors.Wrap(err, "marshal json")
	}

	rawJSON := &json.RawMessage{}
	if err := rawJSON.UnmarshalJSON(b); err != nil {
		return nil, errors.Wrap(err, "unmarshal json")
	}

	return rawJSON, nil
}

// MarshalJSONMasked marshals a config into JSON bytes and excludes any "masked" values.
// The output is meant for display only and can't necessarily be unmarshalled back into the same
// object type because masked values are output as a string value of "***", so if the field type is
// not a string then it will fail.
func MarshalJSONMasked(value interface{}) ([]byte, error) {
	var result []byte
	var fields reflect.Type
	var values reflect.Value
	if reflect.ValueOf(value).Kind() == reflect.Ptr {
		fields = reflect.TypeOf(value).Elem()
		values = reflect.ValueOf(value).Elem()
	} else {
		fields = reflect.TypeOf(value)
		values = reflect.ValueOf(value)
	}

	result = append(result, '{')
	for i := 0; i < fields.NumField(); i++ {
		// get the struct field
		field := fields.Field(i)
		fieldValue := values.Field(i)

		if !fieldValue.CanInterface() {
			continue // not exported
		}

		var b []byte
		var err error
		iface := fieldValue.Interface()
		if field.Tag.Get("masked") == "true" {
			fmt.Printf("field is masked : %s\n", field.Name)
			// Field is masked
			if marshaler, ok := iface.(MaskedJSONMarshaller); ok {
				b, err = marshaler.MarshalJSONMasked()
				if err != nil {
					return nil, errors.Wrapf(err, "marshal masked field: %s", field.Name)
				}
			} else {
				fmt.Printf("field does not have masked marshaller : %s\n", field.Name)
				b = []byte(strconv.Quote("***"))
			}
		} else {
			if marshaler, ok := iface.(json.Marshaler); ok {
				b, err = marshaler.MarshalJSON()
				if err != nil {
					return nil, errors.Wrapf(err, "marshal field: %s", field.Name)
				}
			} else if stringer, ok := iface.(fmt.Stringer); ok {
				b = []byte(strconv.Quote(stringer.String()))
			} else if field.Type.Kind() == reflect.Struct {
				b, err = MarshalJSONMasked(iface)
				if err != nil {
					return nil, errors.Wrapf(err, "marshal struct: %s", field.Name)
				}
			} else {
				b = []byte(strconv.Quote(fmt.Sprintf("%v", iface)))
			}
		}

		if len(result) > 1 {
			result = append(result, ',')
		}

		name := field.Name
		tagName := field.Tag.Get("json")
		if len(tagName) > 0 {
			name = tagName
		}

		result = append(result, []byte(strconv.Quote(name))...)
		result = append(result, ':')
		result = append(result, b...)
	}
	result = append(result, '}')

	return result, nil
}
