package models

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Report represents a report for flagging inappropriate posts or comments
type Report struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID  `gorm:"type:char(36);not null"`
	PostID      *uuid.UUID `gorm:"type:char(36);index;uniqueIndex:unique_report"`
	CommentID   *uuid.UUID `gorm:"type:char(36);index;uniqueIndex:unique_report"`
	Reason      string     `gorm:"type:varchar(64);not null"`
	Reviewed    bool       `gorm:"default:false"`
	ActionTaken string     `gorm:"type:varchar(64);default:''"`

	// Relationships
	Post    *Post    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Comment *Comment `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (Report) TableName() string {
	return "reports"
}

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
