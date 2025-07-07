package test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sajad-dev/gingo-helpers/internal/config"
	"github.com/sajad-dev/gingo-helpers/utils"
	"github.com/stretchr/testify/assert"
)

type InputTest2 struct {
	Name string
}
type InputTest3 struct {
	Name string
}

type InputTest struct {
	Name string
	Inp1 []InputTest2
	Inp2 InputTest3
	A    []string
	File string `file:"yes"`
}

func handlerSendRequest(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(400, gin.H{"error": ""})
		return
	}

	formData := c.Request.MultipartForm.Value



	result := map[string]string{}
	for key, values := range formData {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}

	c.JSON(200, gin.H{"form_fields": result})

}

func TestUtils_SendRequest(t *testing.T) {
	e := utils.CreateServer("/test", "post", []gin.HandlerFunc{handlerSendRequest}, 8080)
	req := InputTest{Name: "haha", A: []string{"a", "b", "c", "d"}, Inp1: []InputTest2{{Name: "ff"}}, File: config.ConfigStUtils.IMAGE_TEST}
	res, err := utils.SendRequest(utils.Request{
		Method:  http.MethodPost,
		Path:    "/test",
		Headers: map[string]string{"Content-Type": "multipart/form-data"},
		Inputs:  req,
		Engin:   e,
	})
	assert.NoError(t, err)
	assert.Equal(t,200,res.Code)
}
