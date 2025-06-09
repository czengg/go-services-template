// this would be in shared package
package upwardli

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// UnmarshalExternal unmarshals JSON data using the "external" struct tags
func UnmarshalExternal(data []byte, v interface{}) error {
	return unmarshalWithTag(data, v, "external")
}

// MarshalExternal marshals data using the "external" struct tags
func MarshalExternal(v interface{}) ([]byte, error) {
	return marshalWithTag(v, "external")
}

// unmarshalWithTag unmarshals JSON using the specified struct tag
func unmarshalWithTag(data []byte, v interface{}, tagName string) error {
	// Get the reflect value and type
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("unmarshal target must be a non-nil pointer")
	}

	// Get the element being pointed to
	rv = rv.Elem()
	rt := rv.Type()

	// Handle slices
	if rt.Kind() == reflect.Slice {
		return unmarshalSliceWithTag(data, rv, tagName)
	}

	// Handle structs
	if rt.Kind() != reflect.Struct {
		// Fallback to regular JSON unmarshaling for non-structs
		return json.Unmarshal(data, v)
	}

	return unmarshalStructWithTag(data, rv, rt, tagName)
}

// unmarshalSliceWithTag handles slice unmarshaling
func unmarshalSliceWithTag(data []byte, rv reflect.Value, tagName string) error {
	// First unmarshal into []json.RawMessage to get individual elements
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}

	// Create a new slice of the correct type
	elementType := rv.Type().Elem()
	newSlice := reflect.MakeSlice(rv.Type(), len(rawMessages), len(rawMessages))

	// Unmarshal each element
	for i, raw := range rawMessages {
		elem := newSlice.Index(i)

		// Create a pointer to the element for unmarshaling
		elemPtr := reflect.New(elementType)
		if err := unmarshalWithTag(raw, elemPtr.Interface(), tagName); err != nil {
			return fmt.Errorf("error unmarshaling slice element %d: %w", i, err)
		}

		// Set the element value
		elem.Set(elemPtr.Elem())
	}

	// Set the slice
	rv.Set(newSlice)
	return nil
}

// unmarshalStructWithTag handles struct unmarshaling
func unmarshalStructWithTag(data []byte, rv reflect.Value, rt reflect.Type, tagName string) error {
	// Create a map to hold the JSON data
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	// Create a mapping from external field names to struct fields
	fieldMap := buildFieldMap(rt, tagName)

	// Set values on the struct
	for externalName, value := range jsonData {
		if fieldInfo, exists := fieldMap[externalName]; exists {
			if err := setFieldValue(rv, fieldInfo, value); err != nil {
				return fmt.Errorf("error setting field %s: %w", fieldInfo.Name, err)
			}
		}
	}

	return nil
}

// fieldInfo holds information about a struct field
type fieldInfo struct {
	Index []int
	Name  string
	Type  reflect.Type
}

// buildFieldMap creates a mapping from external tag names to field info
func buildFieldMap(rt reflect.Type, tagName string) map[string]fieldInfo {
	fieldMap := make(map[string]fieldInfo)
	buildFieldMapRecursive(rt, nil, fieldMap, tagName)
	return fieldMap
}

// buildFieldMapRecursive recursively builds the field map, handling embedded structs
func buildFieldMapRecursive(rt reflect.Type, indexPrefix []int, fieldMap map[string]fieldInfo, tagName string) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// Build the index path for this field
		index := make([]int, len(indexPrefix)+1)
		copy(index, indexPrefix)
		index[len(indexPrefix)] = i

		// Handle embedded structs
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			buildFieldMapRecursive(field.Type, index, fieldMap, tagName)
			continue
		}

		// Get the tag value
		tag := field.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		// Parse tag options (e.g., "field_name,omitempty")
		tagParts := strings.Split(tag, ",")
		externalName := tagParts[0]

		if externalName != "" {
			fieldMap[externalName] = fieldInfo{
				Index: index,
				Name:  field.Name,
				Type:  field.Type,
			}
		}
	}
}

// setFieldValue sets a field value using reflection
func setFieldValue(rv reflect.Value, field fieldInfo, value interface{}) error {
	// Navigate to the field using the index path
	fieldValue := rv
	for _, idx := range field.Index {
		fieldValue = fieldValue.Field(idx)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field %s cannot be set", field.Name)
	}

	// Convert the value to the correct type
	convertedValue, err := convertValue(value, field.Type)
	if err != nil {
		return err
	}

	fieldValue.Set(convertedValue)
	return nil
}

// convertValue converts an interface{} value to the target reflect.Type
func convertValue(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	if value == nil {
		return reflect.Zero(targetType), nil
	}

	sourceValue := reflect.ValueOf(value)
	sourceType := sourceValue.Type()

	// Direct assignment if types match
	if sourceType.AssignableTo(targetType) {
		return sourceValue, nil
	}

	// Handle type conversions
	if sourceType.ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType), nil
	}

	// Handle string to custom types (like enums)
	if sourceType.Kind() == reflect.String && targetType.Kind() == reflect.String {
		return reflect.ValueOf(value).Convert(targetType), nil
	}

	// Handle numeric conversions
	if isNumeric(sourceType) && isNumeric(targetType) {
		return sourceValue.Convert(targetType), nil
	}

	// Handle slice types
	if targetType.Kind() == reflect.Slice && sourceType.Kind() == reflect.Slice {
		return convertSlice(sourceValue, targetType)
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %v to %v", sourceType, targetType)
}

// convertSlice converts a slice to the target slice type
func convertSlice(sourceValue reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	sourceLen := sourceValue.Len()
	targetSlice := reflect.MakeSlice(targetType, sourceLen, sourceLen)
	elementType := targetType.Elem()

	for i := 0; i < sourceLen; i++ {
		sourceElem := sourceValue.Index(i)
		convertedElem, err := convertValue(sourceElem.Interface(), elementType)
		if err != nil {
			return reflect.Value{}, err
		}
		targetSlice.Index(i).Set(convertedElem)
	}

	return targetSlice, nil
}

// isNumeric checks if a type is numeric
func isNumeric(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// marshalWithTag marshals using the specified struct tag
func marshalWithTag(v interface{}, tagName string) ([]byte, error) {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	// Handle pointers
	if rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return json.Marshal(nil)
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}

	// Handle slices
	if rt.Kind() == reflect.Slice {
		return marshalSliceWithTag(rv, tagName)
	}

	// Handle structs
	if rt.Kind() != reflect.Struct {
		return json.Marshal(v)
	}

	return marshalStructWithTag(rv, rt, tagName)
}

// marshalSliceWithTag marshals a slice using the specified tag
func marshalSliceWithTag(rv reflect.Value, tagName string) ([]byte, error) {
	length := rv.Len()
	result := make([]interface{}, length)

	for i := 0; i < length; i++ {
		elem := rv.Index(i)
		marshaledElem, err := marshalWithTag(elem.Interface(), tagName)
		if err != nil {
			return nil, err
		}

		var elemInterface interface{}
		if err := json.Unmarshal(marshaledElem, &elemInterface); err != nil {
			return nil, err
		}
		result[i] = elemInterface
	}

	return json.Marshal(result)
}

// marshalStructWithTag marshals a struct using the specified tag
func marshalStructWithTag(rv reflect.Value, rt reflect.Type, tagName string) ([]byte, error) {
	result := make(map[string]interface{})
	fieldMap := buildFieldMap(rt, tagName)

	// Reverse the field map to go from field info to external name
	for externalName, field := range fieldMap {
		// Get the field value
		fieldValue := rv
		for _, idx := range field.Index {
			fieldValue = fieldValue.Field(idx)
		}

		// Skip zero values for omitempty
		if fieldValue.IsZero() {
			// You could implement omitempty logic here if needed
			continue
		}

		result[externalName] = fieldValue.Interface()
	}

	return json.Marshal(result)
}
