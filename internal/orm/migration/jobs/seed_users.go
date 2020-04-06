package jobs

import (
	"github.com/jinzhu/gorm"
	"github.com/yibozhuang/go-gql-server/internal/orm/models"
	"gopkg.in/gormigrate.v1"
)

var (
	uname       = "Test User"
	fname       = "Test"
	lname       = "User"
	email       = "test@test.com"
	nname       = "Foo Bar"
	description = "This is the first user"
	location    = "Some location"

	firstUser *models.User = &models.User{
		Email:       email,
		Name:        &uname,
		FirstName:   &fname,
		LastName:    &lname,
		NickName:    &nname,
		Description: &description,
		Location:    &location,
	}
)

// SeedUsers inserts the first users
var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		return db.Create(&firstUser).Error
	},
	Rollback: func(db *gorm.DB) error {
		return db.Delete(&firstUser).Error
	},
}
