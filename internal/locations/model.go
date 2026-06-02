package locations

import (
	"OmniClima/internal/user"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Country   string    `gorm:"not null;index"`
	State     string    `gorm:"not null;index"`
	City      string    `gorm:"not null;index"`
	Lat       float64   `gorm:"not null;type:double precision"`
	Lon       float64   `gorm:"not null;type:double precision"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      user.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (s *Location) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
