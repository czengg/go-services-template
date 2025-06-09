package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func UnmarshalExternal(data []byte, v interface{}) error {
	return unmarshalWithTag(data, v, "external")
}

func MarshalExternal(v interface{}) ([]byte, error) {
	return marshalWithTag(v, "external")
}

func unmarshalWithTag(data []byte, v interface{}, tagName string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("unmarshal target must be a non-nil pointer")
	}

	rv = rv.Elem()
	rt := rv.Type()

	if rt.Kind() == reflect.Slice {
		return unmarshalSliceWithTag(data, rv, tagName)
	}

	if rt.Kind() != reflect.Struct {
		return json.Unmarshal(data, v)
	}

	return unmarshalStructWithTag(data, rv, rt, tagName)
}

func unmarshalSliceWithTag(data []byte, rv reflect.Value, tagName string) error {
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}

	elementType := rv.Type().Elem()
	newSlice := reflect.MakeSlice(rv.Type(), len(rawMessages), len(rawMessages))

	for i, raw := range rawMessages {
		elem := newSlice.Index(i)

		elemPtr := reflect.New(elementType)
		if err := unmarshalWithTag(raw, elemPtr.Interface(), tagName); err != nil {
			return fmt.Errorf("error unmarshaling slice element %d: %w", i, err)
		}

		elem.Set(elemPtr.Elem())
	}

	rv.Set(newSlice)
	return nil
}

func unmarshalStructWithTag(data []byte, rv reflect.Value, rt reflect.Type, tagName string) error {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}

	fieldMap := buildFieldMap(rt, tagName)

	for externalName, value := range jsonData {
		if fieldInfo, exists := fieldMap[externalName]; exists {
			if err := setFieldValue(rv, fieldInfo, value); err != nil {
				return fmt.Errorf("error setting field %s: %w", fieldInfo.Name, err)
			}
		}
	}

	return nil
}

type fieldInfo struct {
	Index []int
	Name  string
	Type  reflect.Type
}

func buildFieldMap(rt reflect.Type, tagName string) map[string]fieldInfo {
	fieldMap := make(map[string]fieldInfo)
	buildFieldMapRecursive(rt, nil, fieldMap, tagName)
	return fieldMap
}

func buildFieldMapRecursive(rt reflect.Type, indexPrefix []int, fieldMap map[string]fieldInfo, tagName string) {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		index := make([]int, len(indexPrefix)+1)
		copy(index, indexPrefix)
		index[len(indexPrefix)] = i

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			buildFieldMapRecursive(field.Type, index, fieldMap, tagName)
			continue
		}

		tag := field.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

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

func setFieldValue(rv reflect.Value, field fieldInfo, value interface{}) error {
	fieldValue := rv
	for _, idx := range field.Index {
		fieldValue = fieldValue.Field(idx)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field %s cannot be set", field.Name)
	}

	convertedValue, err := convertValue(value, field.Type)
	if err != nil {
		return err
	}

	fieldValue.Set(convertedValue)
	return nil
}

func convertValue(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	if value == nil {
		return reflect.Zero(targetType), nil
	}

	sourceValue := reflect.ValueOf(value)
	sourceType := sourceValue.Type()

	if sourceType.AssignableTo(targetType) {
		return sourceValue, nil
	}

	if sourceType.ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType), nil
	}

	if sourceType.Kind() == reflect.String && targetType.Kind() == reflect.String {
		return reflect.ValueOf(value).Convert(targetType), nil
	}

	if isNumeric(sourceType) && isNumeric(targetType) {
		return sourceValue.Convert(targetType), nil
	}

	if targetType.Kind() == reflect.Slice && sourceType.Kind() == reflect.Slice {
		return convertSlice(sourceValue, targetType)
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %v to %v", sourceType, targetType)
}

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

func marshalWithTag(v interface{}, tagName string) ([]byte, error) {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	if rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return json.Marshal(nil)
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		return marshalSliceWithTag(rv, tagName)
	}

	if rt.Kind() != reflect.Struct {
		return json.Marshal(v)
	}

	return marshalStructWithTag(rv, rt, tagName)
}

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

func marshalStructWithTag(rv reflect.Value, rt reflect.Type, tagName string) ([]byte, error) {
	result := make(map[string]interface{})
	fieldMap := buildFieldMap(rt, tagName)

	for externalName, field := range fieldMap {
		fieldValue := rv
		for _, idx := range field.Index {
			fieldValue = fieldValue.Field(idx)
		}

		if fieldValue.IsZero() {
			continue
		}

		result[externalName] = fieldValue.Interface()
	}

	return json.Marshal(result)
}
