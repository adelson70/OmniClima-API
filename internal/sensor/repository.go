package sensor

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

func ValidateToken(db *gorm.DB, sensorID, token string) (bool, error) {
	var count int64
	err := db.Model(&Sensor{}).
		Where("id = ? AND token = ?", sensorID, token).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) CreateSensor(s *Sensor) error {
	return r.db.Create(s).Error
}

func (r *Repository) ListSensors(userID uuid.UUID) ([]Sensor, error) {
	var sensors []Sensor
	err := r.db.Where("user_id = ?", userID).Find(&sensors).Error
	return sensors, err
}

func (r *Repository) CreateRecord(input *SensorData) error {
	return r.db.Create(input).Error
}

func (r *Repository) DeleteSensor(sensorID, userID uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", sensorID, userID).Delete(&Sensor{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrSensorNotFound
	}

	return nil
}

func (r *Repository) RenewTokenSensor(sensorID, userID uuid.UUID, newToken string) error {
	result := r.db.Model(&Sensor{}).
		Where("id=? AND user_id=?", sensorID, userID).
		Update("token", newToken)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSensorNotFound
	}
	return nil

}
