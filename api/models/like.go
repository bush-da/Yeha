package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Like represents a 'like' on a post.
type Like struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	AuthorID uuid.UUID `gorm:"type:char(36);not null" json:"author_id"`
	PostID   uuid.UUID `gorm:"type:char(36);not null" json:"post_id"`

	// Relationships
	Author User `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"author"`
	Post   Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post"`
}

// TableName overrides the table name
func (Like) TableName() string {
	return "likes"
}

// BeforeCreate hook to ensure UUID generation for composite keys
func (like *Like) BeforeCreate(tx *gorm.DB) (err error) {
	if like.ID == uuid.Nil {
		like.ID = uuid.New() // Generate a new UUID for AuthorID if not set
	}
	if like.PostID == uuid.Nil {
		like.PostID = uuid.New() // Generate a new UUID for PostID if not set
	}
	return
}
