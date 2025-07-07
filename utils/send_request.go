package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Inputs  any
	Engin   *gin.Engine
}

func sliceHandel(writer *multipart.Writer, value reflect.Value, prefix string) error {
	sliceLen := value.Len()

	for i := 0; i < sliceLen; i++ {
		val := value.Index(i).Interface()
		rt := reflect.TypeOf(val)
		rval := reflect.ValueOf(val)
		var err error
		fieldPrefix := prefix + "[" + strconv.Itoa(i) + "]"
		switch rt.Kind() {
		case reflect.Struct:
			err = writerInput(writer, val, fieldPrefix)
		case reflect.Slice:
			err = sliceHandel(writer, rval, fieldPrefix)
		default:
			err = writer.WriteField(fieldPrefix, rval.Interface().(string))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func writerInput(writer *multipart.Writer, input any, prefix string) error {
	inputType := reflect.TypeOf(input)
	inputValue := reflect.ValueOf(input)
	for i := 0; i < inputType.NumField(); i++ {

		field := inputType.Field(i)
		value := inputValue.Field(i)
		var err error

		fieldName := strcase.ToSnake(field.Name)
		if prefix != "" {
			fieldName = prefix + "[" + fieldName + "]"
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			err = writerInput(writer, value.Interface(), fieldName)
		case reflect.Slice:
			err = sliceHandel(writer, value, fieldName)
		default:
			fileType := field.Tag.Get("file")

			if len(fileType) != 0 {
				file, err := os.Open(value.Interface().(string))
				if err != nil {
					return err
				}

				part, err := writer.CreateFormFile(fieldName, file.Name())
				if err != nil {
					return err
				}

				_, err = io.Copy(part, file)
				if err != nil {
					return err
				}

			} else {
				err = writer.WriteField(fieldName, fmt.Sprintf("%v", value.Interface()))

			}

		}

		if err != nil {
			return err
		}

	}

	return nil
}
func (r *Request) multipartFormData() (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writerInput(writer, r.Inputs, "")
	if err != nil {
		return nil, err
	}
	writer.Close()

	req := httptest.NewRequest(r.Method, r.Path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func (r *Request) jsonFormData() (*http.Request, error) {
	body, _ := json.Marshal(r.Inputs)
	req, _ := http.NewRequest(r.Method, r.Path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (r *Request) switchContentType() (*http.Request, error) {
	switch r.Headers["Content-Type"] {
	case "multipart/form-data":
		return r.multipartFormData()
	default:
		return r.jsonFormData()
	}
}
func SendRequest(request Request) (*httptest.ResponseRecorder, error) {
	req, err := request.switchContentType()
	if err != nil {
		return nil, err
	}

	for key, value := range request.Headers {
		if key != "Content-Type" {
			req.Header.Set(key, value)
		}
	}

	wr := httptest.NewRecorder()
	request.Engin.ServeHTTP(wr, req)
	return wr, err
}
