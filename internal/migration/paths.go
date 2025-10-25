package migration

import (
	"github.com/spitfy/urlshortener/internal/gomodule"
	"os"
	"path/filepath"
)

func getMigrationsDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	moduleRoot, err := gomodule.FindModuleRoot(wd)
	if err != nil {
		return "", err
	}
	migrationsPath := filepath.Join(moduleRoot, "migrations")
	return migrationsPath, nil
}
