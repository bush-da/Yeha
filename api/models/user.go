package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the Yeha application
type User struct {
	ID             uuid.UUID `gorm:"type:char(36);primaryKey"`
	Username       string    `gorm:"size:45;not null"`
	Email          string    `gorm:"size:45;unique;not null"`
	Password       string    `gorm:"size:255;not null"`
	ProfilePicture string    `gorm:"size:255"`
	Gender         string    `gorm:"type:enum('male', 'female', 'other');not null;default:'other'" json:"gender"`
	Bio            string    `gorm:"type:text"`
	IsAdmin        bool      `gorm:"default:false"`

	// Relationships
	Posts     []Post     `gorm:"foreignKey:AuthorID"`
	Comments  []Comment  `gorm:"foreignKey:AuthorID"`
	Likes     []Like     `gorm:"foreignKey:AuthorID"`
	Following []Follower `gorm:"foreignKey:FollowerID"`
	Followers []Follower `gorm:"foreignKey:FollowedID"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook to generate UUIDs before creating a record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID before creating the record
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	validGenders := map[string]bool{"male": true, "female": true, "other": true}
	if _, valid := validGenders[u.Gender]; !valid {
		u.Gender = "other" // Default if invalid
	}
	return nil
}
