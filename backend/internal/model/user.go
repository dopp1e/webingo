package model

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

const (
	ArgonMemory  = 96 * 1024 // 96 MB
	ArgonTime    = 4         // 4 iterations
	ArgonThreads = 1         // 1 thread
	ArgonKeyLen  = 32        // 32 bytes key length
	ArgonSaltLen = 16        // 16 bytes salt length
)

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type User struct {
	Base
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"` // Password stored here, not exposed in JSON output
	IsActive bool   `gorm:"default:false" json:"isActive"`

	RoleID uuid.UUID `json:"roleId"`
	Role   Role      `gorm:"foreignKey:RoleID"`

	Boards        []Board        `gorm:"foreignKey:UserID"`
	Games         []Game         `gorm:"foreignKey:UserID"`
	BoardComments []BoardComment `gorm:"foreignKey:UserID"`
	GameComments  []GameComment  `gorm:"foreignKey:UserID"`
	BoardVotes    []BoardVote    `gorm:"foreignKey:UserID"`
	GameVotes     []GameVote     `gorm:"foreignKey:UserID"`

	BoardCollections []BoardCollection `gorm:"foreignKey:UserID"`
	Notifications    []Notification    `gorm:"foreignKey:UserID"`
	Reports          []Report          `gorm:"foreignKey:ReporterUserID"` // Reports made by this user

	Following []User `gorm:"many2many:follows;joinForeignKey:FollowerID;joinReferences:FollowedUserID"`
	Followers []User `gorm:"many2many:follows;joinForeignKey:FollowedUserID;joinReferences:FollowerID"`
}

func (u *User) SetPassword(text string) error {
	salt, err := generateRandomBytes(ArgonSaltLen)
	if err != nil {
		return err
	}
	hash := argon2.Key([]byte(text), salt, ArgonTime, ArgonMemory, ArgonThreads, ArgonKeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		ArgonMemory,
		ArgonTime,
		ArgonThreads,
		b64Salt,
		b64Hash,
	)

	u.Password = encodedHash

	return nil
}

func (u *User) CheckPassword(text string) bool {
	if u.Password == "" {
		return false
	}

	parts := strings.Split(u.Password, "$")
	if len(parts) != 6 || parts[0] != "$argon2id" {
		return false
	}

	b64Salt := parts[4]
	b64Hash := parts[5]

	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false
	}

	hash := argon2.Key([]byte(text), salt, ArgonTime, ArgonMemory, ArgonThreads, ArgonKeyLen)

	expectedHash, err := base64.RawStdEncoding.DecodeString(b64Hash)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}
