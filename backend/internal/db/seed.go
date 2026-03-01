package db

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"strconv"

	"github.com/dopp1e/webingo/backend/internal/dto"
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

	// ensure default "user" role exists and get its id
	userRoleReq := dto.RoleCreateRequest{
		Name:        "user",
		Description: "Default user role",
		Level:       10,
	}
	userRole, err := s.Roles.CreateRoleIfNotExists(ctx, userRoleReq)
	if err != nil {
		return err
	}

	// create or ensure a seed user exists to own the boards
	seedUser := &model.User{
		Username: "seed-user",
		Email:    "seed@local",
		RoleID:   userRole.ID,
	}
	// set a default password (hashed)
	if err := seedUser.SetPassword("seed-password"); err != nil {
		return err
	}

	createdUser, err := s.Users.CreateIfNotExists(ctx, seedUser)
	if err != nil {
		return err
	}

	boards := generateBoards(10)

	for _, boardReq := range boards {
		if _, err := s.Boards.CreateBoard(ctx, boardReq, createdUser.ID); err != nil {
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

func generateBoards(num int) []*dto.BoardCreateRequest {
	boards := make([]*dto.BoardCreateRequest, num)

	for i := 0; i < num; i++ {
		board := &dto.BoardCreateRequest{
			Name:        "Board " + strconv.FormatInt(int64(i+1), 10),
			Description: "Description for board " + strconv.FormatInt(int64(i+1), 10),
			Spaces:      generateSpaces(10),
		}
		boards[i] = board
	}

	return boards
}

func generateSpaces(num int) []dto.SpaceRequest {
	spaces := make([]dto.SpaceRequest, num)

	for i := 0; i < num; i++ {
		// Randomly select a bingo space
		space := dto.SpaceRequest{
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
