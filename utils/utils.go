package utils

import (
	"reflect"
)

func GetTableName(model interface{}) string {
	results := Indirect(reflect.ValueOf(model))

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

func Indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}
