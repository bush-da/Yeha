package models

import (
	"github.com/google/uuid"
)

// PostTag represents the many-to-many relationship between posts and tags.
type PostTag struct {
	PostID uuid.UUID `gorm:"type:char(36);primaryKey" json:"post_id"`
	TagID  uuid.UUID `gorm:"type:char(36);primaryKey" json:"tag_id"`
}

// TableName overrides the default table name for the join table.
func (PostTag) TableName() string {
	return "post_tags"
}
