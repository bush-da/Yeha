package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents a blog post.
type Post struct {
	ID       uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Title    string    `gorm:"type:varchar(255);not null" json:"title"`
	AuthorID uuid.UUID `gorm:"type:char(36);not null" json:"author_id"`

	// Relationships
	// Reports  []Report  `gorm:"foreignKey:PostID" json:"reports"`
	// Contents []Content `gorm:"foreignKey:PostID" json:"contents"`
	// Comments []Comment `gorm:"foreignKey:PostID" json:"comments"`
	// Likes    []Like    `gorm:"foreignKey:PostID" json:"likes"`
	// Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags"`
	// Post Model Relationships
	Reports  []Report  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"reports"`
	Contents []Content `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"contents"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments"`
	Likes    []Like    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"likes"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags"` // No need for CASCADE here; managed in PostTag

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
func (p *Post) IsLiked(userID uuid.UUID, db *gorm.DB) bool {
	var like Like
	result := db.Where("post_id = ? AND author_id = ?", p.ID, userID).First(&like)
	return result.RowsAffected > 0
}

// AddTag adds a tag to the post.
func (p *Post) AddTag(tag Tag, db *gorm.DB) error {
	err := db.Model(&p).Association("Tags").Append(&tag)
	return err
}

// // RemoveTag removes a tag from the post.
// func (p *Post) RemoveTag(tag Tag, db *gorm.DB) error {
// 	err := db.Model(&p).Association("Tags").Delete(&tag)
// 	return err
// }

// //BeforeDelete cleans up unused tags before deleting the post.
// func (p *Post) BeforeDelete(tx *gorm.DB) error {
// 	var tags []Tag
// 	tx.Model(&p).Association("Tags").Find(&tags)

// 	// Delete the post
// 	if err := tx.Delete(&p).Error; err != nil {
// 		return err
// 	}

// 	// Clean up tags with no posts
// 	for _, tag := range tags {
// 		var count int64
// 		// Correct the call to Count() without passing an argument
// 		tx.Model(&tag).Association("Posts").Count()

// 		// If the count is zero, delete the tag
// 		if count == 0 {
// 			if err := tx.Delete(&tag).Error; err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// BeforeCreate hook ensures UUID is set for ID field
func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if post.ID == uuid.Nil {
		post.ID = uuid.New() // Generate a new UUID for ID if not set
	}
	return
}
