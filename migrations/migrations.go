package migrations

import (
	"github.com/jinzhu/gorm"
	migrate "github.com/rubenv/sql-migrate"
	"sync"
)

type dbmigrations struct {
	m          sync.Mutex
	migrations []*migrate.Migration
}

var instance = &dbmigrations{
	m:          sync.Mutex{},
	migrations: make([]*migrate.Migration, 0),
}

func (m *dbmigrations) add(migration *migrate.Migration) {
	m.m.Lock()
	m.migrations = append(m.migrations, migration)
	m.m.Unlock()
}

func InitMySQL(db *gorm.DB) error {
	_, err := migrate.Exec(db.DB(), "mysql", &migrate.MemoryMigrationSource{
		Migrations: instance.migrations,
	}, migrate.Up)
	return err
}

