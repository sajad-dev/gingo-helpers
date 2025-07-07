package validation

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"github.com/sajad-dev/gingo-helpers/utils"
	// "github.com/sajad-dev/gingo/internal/app/validation"
	// "github.com/sajad-dev/gingo/internal/exception"
)

func sliceHandel(valueOf reflect.Value, key []string, value string) error {

	fieldOf := reflect.New(valueOf.Type().Elem()).Elem()

	switch fieldOf.Kind() {
	case reflect.Slice:
		return sliceHandel(fieldOf, key[1:], value)

	case reflect.Struct:
		num, _ := strconv.Atoi(key[0])
		if valueOf.Len() == num {
			newval := reflect.New(fieldOf.Type()).Elem()
			err := readInput(newval, key[1:], value)
			fieldOf.Set(newval)
			newvalAppend := reflect.Append(valueOf, newval)
			valueOf.Set(newvalAppend)
			return err
		} else {
			valend := valueOf.Index(valueOf.Len() - 1)
			err := readInput(valend, key[1:], value)

			return err

		}
	case reflect.String:
		valueOf.Set(reflect.Append(valueOf, reflect.ValueOf(value)))

	default:
		convertedVal, err := utils.ConvertStringToKind(value, fieldOf.Kind())
		if err != nil {
			return err
		}
		appendCon := reflect.Append(fieldOf, convertedVal)
		fieldOf.Set(appendCon)
	}
	return nil
}

func readInput(valueOf reflect.Value, key []string, value string) error {
	fieldName := strcase.ToCamel(key[0])

	fieldOf := valueOf.FieldByName(fieldName)

	if !fieldOf.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}
	if !fieldOf.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	switch fieldOf.Kind() {
	case reflect.Slice:
		err := sliceHandel(fieldOf, key[1:], value)
		return err
	case reflect.Struct:
		newval := reflect.New(fieldOf.Type()).Elem()
		err := readInput(newval, key[1:], value)
		fieldOf.Set(newval)
		return err
	case reflect.String:
		fieldOf.SetString(value)
	default:
		convertedVal, err := utils.ConvertStringToKind(value, fieldOf.Kind())
		if err != nil {
			return err
		}
		fieldOf.Set(convertedVal)
	}

	return nil
}
func readFile(valueOf reflect.Value, key []string, value []*multipart.FileHeader) error {
	fieldName := strcase.ToCamel(key[0])
	fieldOf := valueOf.FieldByName(fieldName)

	if !fieldOf.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}
	if !fieldOf.CanSet() {
		return fmt.Errorf("cannot set field %s", fieldName)
	}

	switch fieldOf.Kind() {
	case reflect.Struct:
		newval := reflect.New(fieldOf.Type()).Elem()
		err := readFile(newval, key[1:], value)
		fieldOf.Set(newval)
		return err

	default:
		fieldOf.Set(reflect.ValueOf(value))

	}
	return nil
}
func multipartHandel(ctx *gin.Context, fieldValidation any) (any, error) {
	ctx.Request.ParseMultipartForm(10 << 20)

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	result := map[string]string{}

	for key, value := range ctx.Request.MultipartForm.Value {
		result[key] = value[0]
	}

	valueOf := reflect.ValueOf(fieldValidation)

	if valueOf.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("field must be a pointer to a struct")
	}
	valueOf = valueOf.Elem()
	if valueOf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("field must point to a struct")
	}

	for key, value := range result {
		tree := strings.Split(strings.ReplaceAll(key, "]", ""), "[")
		readInput(valueOf, tree, value)
	}

	for key, value := range form.File {
		tree := strings.Split(strings.ReplaceAll(key, "]", ""), "[")
		readFile(valueOf, tree, value)
	}
	return fieldValidation, nil
}
func SwitchHeader(ctx *gin.Context, validationSt any) (reflect.Value, error) {

	formType := reflect.TypeOf(validationSt)
	form := reflect.New(formType)

	content := ctx.GetHeader("Content-Type")
	switch {
	case strings.Contains(content, "multipart/form-data"):
		parm, err := multipartHandel(ctx, form.Interface())
		return reflect.ValueOf(parm).Elem(), err

	default:
		formIn := form.Interface()
		err := ctx.ShouldBindJSON(formIn)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return form.Elem(), nil
	}
}
