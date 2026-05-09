package pg

import (
	"context"
	"database/sql"

	"github.com/pressly/goose"
)

type pgMigrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *pgMigrator {
	return &pgMigrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *pgMigrator) Up(ctx context.Context) error {
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}

	return nil
}

func (m *pgMigrator) Down(ctx context.Context) error {
	return goose.Down(m.db, m.migrationsDir)
}

func (m *pgMigrator) Reset(ctx context.Context) error {
	// goose.Reset = откат всех + применение заново
	return goose.Reset(m.db, m.migrationsDir)
}

func (m *pgMigrator) Version(ctx context.Context) (int64, error) {
	version, err := goose.GetDBVersion(m.db)
	if err != nil {
		return 0, err
	}
	return version, nil
}
