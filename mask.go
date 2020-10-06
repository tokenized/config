package config

import (
	"fmt"
	"reflect"
)

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
