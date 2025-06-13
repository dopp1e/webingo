package db

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"strconv"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/dopp1e/webingo/backend/internal/service"
)

var bingoSpaces = []string{
	"Free Space",
	"Call a friend",
	"Drink water",
	"Take a walk",
	"Read a book",
	"Send an email",
	"Try something new",
	"Listen to music",
	"Write a note",
	"Organize your desk",
	"Compliment someone",
	"Stretch for 5 minutes",
	"Plan your day",
	"Take a deep breath",
	"Share a joke",
	"Help a colleague",
	"Clean your inbox",
	"Update your calendar",
	"Make a to-do list",
	"Smile at someone",
	"Meditate for 10 minutes",
	"Learn a new word",
	"Water your plants",
	"Take a photo",
	"Write down a goal",
	"Eat a healthy snack",
	"Send a thank you message",
	"Declutter a drawer",
	"Watch a TED talk",
	"Draw something",
	"Try a new recipe",
	"Check in with family",
	"Do 10 push-ups",
	"Update your profile picture",
	"Read the news",
	"Donate to charity",
	"Unsubscribe from an email list",
	"Review your budget",
	"Share a positive story",
	"Take a power nap",
}

func Seed(s service.Service) error {
	ctx := context.Background()

	boards := generateBoards(10)

	for _, boardLoad := range boards {
		board := &model.Board{
			BoardPayload: *boardLoad,
		}
		if err := s.Boards.CreateBoard(ctx, board); err != nil {
			log.Println("failed to create board:", err)
		}
	}

	return nil
}

func randBingoSpace() string {
	j, err := rand.Int(rand.Reader, big.NewInt(int64(len(bingoSpaces))))
	if err != nil {
		panic("failed to generate random number")
	}
	return bingoSpaces[j.Int64()]
}

func generateBoards(num int) []*model.BoardPayload {
	boards := make([]*model.BoardPayload, num)

	for i := 0; i < num; i++ {
		board := &model.BoardPayload{
			Name:        "Board " + strconv.FormatInt(int64(i+1), 10),
			Description: "Description for board " + strconv.FormatInt(int64(i+1), 10),
			Spaces:      generateSpaces(10),
		}
		boards[i] = board
	}

	return boards
}

func generateSpaces(num int) []model.Space {
	spaces := make([]model.Space, num)

	for i := 0; i < num; i++ {
		// Randomly select a bingo space
		space := model.Space{
			Content: randBingoSpace(),
		}
		// Ensure the space is unique
		for j := 0; j < i; j++ {
			for space.Content == spaces[j].Content {
				space.Content = randBingoSpace()
			}
		}

		spaces[i] = space
	}

	return spaces
}
