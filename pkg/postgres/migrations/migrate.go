package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

func Up(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		err = fmt.Errorf("failed to set dialect: %v", err)
		logrus.WithError(err).Error("Migration UP: failed to set database dialect")
		return err
	}
	logrus.Info("Migration UP: PostgreSQL dialect set successfully")

	if err := goose.Up(db, "sql"); err != nil {
		err = fmt.Errorf("failed to apply UP migrations: %v", err)
		logrus.WithError(err).Error("Migration UP: failed to apply migrations")
		return err
	}
	logrus.Info("Migration UP: migrations applied successfully")
	
	return nil
}

func Down(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		err = fmt.Errorf("failed to set dialect: %v", err)
		logrus.WithError(err).Error("Migration DOWN: failed to set database dialect")
		return err
	}
	logrus.Info("Migration DOWN: PostgreSQL dialect set successfully")

	if err := goose.Down(db, "sql"); err != nil {
		err = fmt.Errorf("failed to apply DOWN migrations: %v", err)
		logrus.WithError(err).Error("Migration DOWN: failed to apply migrations")
		return err
	}
	logrus.Info("Migration DOWN: migrations rolled back successfully")

	return nil
}

