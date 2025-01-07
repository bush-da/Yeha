package models

import (
	"github.com/google/uuid"
)

// Follower represents the user-following relationship.
type Follower struct {
	FollowerID uuid.UUID `gorm:"type:char(36);not null;primaryKey" json:"follower_id"`
	FollowedID uuid.UUID `gorm:"type:char(36);not null;primaryKey" json:"followed_id"`

	// Relationships
	Follower User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"follower"`
	Followed User `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE" json:"followed"`
}

// TableName overrides the table name.
func (Follower) TableName() string {
	return "followers"
}
