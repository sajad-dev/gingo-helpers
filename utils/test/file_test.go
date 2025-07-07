package test

import (
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sajad-dev/gingo-helpers/core/validation"
	"github.com/sajad-dev/gingo-helpers/internal/config"
	"github.com/sajad-dev/gingo-helpers/utils"
	"github.com/stretchr/testify/assert"
)

type TestStructFileValidSend struct {
	File string `file:"yes"`
}
type TestStructFileValid struct {
	File []*multipart.FileHeader `json:"file" validate_file:"size=10"`
}

func handlerFileRequest(ctx *gin.Context) {
	reqParams, _ := ctx.Get("req")

	req := reqParams.(TestStructFileValid)
	path, err := utils.SaveFile(ctx, "", req.File, "sajad")
	if err != nil {
		ctx.JSON(500, "")
	}
	ctx.JSON(200, path)
}

func ValidationMiddlewareTe(validationSt any) func(*gin.Context) {

	return func(ctx *gin.Context) {

		st, err := validation.SwitchHeader(ctx, validationSt)
		if err != nil {
			ctx.JSON(500, "")
			return
		}

		ctx.Set("req", st.Interface())
		ctx.Next()
	}

}

func TestFile(t *testing.T) {

	r := utils.CreateServer("/test_storage", "post", []gin.HandlerFunc{ValidationMiddlewareTe(TestStructFileValid{}), handlerFileRequest}, 8080)

	w, err := utils.SendRequest(utils.Request{Method: http.MethodPost, Path: "/test_storage", Engin: r, Headers: map[string]string{"Content-Type": "multipart/form-data"},
		Inputs: TestStructFileValidSend{File: config.ConfigStUtils.IMAGE_TEST}})

	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "sajad-")
}
