package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag represents a tag for posts (e.g., 'Python', 'Web Development').
type Tag struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(45);unique;not null;index" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Many-to-Many Relationship with Post
	Posts []Post `gorm:"many2many:post_tags;" json:"posts"`
}

// BeforeCreate hook to generate UUID before saving to DB
func (tag *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if tag.ID == uuid.Nil {
		tag.ID = uuid.New() // Generate a new UUID if not set
	}
	return
}
