package test

import (
	"testing"

	"github.com/sajad-dev/gingo-helpers/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	token, err := utils.CreateJWT(map[string]any{"email":"example@gmail.com"})
	assert.NoError(t, err)
	
	parms, ok, err := utils.ValidJWT(token)
	
	assert.True(t, ok)
	assert.NoError(t, err)

	assert.Equal(t, "example@gmail.com", parms.Parameters["email"])
}
