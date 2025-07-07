package utils

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sajad-dev/gingo-helpers/internal/config"
)

func SaveFile(ctx *gin.Context, namespace string, files []*multipart.FileHeader, name string) ([]string, error) {
	timestamp := time.Now().Unix()
	var dst []string
	for _, file := range files {
		indexSt := strings.LastIndex(file.Filename, ".")
		name := namespace + "/" + name + "-" + strconv.FormatInt(timestamp, 10) + "." + file.Filename[indexSt+1:]
		dstPath := config.ConfigStUtils.STORAGE_PATH + name
		dst = append(dst, name)
		err := ctx.SaveUploadedFile(file, dstPath)
		if err != nil {
			return []string{}, err
		}
	}

	return dst, nil
}

func SetMultipartFields(ctx *gin.Context, obj any) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Not valid :(")
	}
	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		if fieldType.Type == reflect.TypeOf((*multipart.FileHeader)(nil)) {
			if field.IsNil() {
				field.Set(reflect.New(fieldType.Type.Elem())) // Initialize the field as a pointer
			}
			if file, err := ctx.FormFile(string(fieldType.Tag.Get("json"))); err == nil {
				field.Set(reflect.ValueOf(file))
			}
		}
	}
	return nil
}
