// Package faces filters a structure using field tags
package faces

import (
	"reflect"
	"strings"
)

const FILTER_TAG = "faces"

// FilterWithTags resets the...
// input must be a pointer to a structure.
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

	// Get the type of the value
	t := v.Type()

	// Loop over the fields of the structure
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// We cannot access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}

		fieldTags := field.Tag.Get(FILTER_TAG)

		// Keep the original value if the field has no tags
		if fieldTags == "" {
			continue
		}

		// Recursively filter nested structures
		if field.Type.Kind() == reflect.Struct {
			FilterWithTags(v.Field(i), tags...)
		}

		// Filter slice elements
		if field.Type.Kind() == reflect.Slice {
			for j := 0; j < v.Field(i).Len(); j++ {
				FilterWithTags(v.Field(i).Index(j), tags...)
			}
		}

		// Exclude the field if it does not match any of the given tags
		if !matchTags(strings.Split(fieldTags, ","), tags) {
			v.Field(i).Set(reflect.New(field.Type).Elem())
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
