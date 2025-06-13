package model

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType defines categories of notifications (e.g., "COMMENT", "FOLLOW", "VOTE", "SYSTEM")
type NotificationType string

const (
	NotificationTypeComment NotificationType = "COMMENT"
	NotificationTypeFollow  NotificationType = "FOLLOW"
	NotificationTypeVote    NotificationType = "VOTE"
	NotificationTypeSystem  NotificationType = "SYSTEM"
	NotificationTypeUpdate  NotificationType = "UPDATE" // e.g., board you favorited updated
)

// Notification represents an alert for a user.
type Notification struct {
	Base
	UserID            uuid.UUID        `gorm:"not null;index" json:"userId"` // User to whom the notification belongs
	User              User             `gorm:"foreignKey:UserID"`
	Type              NotificationType `gorm:"not null" json:"type"`        // Type of notification
	Message           string           `gorm:"not null" json:"message"`     // Display message
	IsRead            bool             `gorm:"default:false" json:"isRead"` // Has the user read it?
	RelatedEntityID   *uuid.UUID       `json:"relatedEntityId,omitempty"`   // Optional ID of related entity (e.g., BoardID, CommentID)
	RelatedEntityType *string          `json:"relatedEntityType,omitempty"` // Optional type of related entity (e.g., "Board", "Comment")
	Link              *string          `json:"link,omitempty"`              // Optional URL to navigate to
	AcknowledgedAt    *time.Time       `json:"acknowledgedAt,omitempty"`    // When the user dismissed/acted on it
}
