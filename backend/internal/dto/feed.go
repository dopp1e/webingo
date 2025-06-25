package dto

import (
	"time"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
)

// FeedUserData represents a user in the feed, containing minimal information.
type FeedUserData struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// FeedItem represents a single entry in a user's activity feed.
// It combines the activity metadata with the actual target entity's details.
type FeedItem struct {
	ID          uuid.UUID                 `json:"id"`
	Actor       FeedUserData              `json:"actor"`
	Verb        model.ActivityVerb        `json:"verb"`
	CreatedAt   time.Time                 `json:"createdAt"`
	TargetID    uuid.UUID                 `json:"targetId"`
	TargetType  model.ActivityTargetType  `json:"targetType"`
	ContextID   *uuid.UUID                `json:"contextId,omitempty"`
	ContextType *model.ActivityTargetType `json:"contextType,omitempty"`

	// This field will hold the actual hydrated entity (Board, Game, Comment etc.)
	// Use an interface{} and rely on JSON marshalling to handle different types.
	Target interface{} `json:"target"`
}

type BoardFeedItem struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
}

type GameFeedItem struct {
	ID           uuid.UUID `json:"id"`
	YoutubeID    *string   `json:"youtubeId"`
	Height       int       `json:"height"`
	Width        int       `json:"width"`
	MarkedFields int       `json:"markedFields"`
	BingoCount   int       `json:"bingoCount"`
}

type CommentFeedItem struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}
