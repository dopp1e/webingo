package model

import "github.com/google/uuid"

// ActivityVerb defines the type of action performed in an activity.
type ActivityVerb string

const (
	ActivityVerbCreatedBoard     ActivityVerb = "created_board"
	ActivityVerbCreatedGame      ActivityVerb = "created_game"
	ActivityVerbCommentedOnBoard ActivityVerb = "commented_on_board"
	ActivityVerbCommentedOnGame  ActivityVerb = "commented_on_game"
	ActivityVerbVotedOnBoard     ActivityVerb = "voted_on_board"
	ActivityVerbVotedOnGame      ActivityVerb = "voted_on_game"
	// More verbs might be added here.
)

// ActivityTargetType defines the type of entity targeted by an activity.
type ActivityTargetType string

const (
	ActivityTargetTypeBoard   ActivityTargetType = "board"
	ActivityTargetTypeGame    ActivityTargetType = "game"
	ActivityTargetTypeComment ActivityTargetType = "comment"
	// More target types might be added here.
)

// Activity represents a specific action or event in the system, suitable for a user's feed.
type Activity struct {
	Base
	ActorID     uuid.UUID           `gorm:"not null;index" json:"actorId"` // The user who performed the action
	Actor       User                `gorm:"foreignKey:ActorID"`
	Verb        ActivityVerb        `gorm:"not null" json:"verb"`           // What action was performed (e.g., "created_board")
	TargetID    uuid.UUID           `gorm:"not null;index" json:"targetId"` // The ID of the primary entity involved (e.g., Board ID, Comment ID)
	TargetType  ActivityTargetType  `gorm:"not null" json:"targetType"`     // The type of the target entity (e.g., "board", "comment")
	ContextID   *uuid.UUID          `json:"contextId,omitempty"`            // Optional: ID of a related entity (e.g., BoardID if Target is a Comment)
	ContextType *ActivityTargetType `json:"contextType,omitempty"`          // Optional: Type of related entity (e.g., "board" if Target is a Comment)

	// Optional TODO: Denormalized fields for performance
}
