package model

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type FeedEntryType string

const (
	FeedEntryBoard   FeedEntryType = "board"
	FeedEntryComment FeedEntryType = "comment"
)

type FeedEntry struct {
	ID        string        `json:"id"`
	Type      FeedEntryType `json:"type"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	Board     *Board        `json:"board,omitempty"`
	Comment   *BoardComment `json:"comment,omitempty"`
}

type CursorQuery struct {
	Limit  int    `json:"limit" validate:"required,gte=1,lte=20"`
	Cursor string `json:"cursor"` // Optional, base64-encoded
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

type FeedCursor struct {
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`
}

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"required,gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	return fq, nil
}

func EncodeCursor(c FeedCursor) string {
	b, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(b)
}

func DecodeCursor(s string) (*FeedCursor, error) {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var cursor FeedCursor
	err = json.Unmarshal(data, &cursor)
	return &cursor, err
}
