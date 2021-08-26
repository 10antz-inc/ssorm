package utils

import (
	"cloud.google.com/go/spanner"
	"fmt"
	"reflect"
	"time"
)

func GetTableName(model interface{}) string {
	results := reflect.Indirect(reflect.ValueOf(model))

	if reflect.TypeOf(model).Kind() == reflect.String {
		return model.(string)
	}

	if kind := results.Kind(); kind == reflect.Slice {
		resultType := results.Type().Elem()
		if resultType.Kind() == reflect.Ptr {
			resultType = resultType.Elem()
			elem := reflect.New(resultType).Interface()
			return reflect.TypeOf(elem).Elem().Name()
		} else {
			return resultType.Name()
		}
	}

	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		modelType := reflect.TypeOf(model)
		modelValue := reflect.New(modelType.Elem()).Interface()
		return reflect.TypeOf(modelValue).Elem().Name()
	}

	return results.Type().Name()
}

func GetDeleteColumnName(model interface{}) string {
	results := reflect.Indirect(reflect.ValueOf(model))
	if kind := results.Kind(); kind == reflect.Slice {
		resultType := results.Type().Elem()
		if resultType.Kind() == reflect.Ptr {
			resultType = resultType.Elem()
			elem := reflect.New(resultType).Interface()
			e := reflect.Indirect(reflect.ValueOf(elem))
			for i := 0; i < e.NumField(); i++ {
				tag, varName, _, _ := ReflectValues(e, i)
				if tag.Get(SSORM_TAG_KEY) == SSORM_TAG_DELETE_TIME {
					return varName
				}
			}
		}
		return ""
	}

	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		modelType := reflect.TypeOf(model)
		modelValue := reflect.New(modelType.Elem()).Interface()
		e := reflect.Indirect(reflect.ValueOf(modelValue))
		for i := 0; i < e.NumField(); i++ {
			tag, varName, _, _ := ReflectValues(e, i)
			if tag.Get(SSORM_TAG_KEY) == SSORM_TAG_DELETE_TIME {
				return varName
			}
		}
	}

	if reflect.TypeOf(model).Kind() == reflect.Struct {
		e := reflect.Indirect(reflect.ValueOf(model))
		for i := 0; i < e.NumField(); i++ {
			tag, varName, _, _ := ReflectValues(e, i)
			if tag.Get(SSORM_TAG_KEY) == SSORM_TAG_DELETE_TIME {
				return varName
			}
		}
	}

	return ""
}

func ArrayContains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func ReflectValues(reflectValue reflect.Value, i int) (reflect.StructTag, string, interface{}, reflect.Type) {
	value := reflect.TypeOf(reflectValue.Interface())
	tag := value.Field(i).Tag
	varName := reflectValue.Type().Field(i).Name
	varType := reflectValue.Type().Field(i).Type
	varValue := reflectValue.Field(i).Interface()
	return tag, varName, varValue, varType
}

func GetTimestampStr(value interface{}) string {
	timestampStr := "NULL"
	switch v := value.(type) {
	case time.Time:
		if !v.IsZero() {
			timestampStr = fmt.Sprintf("TIMESTAMP_MILLIS(%d)", v.UnixNano()/int64(time.Millisecond))
		}
	case spanner.NullTime:
		if !v.IsNull() {
			timestampStr = fmt.Sprintf("TIMESTAMP_MILLIS(%d)", v.Time.UnixNano()/int64(time.Millisecond))
		}
	}
	return timestampStr
}

func IsTime(value interface{}) bool {
	switch value.(type) {
	case time.Time:
		return true
	case spanner.NullTime:
		return true
	}
	return false
}
