package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sajad-dev/gingo-helpers/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConvertStringToKind(s string, kind reflect.Kind) (reflect.Value, error) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i).Convert(reflect.TypeOf(int(0))), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(u).Convert(reflect.TypeOf(uint(0))), nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f).Convert(reflect.TypeOf(float64(0))), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b), nil
	case reflect.String:
		return reflect.ValueOf(s), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported kind %v", kind)
	}
}

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}
	err = db.AutoMigrate(config.ConfigStUtils.DATABASE...)
	if err != nil {
		panic(err.Error())
	}
	return db
}
func CreateServer(path string, method string, controller []gin.HandlerFunc, port int) *gin.Engine {
	http := gin.Default()
	switch method {
	case "get":
		http.GET(path, controller...)
	case "post":
		http.POST(path, controller...)
	case "put":
		http.PUT(path, controller...)
	case "patch":
		http.PATCH(path, controller...)
	case "delete":
		http.DELETE(path, controller...)

	}
	return http
}

func CheckValidationErr(res string, filed string, tag string) (bool, error) {
	var dec map[string]any

	err := json.Unmarshal([]byte(res), &dec)
	if err != nil {
		return false, err
	}

	var body []any

	errorsList, ok := dec["errors"]
	if !ok {
		return false, errors.New("errors not in response ")
	}

	body, ok = errorsList.([]any)
	if !ok {
		return false, errors.New("error pattern is not valid")
	}

	for _, value := range body {
		val, ok := value.([]any)
		if !ok {
			return false, errors.New("error pattern is not valid")
		}

		fieldName, _ := val[0].(string)
		errorMesage, _ := val[1].(string)
		if fieldName == filed && strings.Contains(errorMesage, fmt.Sprintf("'%s' tag", tag)) {
			return true, nil
		}
	}

	return false, nil
}
