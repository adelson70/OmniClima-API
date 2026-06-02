package devices

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateDevice(device *Device) (uuid.UUID, error) {
	err := r.db.Create(device).Error
	if err != nil {
		return uuid.Nil, err
	}
	return device.ID, nil
}

func (r *Repository) GetDeviceByID(deviceID uuid.UUID) (string, error) {
	var device Device
	if err := r.db.Where("id = ?", deviceID).First(&device).Error; err != nil {
		return "", err
	}
	return device.PublicKey, nil
}
