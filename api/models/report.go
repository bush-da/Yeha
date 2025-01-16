package models

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Report struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updated_at"`

	UserID      uuid.UUID  `gorm:"type:char(36);not null" json:"user_id"`                             // Reporter ID
	AuthorID    uuid.UUID  `gorm:"type:char(36);not null" json:"author_id"`                           // Author of the reported post or comment
	PostID      *uuid.UUID `gorm:"type:char(36);index;constraint:OnDelete:SET NULL" json:"post_id"`   // Nullable post ID
	CommentID   *uuid.UUID `gorm:"type:char(36);index;constraint:OnDelete:CASCADE" json:"comment_id"` // Nullable comment ID
	Reason      string     `gorm:"type:varchar(64);not null" json:"reason"`                           // Report reason
	Reviewed    bool       `gorm:"default:false" json:"reviewed"`
	ActionTaken string     `gorm:"type:varchar(64);default:''" json:"action_taken"`

	// Relationships
	Post    *Post    `gorm:"foreignKey:PostID" json:"post"`
	Comment *Comment `gorm:"foreignKey:CommentID" json:"comment"`
}

// TableName overrides the table name
func (Report) TableName() string {
	return "reports"
}

// BeforeCreate hook to validate Report fields and generate UUIDs before saving to DB
func (report *Report) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID if not set
	if report.ID == uuid.Nil {
		report.ID = uuid.New()
	}

	// Validate that either PostID or CommentID is set, but not both
	if (report.PostID == nil || *report.PostID == uuid.Nil) && (report.CommentID == nil || *report.CommentID == uuid.Nil) {
		return fmt.Errorf("either PostID or CommentID must be set")
	}
	if (report.PostID != nil && *report.PostID != uuid.Nil) && (report.CommentID != nil && *report.CommentID != uuid.Nil) {
		return fmt.Errorf("cannot have both PostID and CommentID set")
	}

	// Verify PostID or CommentID exists in the database
	if report.PostID != nil && *report.PostID != uuid.Nil {
		var post Post
		if err := tx.First(&post, "id = ?", report.PostID).Error; err != nil {
			return fmt.Errorf("invalid PostID")
		}
	}
	if report.CommentID != nil && *report.CommentID != uuid.Nil {
		var comment Comment
		if err := tx.First(&comment, "id = ?", report.CommentID).Error; err != nil {
			return fmt.Errorf("invalid CommentID")
		}
	}

	return nil
}
