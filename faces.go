// Package faces filters a structure using field tags
package faces

import (
	"reflect"
	"strings"
)

const revealTag = "faces"

// Reveal reveals faces of the input based on tags. The input can be:
//
// - A pointer to a structure. In this case, the function will be applied
// directly to the structure.
//
// - An array or a slice. In this case, the function will be applied to each
// element.
//
// - A map. In this case, the function will be applied to each value (of a
// key/value pair).
//
// Examples of struct field tags and their meanings:
//
//   // This field will be revealed with the public face of the structure.
//   Field string `faces:"public"`
//
//   // This field will be revealed with the private and the public faces of the structure.
//   Field string `faces:"private,public"`
//
//   // This field will be revealed with any faces of the structure.
//   Field string
func Reveal(input interface{}, tags ...string) {
	if tags == nil {
		return
	}

	v := reflect.ValueOf(input)

	revealValue(v, tags...)
}

func revealValue(v reflect.Value, tags ...string) {
	// Get the underlying elements if the value is a pointer
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// If v is an array or a slice, reveal each element
	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			revealValue(v.Index(i), tags...)
		}
	}

	// If v is a map, reveal each value
	if v.Kind() == reflect.Map {
		for _, vElement := range v.MapKeys() {
			revealValue(v.MapIndex(vElement), tags...)
		}
	}

	// Finally, make sure that the value is a structure
	if v.Kind() != reflect.Struct {
		return
	}

	// Loop over the fields of the structure
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Recursively reveal nested objects
		if field.Kind() == reflect.Array ||
			field.Kind() == reflect.Map ||
			field.Kind() == reflect.Slice ||
			field.Kind() == reflect.Struct {
			revealValue(field, tags...)
		}

		structField := v.Type().Field(i)

		// Unexported fields are not accessible
		if structField.PkgPath != "" {
			continue
		}

		fieldTags := structField.Tag.Get(revealTag)

		// Keep the original value if the field has no tags and reset the field if
		// it does not match any of the given tags
		if field.CanSet() &&
			fieldTags != "" &&
			!matchTags(strings.Split(fieldTags, ","), tags) {
			field.Set(reflect.New(structField.Type).Elem())
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
