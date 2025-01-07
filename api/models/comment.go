package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Comment represents a comment on a post.
type Comment struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Content  string    `gorm:"type:text;not null" json:"content"`
	AuthorID uuid.UUID `gorm:"type:char(36);not null" json:"author_id"`
	PostID   uuid.UUID `gorm:"type:char(36);not null" json:"post_id"`

	// Relationships
	Author  User     `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"author"`
	Post    Post     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post"`
	Reports []Report `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE" json:"reports"`
}

// BeforeCreate hook to generate UUID before saving to DB
func (comment *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New() // Generate a new UUID if not set
	}
	return
}

// TableName overrides the table name for the Comment model
func (Comment) TableName() string {
	return "comments"
}
