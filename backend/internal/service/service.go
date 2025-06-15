package service

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/dto"
	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	Boards interface {
		CreateBoard(context.Context, *dto.BoardCreateRequest, uuid.UUID) (*model.Board, error)
		GetById(context.Context, uuid.UUID) (*model.Board, error)
		UpdateBoard(context.Context, *dto.BoardUpdateRequest, uuid.UUID) (*model.Board, error)
		DeleteBoard(context.Context, uuid.UUID) error
		AddComment(ctx context.Context, boardId uuid.UUID, userId uuid.UUID, request *dto.CommentPutRequest) (*model.BoardComment, error)
	}
	Users interface {
		Create(context.Context, *model.User) error
		GetByUsername(context.Context, string) (*model.User, error)
		Exists(context.Context, string) (bool, error)
		CreateIfNotExists(context.Context, *model.User) (*model.User, error)
		GetByID(context.Context, uuid.UUID) (*model.User, error)
		FollowUser(ctx context.Context, followerID, followedID uuid.UUID) error
		UnfollowUser(ctx context.Context, followerID, followedID uuid.UUID) error
	}
	Roles interface {
		CreateRole(context.Context, *model.Role) error
		GetRoleByName(context.Context, string) (*model.Role, error)
		CreateRoleIfNotExists(context.Context, *model.Role) (*model.Role, error)
	}
	Feed interface {
		GetFeed(ctx context.Context, userId uuid.UUID) ([]model.Board, error)
	}
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
		Boards: &BoardService{
			db: db,
		},
		Users: &UserService{
			db: db,
		},
		Roles: &RoleService{
			db: db,
		},
		Feed: &FeedService{
			db: db,
		},
	}
}
