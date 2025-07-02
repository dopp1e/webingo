package store

import (
	"context"

	"github.com/dopp1e/webingo/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardStoreInterface interface {
	Create(context.Context, *model.Board) error
	GetById(context.Context, uuid.UUID) (*model.Board, error)
	Update(context.Context, *model.Board) error
	Delete(context.Context, *model.Board) error
	Exists(context.Context, uuid.UUID) (bool, error)
	GetByCreatorIDs(ctx context.Context, creatorIDs []uuid.UUID, offset int, limit int, sort string) ([]model.Board, error)
}

type BoardCommentStoreInterface interface {
	Create(context.Context, *model.BoardComment) error
	GetByID(context.Context, uuid.UUID) (*model.BoardComment, error)
}

type UserStoreInterface interface {
	Create(context.Context, *model.User) error
	GetByUsername(context.Context, string) (*model.User, error)
	Exists(context.Context, string) (bool, error)
	GetByID(context.Context, uuid.UUID) (*model.User, error)
	GetFollowedUserIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
}

type SpaceStoreInterface interface {
	FindOrCreateByContent(ctx context.Context, space *model.Space) (*model.Space, error)
}

type GameStoreInterface interface {
	Create(context.Context, *model.Game) error
	GetByID(context.Context, uuid.UUID) (*model.Game, error)
	GetGameMarkedFieldAndBingoCount(ctx context.Context, gameID uuid.UUID) (int, int, error)
}

type RoleStoreInterface interface {
	Create(context.Context, *model.Role) error
	GetByName(context.Context, string) (*model.Role, error)
	Exists(context.Context, string) (bool, error)
}

type FollowStoreInterface interface {
	Create(context.Context, *model.Follow) error
	GetFollowers(context.Context, uuid.UUID) ([]model.Follow, error)
	GetFollowing(context.Context, uuid.UUID) ([]model.Follow, error)
	Delete(context.Context, *model.Follow) error
}

type ActivityStoreInterface interface {
	Create(context.Context, *model.Activity) error
	GetPaginatedActivitiesByActors(ctx context.Context, actorIDs []uuid.UUID, fq model.PaginatedFeedQuery) ([]model.Activity, error)
}

type InvitationStoreInterface interface {
	Create(context.Context, *model.Invitation) error
}

type StorageInterface interface {
	Boards() BoardStoreInterface
	BoardComments() BoardCommentStoreInterface
	Users() UserStoreInterface
	Spaces() SpaceStoreInterface
	Games() GameStoreInterface
	Roles() RoleStoreInterface
	Follows() FollowStoreInterface
	Activities() ActivityStoreInterface
	Invitations() InvitationStoreInterface
}

type Storage struct {
	boards        BoardStoreInterface
	boardComments BoardCommentStoreInterface
	users         UserStoreInterface
	spaces        SpaceStoreInterface
	games         GameStoreInterface
	roles         RoleStoreInterface
	follows       FollowStoreInterface
	activity      ActivityStoreInterface
	invitation    InvitationStoreInterface
}

func (s *Storage) Boards() BoardStoreInterface               { return s.boards }
func (s *Storage) BoardComments() BoardCommentStoreInterface { return s.boardComments }
func (s *Storage) Users() UserStoreInterface                 { return s.users }
func (s *Storage) Spaces() SpaceStoreInterface               { return s.spaces }
func (s *Storage) Games() GameStoreInterface                 { return s.games }
func (s *Storage) Roles() RoleStoreInterface                 { return s.roles }
func (s *Storage) Follows() FollowStoreInterface             { return s.follows }
func (s *Storage) Activities() ActivityStoreInterface        { return s.activity }
func (s *Storage) Invitations() InvitationStoreInterface     { return s.invitation }

type TransactionalStorage interface {
	StorageInterface
	Commit() error
	Rollback() error
	Model(interface{}) *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		boards:        &BoardStore{db},
		boardComments: &BoardCommentStore{db},
		users:         &UserStore{db},
		spaces:        &SpaceStore{db},
		games:         &GameStore{db},
		roles:         &RoleStore{db},
		follows:       &FollowStore{db},
		activity:      &ActivityStore{db},
		invitation:    &InvitationStore{db},
	}
}

type storageTx struct {
	*Storage
	tx *gorm.DB
}

func (s *storageTx) Commit() error {
	return s.tx.Commit().Error
}

func (s *storageTx) Rollback() error {
	return s.tx.Rollback().Error
}

func (s *storageTx) Model(value interface{}) *gorm.DB {
	return s.tx.Model(value)
}

func StartTransaction(ctx context.Context, db *gorm.DB) (TransactionalStorage, error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &storageTx{
		Storage: NewStorage(tx),
		tx:      tx,
	}, nil
}
