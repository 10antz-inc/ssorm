package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
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
	varName := tag.Get(SPANER_KEY)
	if varName == "" {
		varName = reflectValue.Type().Field(i).Name
	}
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

func IsNullable(value interface{}) bool {
	switch value.(type) {
	case spanner.NullInt64, spanner.NullFloat64, spanner.NullString, spanner.NullDate, spanner.NullTime, spanner.NullRow,
		*spanner.NullInt64, *spanner.NullFloat64, *spanner.NullString, *spanner.NullDate, *spanner.NullTime, *spanner.NullRow:
		return true
	}
	return false
}

func IsValid(value interface{}) bool {
	if value == nil {
		return false
	}
	switch value.(type) {
	case spanner.NullInt64:
		return value.(spanner.NullInt64).Valid
	case *spanner.NullInt64:
		return value.(*spanner.NullInt64).Valid
	case spanner.NullFloat64:
		return value.(spanner.NullFloat64).Valid
	case *spanner.NullFloat64:
		return value.(*spanner.NullFloat64).Valid
	case spanner.NullString:
		return value.(spanner.NullString).Valid
	case *spanner.NullString:
		return value.(*spanner.NullString).Valid
	case spanner.NullDate:
		return value.(spanner.NullDate).Valid
	case *spanner.NullDate:
		return value.(*spanner.NullDate).Valid
	case spanner.NullTime:
		return value.(spanner.NullTime).Valid
	case *spanner.NullTime:
		return value.(*spanner.NullTime).Valid
	case spanner.NullRow:
		return value.(spanner.NullRow).Valid
	case *spanner.NullRow:
		return value.(*spanner.NullRow).Valid
	}
	return false
}

func GetArrayStr(value interface{}, valType reflect.Type) string {
	var res string
	var stringVal []string
	var valFormat string
	switch valType.String() {
	case "[]string", "[]*string", "[]spanner.NullString", "[]*spanner.NullString", "[]civil.Date", "[]*civil.Date", "[]spanner.NullDate", "[]*spanner.NullDate":
		valFormat = `"%v"`
	default:
		valFormat = "%v"
	}

	val := reflect.ValueOf(value)
	for i := 0; i < val.Len(); i++ {
		stringVal = append(stringVal, fmt.Sprintf(valFormat, reflect.Indirect(val.Index(i)).Interface()))
	}
	elms := strings.Join(stringVal, ",")
	res = fmt.Sprintf("[%v]", elms)
	return res
}

func IsTypeString(valType reflect.Type) bool {
	if valType.Kind() == reflect.String {
		return true
	}
	switch valType.String() {
	case "spanner.NullString", "*spanner.NullString", "civil.Date", "spanner.NullDate", "*civil.Date", "*spanner.NullDate":
		return true
	}
	return false
}

