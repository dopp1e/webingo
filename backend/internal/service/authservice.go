package service

import (
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}
