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
}

type UserStoreInterface interface {
	Create(context.Context, *model.User) error
	GetByUsername(context.Context, string) (*model.User, error)
	Exists(context.Context, string) (bool, error)
}

type SpaceStoreInterface interface {
	FindOrCreateByContent(ctx context.Context, space *model.Space) (*model.Space, error)
}

type GameStoreInterface interface {
	Create(context.Context, *model.Game) error
}

type RoleInterface interface {
	Create(context.Context, *model.Role) error
	GetByName(context.Context, string) (*model.Role, error)
	Exists(context.Context, string) (bool, error)
}

type StorageInterface interface {
	Boards() BoardStoreInterface
	Users() UserStoreInterface
	Spaces() SpaceStoreInterface
	Games() GameStoreInterface
	Roles() RoleInterface
}

type Storage struct {
	boards BoardStoreInterface
	users  UserStoreInterface
	spaces SpaceStoreInterface
	games  GameStoreInterface
	roles  RoleInterface
}

func (s *Storage) Boards() BoardStoreInterface { return s.boards }
func (s *Storage) Users() UserStoreInterface   { return s.users }
func (s *Storage) Spaces() SpaceStoreInterface { return s.spaces }
func (s *Storage) Games() GameStoreInterface   { return s.games }
func (s *Storage) Roles() RoleInterface        { return s.roles }

type TransactionalStorage interface {
	StorageInterface
	Commit() error
	Rollback() error
	Model(interface{}) *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		boards: &BoardStore{db},
		users:  &UserStore{db},
		spaces: &SpaceStore{db},
		games:  &GameStore{db},
		roles:  &RoleStore{db},
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
