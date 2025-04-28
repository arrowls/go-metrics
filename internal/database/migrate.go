package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func TryMigrations(dbURL string, loggerInst *logrus.Logger) {
	migrationsPath, err := filepath.Abs("./internal/database/migrations")
	if err != nil {
		return
	}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		return
	}
	defer func() {
		if errClose, errDatabase := m.Close(); errClose != nil || errDatabase != nil {
			loggerInst.Errorf("error migrating: %v, %v", errClose, errDatabase)
		}
	}()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return
	}
}
