package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameStore struct {
	db *gorm.DB
}

func (s *GameStore) Create(ctx context.Context, game *model.Game) error {
	result := s.db.WithContext(ctx).Create(game)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *GameStore) GetByID(ctx context.Context, id uuid.UUID) (*model.Game, error) {
	var game model.Game
	result := s.db.WithContext(ctx).First(&game, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found
		}
		return nil, result.Error
	}

	return &game, nil
}

func fillBoard(board [][]bool, moves []model.PlayerMove, height int, width int) (int, int) {
	moveCount := 0
	bingoCount := 0
	for _, move := range moves {
		okHeight := move.CellH >= 0 && move.CellH < height
		okWidth := move.CellW >= 0 && move.CellW < width
		if !okHeight || !okWidth {
			continue // Skip invalid moves
		}

		board[move.CellH][move.CellW] = !board[move.CellH][move.CellW] // Toggle the marked state
		if board[move.CellH][move.CellW] {
			moveCount++ // Increment count if the field is marked
		} else {
			moveCount-- // Decrement count if the field is unmarked
		}
	}

	return moveCount, bingoCount
}

func (s *GameStore) GetGameMarkedFieldAndBingoCount(ctx context.Context, gameID uuid.UUID) (int, int, error) {
	var moves []model.PlayerMove
	result := s.db.WithContext(ctx).
		Model(&model.PlayerMove{}).
		Where("game_id = ?", gameID).
		Find(&moves)

	if result.Error != nil {
		return 0, 0, result.Error
	}

	game, err := s.GetByID(ctx, gameID)
	if err != nil || game == nil {
		return 0, 0, err
	}

	var markedFields [][]bool = make([][]bool, game.Height)
	for i := range markedFields {
		markedFields[i] = make([]bool, game.Width)
	}
	moveCount, bingoCount := fillBoard(markedFields, moves, game.Height, game.Width)

	return moveCount, bingoCount, nil
}
