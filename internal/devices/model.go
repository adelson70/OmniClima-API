package devices

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Device struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	PublicKey string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *Device) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
