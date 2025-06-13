package dto

// CommentPutRequest for creating a new comment.
type CommentPutRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}
