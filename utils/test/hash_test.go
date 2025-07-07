package test

import (
	"testing"

	"github.com/sajad-dev/gingo-helpers/utils"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	passHash := utils.PasswordHashing("haha Your funny")
	assert.Equal(t, passHash, "7a16805bd669693cc8d365269b32dabf5f694e8950e1533628d300e5c4392200")
}
