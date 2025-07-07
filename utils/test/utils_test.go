package test

import (
	"testing"

	"github.com/sajad-dev/gingo-helpers/utils"
	"github.com/stretchr/testify/assert"
)

func TestUtils_GenerateToken(t *testing.T) {
	t1 := utils.GenerateToken()
	t2 := utils.GenerateToken()
	assert.NotEqual(t,t1,t2)
}
