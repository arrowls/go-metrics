package database

import (
	"fmt"
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

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dbURL)
	if err != nil {
		return
	}

	defer func() {
		if errClose, errDatabase := m.Close(); errClose != nil || errDatabase != nil {
			loggerInst.Errorf("error migrating: %v, %v", errClose, errDatabase)
		}
	}()

	err = m.Up()
	if err != nil {
		loggerInst.Errorf("error migrating: %v", err)
	}
}
