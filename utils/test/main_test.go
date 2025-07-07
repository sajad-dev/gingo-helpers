package test

import (
	"testing"

	"github.com/sajad-dev/gingo-helpers/core/bootsterap"
	"github.com/sajad-dev/gingo-helpers/types"
)

func TestMain(m *testing.M) {
	bootsterap.Boot(types.Bootsterap{
		Config: types.ConfigUtils{
			STORAGE_PATH: "./storage_test",
			JWT:          "test",
			DATABASE:     []any{},
		},
	})

	m.Run()
}
