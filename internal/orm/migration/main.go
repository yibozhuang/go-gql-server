package migration

import (
	"fmt"

	log "github.com/yibozhuang/go-gql-server/internal/logger"

	"github.com/jinzhu/gorm"
	"github.com/yibozhuang/go-gql-server/internal/orm/migration/jobs"
	"github.com/yibozhuang/go-gql-server/internal/orm/models"
	"gopkg.in/gormigrate.v1"
)

func updateMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
	).Error
}

// ServiceAutoMigration migrates all the tables and modifications to the connected source
func ServiceAutoMigration(db *gorm.DB) error {
	// Keep a list of migrations here
	m := gormigrate.New(db, gormigrate.DefaultOptions, nil)

	m.InitSchema(func(db *gorm.DB) error {
		log.Info("[Migration.InitSchema] Initializing database schema")
		switch db.Dialect().GetName() {
		case "postgres":
			// Let's create the UUID extension, the user has to have superuser permission for this
			db.Exec("CREATE EXTENSION \"uuid-ossp\";")
		}
		if err := updateMigration(db); err != nil {
			return fmt.Errorf("[Migration.InitSchema]: %v", err)
		}
		// Add more jobs, etc. here ...
		return nil
	})
	m.Migrate()

	if err := updateMigration(db); err != nil {
		return err
	}
	m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		jobs.SeedUsers,
	})
	return m.Migrate()
}
