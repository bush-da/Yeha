package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents a blog post.
type Post struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	AuthorID  uuid.UUID `gorm:"type:char(36);not null" json:"author_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Reports  []Report  `gorm:"foreignKey:PostID;constraint:OnDelete:SET NULL" json:"reports"`
	Contents []Content `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"contents"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:SET NULL" json:"comments"`
	Likes    []Like    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"likes"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags"`

	Author User `gorm:"foreignKey:AuthorID" json:"author"`
}

// CountComments returns the number of comments for the post.
func (p *Post) CountComments() int {
	return len(p.Comments)
}

// CountLikes returns the number of likes for the post.
func (p *Post) CountLikes() int {
	return len(p.Likes)
}

// IsLiked checks whether a user has liked the post.
func (p *Post) is_liked(userID uuid.UUID, db *gorm.DB) bool {
	var like Like
	result := db.Where("post_id = ? AND author_id = ?", p.ID, userID).First(&like)
	return result.RowsAffected > 0
}

// AddTag adds a tag to the post.
func (p *Post) AddTag(tag Tag, db *gorm.DB) error {
	err := db.Model(&p).Association("Tags").Append(&tag)
	return err
}

// RemoveTag removes a tag from the post.
func (p *Post) RemoveTag(tag Tag, db *gorm.DB) error {
	err := db.Model(&p).Association("Tags").Delete(&tag)
	return err
}

// BeforeCreate hook ensures UUID is set for ID field
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New() // Generate a new UUID for ID if not set
	}
	return
}
