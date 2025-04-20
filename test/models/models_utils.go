package models_testing

import (
	"database/sql/driver"
	"reflect"
	"strings"
)

func StructToSlice(s interface{}) []driver.Value {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		panic("The input parameter must be a structure type")
	}

	numFields := v.NumField()
	result := make([]interface{}, numFields)

	for i := 0; i < numFields; i++ {
		result[i] = v.Field(i).Interface()
	}

	return convertToDriverValues(result)
}

func StructToUpdateSlice(s interface{}) []driver.Value {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		panic("The input parameter must be a structure type")
	}

	numFields := v.NumField()
	result := make([]interface{}, 0, numFields)

	var idValue interface{}
	foundID := false

	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		if strings.ToLower(field.Name) == "id" {
			idValue = value
			foundID = true
		} else {
			result = append(result, value)
		}
	}

	if foundID {
		result = append(result, idValue)
	}

	return convertToDriverValues(result)
}

func convertToDriverValues(input []interface{}) []driver.Value {
	values := make([]driver.Value, len(input))
	for i, v := range input {
		values[i] = driver.Value(v)
	}
	return values
}
