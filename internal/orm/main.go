// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	"errors"

	"github.com/markbates/goth"
	"github.com/yibozhuang/go-gql-server/internal/logger"

	"github.com/yibozhuang/go-gql-server/internal/gql/resolvers/transformations"
	"github.com/yibozhuang/go-gql-server/internal/orm/models"

	"github.com/yibozhuang/go-gql-server/internal/orm/migration"

	"github.com/yibozhuang/go-gql-server/pkg/utils"

	// Imports the database dialect of choice
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

var autoMigrate, logMode, seedDB bool
var dsn, dialect string

// ORM struct to holds the gorm pointer to db
type ORM struct {
	DB *gorm.DB
}

func init() {
	dialect = utils.MustGet("GORM_DIALECT")
	dsn = utils.MustGet("GORM_CONNECTION_DSN")
	seedDB = utils.MustGetBool("GORM_SEED_DB")
	logMode = utils.MustGetBool("GORM_LOGMODE")
	autoMigrate = utils.MustGetBool("GORM_AUTOMIGRATE")
}

// Factory creates a db connection with the selected dialect and connection string
func Factory(cfg *utils.ServerConfig) (*ORM, error) {
	db, err := gorm.Open(cfg.Database.Dialect, cfg.Database.DSN)
	if err != nil {
		logger.Panic("[ORM] err: ", err)
	}
	orm := &ORM{DB: db}
	// Log every SQL command on dev, @prod: this should be disabled? Maybe.
	db.LogMode(cfg.Database.LogMode)
	// Automigrate tables
	if cfg.Database.AutoMigrate {
		err = migration.ServiceAutoMigration(orm.DB)
		if err != nil {
			logger.Error("[ORM.autoMigrate] err: ", err)
		}
	}
	logger.Info("[ORM] Database connection initialized.")
	return orm, nil
}

// FindUserByAPIKey finds the user that is related to the API key
func (o *ORM) FindUserByAPIKey(apiKey string) (*models.User, error) {
	db := o.DB.New()
	uak := &models.UserAPIKey{}
	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}
	if err := db.Preload("User").Where("api_key = ?", apiKey).Find(uak).Error; err != nil {
		return nil, err
	}
	return &uak.User, nil
}

// FindUserByJWT finds the user that is related to the APIKey token
func (o *ORM) FindUserByJWT(email string, provider string, userID string) (*models.User, error) {
	db := o.DB.New()
	up := &models.UserProfile{}
	if provider == "" || userID == "" {
		return nil, errors.New("provider or userId empty")
	}
	if err := db.Preload("User").Where("email  = ? AND provider = ? AND external_user_id = ?", email, provider, userID).First(up).Error; err != nil {
		return nil, err
	}
	return &up.User, nil
}

// UpsertUserProfile saves the user if doesn't exists and adds the OAuth profile
func (o *ORM) UpsertUserProfile(input *goth.User) (*models.User, error) {
	db := o.DB.New()
	u := &models.User{}
	up := &models.UserProfile{}
	u, err := transformations.GothUserToDBUser(input, false)
	if err != nil {
		return nil, err
	}
	if tx := db.Where("email = ?", input.Email).First(u); !tx.RecordNotFound() && tx.Error != nil {
		return nil, tx.Error
	}
	if tx := db.Model(u).Save(u); tx.Error != nil {
		return nil, err
	}
	if tx := db.Where("email = ? AND provider = ? AND external_user_id = ?",
		input.Email, input.Provider, input.UserID).First(up); !tx.RecordNotFound() && tx.Error != nil {
		return nil, err
	}
	up, err = transformations.GothUserToDBUserProfile(input, false)
	if err != nil {
		return nil, err
	}
	up.User = *u
	if tx := db.Model(up).Save(up); tx.Error != nil {
		return nil, tx.Error
	}
	logger.Info(u.ID)
	return u, nil
}
