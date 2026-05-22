package sensor

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sensor struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UsuarioID uuid.UUID `gorm:"type:uuid;not null;index"`
	Lat       *float64  `gorm:"type:double precision"`
	Lon       *float64  `gorm:"type:double precision"`
	Token     string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SensorDados struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	SensorID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Temp      *float64  `gorm:"type:double precision"`
	Umid      *float64  `gorm:"type:double precision"`
	Rain      bool      `gorm:"not null;index;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Sensor) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (s *SensorDados) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
