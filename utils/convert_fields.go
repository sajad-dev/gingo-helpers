package utils

import (
	"encoding/json"
	"mime/multipart"
	"reflect"

	"gorm.io/datatypes"
)

func checkJson(field reflect.Value, value reflect.Value) error {
	kind := field.Kind()
	if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
		return nil
	}
	if kind == reflect.Struct || kind == reflect.Slice {

		json, err := json.Marshal(field.Interface())
		if err != nil {
			return err
		}

		var jsonData datatypes.JSON = json
		value.Set(reflect.ValueOf(jsonData))

	} else {
		value.Set(reflect.ValueOf(field.Interface()))
	}
	return nil
}

func ConvertValidationToTable(validationSt any, otherField any, tableSt any) error {
	valueOf := reflect.ValueOf(tableSt).Elem()
	typeOf := reflect.TypeOf(tableSt).Elem()

	valOfValidation := reflect.ValueOf(validationSt).Elem()
	valOfOther := reflect.ValueOf(otherField).Elem()

	for i := 0; i < valueOf.NumField(); i++ {
		fieldValue := valueOf.Field(i)
		fieldType := typeOf.Field(i)
		name := fieldType.Name

		switch {
		case validationSt != nil && valOfValidation.FieldByName(name).IsValid():
			err := checkJson(valOfValidation.FieldByName(name), fieldValue)
			if err != nil {
				return err
			}
		case otherField != nil && valOfOther.FieldByName(name).IsValid():
			err := checkJson(valOfOther.FieldByName(name), fieldValue)
			if err != nil {
				return err
			}

		}
	}

	return nil
}
