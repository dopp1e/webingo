package model

import "github.com/google/uuid"

// ReportType defines the type of entity being reported.
type ReportedEntityType string

const (
	ReportedEntityTypeBoard   ReportedEntityType = "BOARD"
	ReportedEntityTypeComment ReportedEntityType = "COMMENT" // For both BoardComment and GameComment
	ReportedEntityTypeUser    ReportedEntityType = "USER"
	ReportedEntityTypeGame    ReportedEntityType = "GAME"
)

// ReportReason defines common reasons for reporting.
type ReportReason string

const (
	ReportReasonSpam       ReportReason = "SPAM"
	ReportReasonHarassment ReportReason = "HARASSMENT"
	ReportReasonNudity     ReportReason = "NUDITY"
	ReportReasonHateSpeech ReportReason = "HATE_SPEECH"
	ReportReasonOther      ReportReason = "OTHER"
)

// Report represents a user's report against another user or content.
type Report struct {
	Base
	ReporterUserID     uuid.UUID          `gorm:"not null;index" json:"reporterUserId"` // User who submitted the report
	ReporterUser       User               `gorm:"foreignKey:ReporterUserID"`
	ReportedEntityID   uuid.UUID          `gorm:"not null;index" json:"reportedEntityId"` // ID of the entity being reported
	ReportedEntityType ReportedEntityType `gorm:"not null" json:"reportedEntityType"`     // Type of the entity (e.g., "BOARD", "COMMENT", "USER")
	Reason             ReportReason       `gorm:"not null" json:"reason"`                 // Reason for the report
	Description        string             `json:"description,omitempty"`                  // Optional detailed description
	IsResolved         bool               `gorm:"default:false" json:"isResolved"`        // Has the report been actioned by moderation?
	ModeratorUserID    *uuid.UUID         `json:"moderatorUserId,omitempty"`              // ID of the moderator who handled it
	ModeratorUser      *User              `gorm:"foreignKey:ModeratorUserID"`
	ResolutionNotes    *string            `json:"resolutionNotes,omitempty"` // Notes from the moderator
}
