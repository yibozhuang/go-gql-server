package models

import "github.com/gofrs/uuid"

// User defines a user for the app
type User struct {
	BaseModelSoftDelete
	Email       string  `gorm:"not null;unique_index:idx_email"`
	UserID      *string // External user ID
	Name        *string
	NickName    *string
	FirstName   *string
	LastName    *string
	Location    *string `gorm:"size:1048"`
	Description *string `gorm:"size:1048"`
}

// UserProfile saves all the related OAuth Profiles
type UserProfile struct {
	BaseModelSeq
	Email          string    `gorm:"unique_index:idx_email_provider_external_user_id"`
	UserID         uuid.UUID `gorm:"not null;index"`
	User           User      `gorm:"association_autocreate:false;association_autoupdate:false"`
	Provider       string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id;default:'DB'"` // DB means database or no ExternalUserID
	ExternalUserID string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id"`              // User ID
	Name           string
	NickName       string
	FirstName      string
	LastName       string
	Location       string `gorm:"size:512"`
	AvatarURL      string `gorm:"size:1024"`
	Description    string `gorm:"size:1024"`
	CreatedBy      *User  `gorm:"association_autoupdate:false;association_autocreate:false"`
	UpdatedBy      *User  `gorm:"association_autoupdate:false;association_autocreate:false"`
}

// UserAPIKey generated api keys for the users
type UserAPIKey struct {
	BaseModelSeq
	Name        string
	User        User         `gorm:"association_autocreate:false;association_autoupdate:false"`
	UserID      uuid.UUID    `gorm:"not null;index"`
	APIKey      string       `gorm:"size:128;unique_index"`
	Permissions []Permission `gorm:"many2many:user_api_key_permissions;association_autocreate:false;association_autoupdate:false"`
}
