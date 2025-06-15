package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	// No Base struct here, as it's a simple join table for a relationship.

	FollowerID uuid.UUID `gorm:"primaryKey;type:uuid;not null;index:idx_follower_followed,composite:FollowerID,FollowedUserID"`
	Follower   User      `gorm:"foreignKey:FollowerID"`

	FollowedUserID uuid.UUID `gorm:"primaryKey;type:uuid;not null;index:idx_follower_followed,composite:FollowerID,FollowedUserID"`
	FollowedUser   User      `gorm:"foreignKey:FollowedUserID"`

	// Add CreatedAt to track when the follow happened (optional but good for activity feeds)
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
}
