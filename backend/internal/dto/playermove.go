package dto

// PlayerMoveCreateRequest for creating a player move.
type PlayerMoveCreateRequest struct {
	CellH     int `json:"cellH" validate:"required"`
	CellW     int `json:"cellW" validate:"required"`
	Timestamp int `json:"timestamp" validate:"required"`
}
