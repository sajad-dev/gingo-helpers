package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sajad-dev/gingo-helpers/types"
)

var ConfigStUtils types.ConfigUtils

func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod file not found")
		}
		dir = parent
	}

}

func BootConfig(configSt types.ConfigUtils) {
	ConfigStUtils = configSt

	dir, _ := os.Getwd()
	filepath.Join(dir, "go.mod")

	ConfigStUtils.PROJECT_PATH, _ = FindProjectRoot()

	ConfigStUtils.IMAGE_TEST = fmt.Sprintf("%s/assets/image.png", ConfigStUtils.PROJECT_PATH)
}
