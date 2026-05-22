package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Nome      string    `gorm:"not null;uniqueIndex"`
	Email     string    `gorm:"not null;uniqueIndex"`
	Senha     string    `gorm:"not null"`
	Admin     bool      `gorm:"not null;uniqueIndex;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *User) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
