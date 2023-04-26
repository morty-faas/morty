package config

import (
	"fmt"
	"reflect"
	"strings"
)

// isDefined check if a configuration is defined using the reflection API.
func isDefined(v any) bool {
	if v == nil {
		return false
	}
	return !reflect.ValueOf(v).IsZero()
}

// ensureKeyHasSingleSubKey check using the reflection API, if the key has only one sub key defined.
func ensureKeyHasSingleSubKey(key interface{}) error {
	var fields []string
	v := reflect.ValueOf(key)
	t := reflect.Indirect(v)

	// Iterate over key fields
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fName := strings.ToLower(t.Type().Field(i).Name)

		// If the field isn't zero (e.g: {}) and
		// we didn't match a sub-key already
		if !f.IsZero() {
			fields = append(fields, fName)
			if len(fields) > 1 {
				return fmt.Errorf("found multiple sub-keys defined for configuration key `%s`, but only one is allowed: %v", strings.ToLower(v.Type().Name()), fields)
			}
		}
	}

	return nil
}
