package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Content represents the content that will be posted on the website.
type Content struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	FileURL   string    `gorm:"type:text" json:"file_url"`
	Paragraph int       `gorm:"not null" json:"paragraph"`
	Content   string    `gorm:"type:text" json:"content"`
	PostID    uuid.UUID `gorm:"type:char(36);not null;index" json:"post_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship with Post
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post"`
}

// BeforeCreate hook to generate UUID before saving to DB
func (content *Content) BeforeCreate(tx *gorm.DB) (err error) {
	if content.ID == uuid.Nil {
		content.ID = uuid.New() // Generate a new UUID if not set
	}
	return
}

// TableName overrides the table name.
func (Content) TableName() string {
	return "contents"
}
