package dto

import (
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
)

// ReportCreateRequest for submitting a new report.
type ReportCreateRequest struct {
	ReportedEntityID   uuid.UUID                `json:"reportedEntityId" validate:"required"`
	ReportedEntityType model.ReportedEntityType `json:"reportedEntityType" validate:"required,oneof=BOARD COMMENT USER GAME"`
	Reason             model.ReportReason       `json:"reason" validate:"required,oneof=SPAM HARASSMENT NUDITY HATE_SPEECH OTHER"`
	Description        string                   `json:"description,omitempty" validate:"max=1000"`
}
