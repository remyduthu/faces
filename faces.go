// Package faces filters a structure using field tags
package faces

import (
	"reflect"
	"strings"
)

const FILTER_TAG = "faces"

// FilterWithTags resets the values of fields of the input based on tags. It
// panics if the input parameter is not an adress to a structure.
func FilterWithTags(input interface{}, tags ...string) {
	if tags == nil {
		return
	}

	v := reflect.ValueOf(input)

	filterValue(v, tags...)
}

func filterValue(v reflect.Value, tags ...string) {
	// Get the underlying elements if the value is a pointer
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Make sure that the value is a structure
	if v.Kind() != reflect.Struct {
		return
	}

	// Loop over the fields of the structure
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Recursively filter nested structures
		if field.Kind() == reflect.Struct {
			filterValue(v.Field(i), tags...)
		}

		// Filter slice elements
		if field.Kind() == reflect.Slice {
			for j := 0; j < v.Field(i).Len(); j++ {
				filterValue(v.Field(i).Index(j), tags...)
			}
		}

		structField := v.Type().Field(i)

		// We cannot access the value of unexported fields
		if structField.PkgPath != "" {
			continue
		}

		fieldTags := structField.Tag.Get(FILTER_TAG)

		// Keep the original value if the field has no tags and exclude the field if
		// it does not match any of the given tags
		if fieldTags != "" && !matchTags(strings.Split(fieldTags, ","), tags) {
			v.Field(i).Set(reflect.New(structField.Type).Elem())
		}
	}
}

func matchTags(fieldTags []string, tags []string) bool {
	for _, fieldTag := range fieldTags {
		for _, tag := range tags {
			if fieldTag == tag {
				return true
			}
		}
	}

	return false
}
