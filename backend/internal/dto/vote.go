package dto

// VoteRequest is used for both creating and updating a vote.
// The target entity ID (BoardID or GameID) and UserID will come from the URL/context.
type VoteRequest struct {
	IsUpvote bool `json:"isUpvote" validate:"boolean"` // true for upvote, false for downvote
}
