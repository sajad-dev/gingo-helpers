package bootstrap

import (
	"github.com/sajad-dev/gingo-helpers/internal/config"
	"github.com/sajad-dev/gingo-helpers/types"
)

func Boot(bootSt types.Bootsterap) {
	config.BootConfig(bootSt.Config)
}
